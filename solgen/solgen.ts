import { FunctionAnalysis, getFunctionsFromGo } from "./contracts_parser";
import * as fs from "fs";
import {
  AsTypeFunction,
  BindingLibraryType,
  OperatorBinding,
  OperatorOverloadDecl,
  PostFix,
  preamble,
  SolTemplate1Arg,
  SolTemplate2Arg,
  SolTemplate3Arg,
  genAbiFile,
  CastBinding,
  SealFromType,
  DecryptBinding,
  IsOperationAllowed,
  RandomGenericFunction,
  RandomFunctions,
  SolTemplateDecrypt,
  SealTypedFromType,
} from "./templates/library";

import {
  testContract0Args,
  testContract1Arg,
  testContract2Arg,
  testContract3Arg,
  testContract2ArgBoolRes,
  testContractReencrypt,
  testContractReq,
  AsTypeTestingContract,
  testContractDecrypt,
  testContractSealTyped,
} from "./templates/testContracts";

import {
  AsTypeBenchmarkContract,
  benchContract1Arg,
  benchContract2Arg,
  benchContract3Arg,
  benchContractReencrypt,
} from "./templates/benchContracts";

import {
  AllTypes,
  BindMathOperators,
  bitwiseAndLogicalOperators,
  EInputType,
  EPlaintextType,
  ShorthandOperations,
  valueIsEncrypted,
  isComparisonType,
  isBitwiseOp,
  SEALING_FUNCTION_NAME,
  capitalize,
  UTypeSealedOutputMap,
  EUintType,
  SEALING_TYPED_FUNCTION_NAME,
} from "./common";

interface FunctionMetadata {
  functionName: string;
  inputCount: number;
  hasDifferentInputTypes: boolean;
  returnValueType?: string;
  inputs: AllTypes[];
  isBooleanMathOp: boolean;
}

const generateMetadataPayload = async (): Promise<FunctionMetadata[]> => {
  const result = await getFunctionsFromGo("../precompiles/contracts.go");
  const resultWithInjected = injectMetadataAdditionalFunctions(result);

  return resultWithInjected.map((value) => {
    return {
      functionName: value.name,
      hasDifferentInputTypes: !value.needsSameType,
      inputCount: value.paramsCount,
      returnValueType: value.returnType,
      inputs: value.inputTypes,
      isBooleanMathOp: value.isBooleanMathOp,
    };
  });
};

const injectMetadataAdditionalFunctions = (fns: FunctionAnalysis[]) => {
  // List of additional functions to be generated that depend upon the parsed `go` functions
  // Dependents will be inserted in the generated contract immediately following the parent function
  const fnDependents: Record<string, FunctionAnalysis[]> = {
    [SEALING_FUNCTION_NAME]: [
      {
        name: "sealoutputTyped",
        paramsCount: 2,
        needsSameType: false,
        // Is replaced in `getReturnType` with `SealedBool`/`SealedUint`/`SealedAddress` based on input0
        returnType: "SealedStruct",
        inputTypes: ["encrypted", "bytes32"],
        isBooleanMathOp: true,
      },
    ],
  };
  return fns.flatMap((fn) => fnDependents[fn.name] != null ? [fn, ...fnDependents[fn.name]] : fn);
}

// Function to generate all combinations of parameters.
function generateCombinations(
  arr: string[][],
  current: string[] = [],
  index: number = 0
): string[][] {
  if (index === arr.length) {
    // Only add to the result if there are elements in the current combination
    return current.length === 0 ? [] : [current];
  }

  let result: string[][] = [];
  // Add each element of the current sub-array to the combination
  for (let item of arr[index]) {
    result = result.concat(
      generateCombinations(arr, current.concat(item), index + 1)
    );
  }
  return result;
}

const getReturnType = (
  inputs: string[],
  isBooleanMathOp: boolean,
  returnType?: string
) => {
  if (returnType === "plaintext") {
    if (inputs.length != 1) {
      throw new Error(
        "expecting exactly one input for functions returning plaintext"
      );
    }

    let inputType = inputs[0].split(" ")[1];
    if (inputType[0] !== "e") {
      throw new Error("expecting encrypted input for plaintext output");
    }

    return inputType.slice(1);
  }

  // `sealoutputTyped` determine and replace output type based on input0 type
  if (returnType && returnType === "SealedStruct") {
    return `${
      UTypeSealedOutputMap[inputs[0].replace("input0 ", "") as EUintType]
    } memory`;
  }

  if (returnType && returnType !== "encrypted") {
    return returnType;
  }

  if (inputs.includes("bytes") || inputs.includes("bytes32")) {
    return "bytes";
  }

  // if (isBooleanMathOp) {
  //     return "ebool";
  // }

  let maxRank = 0;
  for (let input of inputs) {
    let inputType = input.split(" ")[1];
    maxRank = Math.max(
      EInputType.indexOf(inputType),
      EPlaintextType.indexOf(inputType),
      maxRank
    );
  }

  return EInputType[maxRank];
};

function getAllFunctionDeclarations(
  functionName: string,
  functions: string[][],
  isBooleanMathOp: boolean,
  returnValueType?: string
): string[] {
  let functionDecl = `function ${functionName}`;

  // Generate all combinations of input parameters.
  let allCombinations = generateCombinations(functions);

  // Create function declarations for each combination.
  return allCombinations.map((combination) => {
    let returnType = getReturnType(
      combination,
      isBooleanMathOp,
      returnValueType
    );
    let returnStr = `internal pure returns (${returnType});`;

    return `${functionDecl}(${combination.join(", ")}) ${returnStr}`;
  });
}

const getOperator = (functionName: string): string | undefined => {
  return (
    ShorthandOperations.find((operation) => operation.func === functionName)
      ?.operator ?? undefined
  );
};

/** Generates a Solidity test contract based on the provided metadata */
const generateSolidityTestContract = (metadata: FunctionMetadata): string[] => {
  const {
    functionName,
    inputCount,
    hasDifferentInputTypes,
    returnValueType,
    inputs,
    isBooleanMathOp,
  } = metadata;

  if (functionName === "req") {
    return testContractReq();
  }

  if (functionName === "decrypt") {
    return testContractDecrypt();
  }

  if (functionName === SEALING_FUNCTION_NAME) {
    return testContractReencrypt();
  }

  if (functionName === SEALING_TYPED_FUNCTION_NAME) {
    return testContractSealTyped();
  }

  if (inputCount === 0) {
    return testContract0Args(functionName);
  }

  if (
    inputCount === 1 &&
    inputs[0] === "encrypted" &&
    returnValueType === "encrypted"
  ) {
    return testContract1Arg(functionName);
  }

  if (
    inputCount === 2 &&
    inputs[0] === "encrypted" &&
    inputs[1] === "encrypted"
  ) {
    if (returnValueType === "ebool") {
      return testContract2ArgBoolRes(functionName, isBooleanMathOp);
    }
    return testContract2Arg(
      functionName,
      isBooleanMathOp,
      getOperator(functionName)
    );
  }

  if (inputCount === 3) {
    return testContract3Arg(functionName);
  }

  console.log(
    `Function ${functionName} with ${inputCount} inputs that are ${inputs} is not implemented`
  );

  return ["", ""];
};

/** Generates a Solidity bench contract based on the provided metadata */
const generateSolidityBenchContract = (metadata: FunctionMetadata): string => {
  const { functionName, inputCount, inputs } = metadata;

  if (functionName === "random" || functionName === "decrypt") {
    // todo: bench random/decrypt function
    return "";
  }

  if (functionName === SEALING_FUNCTION_NAME) {
    return benchContractReencrypt();
  }

  if (functionName === SEALING_TYPED_FUNCTION_NAME) {
    // `sealoutputTyped` is a wrapper around `sealoutput`, and does not need to be benchmarked directly
    return "";
  }

  if (
    inputCount === 2 &&
    inputs[0] === "encrypted" &&
    inputs[1] === "encrypted"
  ) {
    return benchContract2Arg(functionName);
  }

  if (inputCount === 1) {
    return benchContract1Arg(functionName);
  }

  if (inputCount === 3) {
    return benchContract3Arg(functionName);
  }

  console.log(
    `Function ${functionName} with ${inputCount} inputs that are ${inputs} is not implemented`
  );

  return "";
};

/**
 * Generates a Solidity function based on the provided metadata
 * This generates all the different types of function headers that can exist
 */
const genSolidityFunctionHeaders = (metadata: FunctionMetadata): string[] => {
  const {
    functionName,
    inputCount,
    hasDifferentInputTypes,
    returnValueType,
    inputs,
    isBooleanMathOp,
  } = metadata;

  let functions: string[][] = [];

  inputs.forEach((input, idx) => {
    let inputVariants = [];
    switch (input) {
      case "encrypted":
        let index = 0;
        for (let inputType of EInputType) {
          if (!IsOperationAllowed(functionName, index++)) {
            // Skip unallowed operations based on FheEngine's operation_is_allowed
            continue;
          }
          if (
            inputs.length === 2 &&
            !isBooleanMathOp &&
            isComparisonType(inputType)
          ) {
            continue;
          }
          inputVariants.push(`input${idx} ${inputType}`);
        }
        break;
      case "plaintext":
        for (let inputType of EPlaintextType) {
          inputVariants.push(`input${idx} ${inputType}`);
        }
        break;
      default:
        inputVariants.push(`input${idx} ${input}`);
    }
    functions.push(inputVariants);
  });

  return getAllFunctionDeclarations(
    functionName,
    functions,
    isBooleanMathOp,
    returnValueType
  );
};

type ParsedFunction = {
  funcName: string;
  inputs: AllTypes[];
  returnType: AllTypes;
  inputPlaintext: string;
};

// Regular expression to match the Solidity function signature pattern
const functionPattern =
  /function (\w+)\((.*?)\) internal pure returns \((.*?)\);/;

/**
 * Parses a Solidity function definition into its components.
 *
 * @param funcDef Solidity function definition as string.
 * @returns An object containing the functionName, inputTypes, and outputType.
 */
const parseFunctionDefinition = (funcDef: string): ParsedFunction => {
  const match = funcDef.match(functionPattern);

  if (!match) {
    throw new Error(`Invalid function definition format for ${funcDef}`);
  }

  const [, functionName, inputs, outputType] = match;

  const inputTypes = inputs.split(",").map((input) => {
    return input.trim().split(/\s+/).pop();
  });

  const inputPlaintext = inputTypes[0]!.startsWith("e") ? "none" : "all";

  return <ParsedFunction>{
    funcName: functionName,
    inputs: inputTypes,
    returnType: outputType,
    inputPlaintext,
  };
};

// Helper function to capitalize type name for asEuintX function call.

// This will generate the Solidity function body based on the function definition provided.
const generateSolidityFunction = (parsedFunction: ParsedFunction): string => {
  const { funcName, inputs, returnType } = parsedFunction;
  switch (inputs.length) {
    case 1:
      if (funcName === "decrypt") {
        return SolTemplateDecrypt(inputs[0], returnType);
      }
      return SolTemplate1Arg(funcName, inputs[0], returnType);
    case 2:
      return SolTemplate2Arg(funcName, inputs[0], inputs[1], returnType);
    case 3:
      return SolTemplate3Arg(
        funcName,
        inputs[0],
        inputs[1],
        inputs[2],
        returnType
      );
    default:
      throw new Error("Unknown number of inputs");
  }
};

const main = async () => {
  let metadata = await generateMetadataPayload();
  let solidityHeaders: string[] = [];
  const testContracts: Record<string, string> = {};
  const benchContracts: Record<string, string> = {};
  let testContractsAbis = "";
  let importLineHelper: string = "import { ";

  for (let func of metadata) {
    // this generates test contract for every function
    const testContract = generateSolidityTestContract(func);
    const benchContract = generateSolidityBenchContract(func);

    if (testContract[0] !== "") {
      testContracts[capitalize(func.functionName)] = testContract[0];
      testContractsAbis += testContract[1];
      importLineHelper += `${capitalize(func.functionName)}TestType,\n`;
    }

    if (benchContract !== "") {
      benchContracts[capitalize(func.functionName)] = benchContract;
    }
    // this generates solidity header functions for all the different possible types
    solidityHeaders = solidityHeaders.concat(genSolidityFunctionHeaders(func));
  }

  //console.log(solidityHeaders.filter(name => name.includes('cmux')).map(item => parseFunctionDefinition(item)));

  let outputFile = preamble();
  for (let fn of solidityHeaders) {
    const funcDefinition = generateSolidityFunction(
      parseFunctionDefinition(fn)
    );
    outputFile += funcDefinition;
  }
  outputFile += RandomGenericFunction();
  outputFile += RandomFunctions();
  outputFile += `\n\n    // ********** TYPE CASTING ************* //`;

  // generate casting functions
  for (let fromType of EInputType.concat("uint256", "bytes memory")) {
    for (let toType of EInputType) {
      if (fromType === toType) {
        // todo: this is a bit weird, but I'm using this place to create asXXX functions for the cast from the input types (inXXX)
        const inputTypeName = `in${capitalize(fromType)}`;
        outputFile += AsTypeFunction(inputTypeName, toType);
        continue;
      }

      outputFile += AsTypeFunction(fromType, toType);
    }
  }
  // For a better UX, allow casting from address to eaddress:
  outputFile += AsTypeFunction("address", "eaddress");

  for (let type of EInputType) {
    const functionName = `as${capitalize(type)}`;
    const testContract = AsTypeTestingContract(type);

    testContracts[functionName] = testContract[0];
    benchContracts[functionName] = AsTypeBenchmarkContract(type);

    testContractsAbis += testContract[1];
    importLineHelper += `${capitalize(functionName)}TestType,\n`;
  }

  importLineHelper = importLineHelper.slice(0, -2) + " } from './abis';\n";

  outputFile += AsTypeFunction("bool", "ebool");

  outputFile += PostFix();

  outputFile += `\n\n// ********** OPERATOR OVERLOADING ************* //\n`;

  // generate operator overloading
  ShorthandOperations.filter((v) => v.operator !== null).forEach((value) => {
    let idx = 0;
    for (let encType of EInputType) {
      if (!valueIsEncrypted(encType)) {
        throw new Error("InputType mismatch");
      }

      if (!IsOperationAllowed(value.func, idx++)) {
        // Skip unallowed operations based on FheEngine's operation_is_allowed
        continue;
      }
      if (!isComparisonType(encType) || isBitwiseOp(value.func)) {
        outputFile += OperatorOverloadDecl(
          value.func,
          value.operator!,
          encType,
          value.unary,
          value.returnsBool
        );
      }
    }
  });

  outputFile += `\n// ********** BINDING DEFS ************* //`;

  EInputType.forEach((encryptedType) => {
    outputFile += BindingLibraryType(encryptedType);

    BindMathOperators.forEach((fnToBind) => {
      let foundFnDef = solidityHeaders.find((funcHeader) => {
        const fnDef = parseFunctionDefinition(funcHeader);
        const input = fnDef.inputs[0];

        if (!EInputType.includes(input)) {
          return false;
        }

        if (
          !IsOperationAllowed(fnDef.funcName, EInputType.indexOf(encryptedType))
        ) {
          return false;
        }
        return (
          fnDef.funcName === fnToBind &&
          fnDef.inputs.every((item) => item === input)
        );
      });

      if (foundFnDef) {
        const fnDef = parseFunctionDefinition(foundFnDef);
        if (
          !isComparisonType(encryptedType) ||
          fnDef.inputs.every(isComparisonType)
        ) {
          outputFile += OperatorBinding(
            fnDef.funcName,
            encryptedType,
            fnDef.inputs.length === 1,
            fnDef.returnType === "ebool" &&
              !bitwiseAndLogicalOperators.includes(fnDef.funcName)
          );
        }
      }
    });

    EInputType.filter((otherType) => otherType !== encryptedType).forEach(
      (otherType) => {
        outputFile += CastBinding(encryptedType, otherType);
      }
    );

    outputFile += SealFromType(encryptedType);
    outputFile += SealTypedFromType(encryptedType);
    outputFile += DecryptBinding(encryptedType);

    outputFile += PostFix();
  });

  await fs.promises.writeFile("FHE.sol", outputFile);
  for (const testContract of Object.entries(testContracts)) {
    fs.writeFileSync(
      `../solidity/tests/contracts/${testContract[0]}.sol`,
      testContract[1]
    );
  }

  for (const benchContract of Object.entries(benchContracts)) {
    fs.writeFileSync(
      `../solidity/tests/contracts/bench/${benchContract[0]}.sol`,
      benchContract[1]
    );
  }

  fs.writeFileSync("../solidity/tests/abis.ts", genAbiFile(testContractsAbis));
  console.log(importLineHelper);
};

main();
