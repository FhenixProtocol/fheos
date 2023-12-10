import {getFunctionsFromGo} from "./contracts_parser";
import * as fs from 'fs';
import {
    AsTypeFunction,
    BindingLibraryType, BindingsWithoutOperator,
    OperatorBinding,
    OperatorOverloadDecl,
    PostFix,
    preamble,
    SolTemplate1Arg,
    SolTemplate2Arg,
    SolTemplate3Arg
} from "./templates";
import {AllTypes, BindMathOperators, EInputType, EPlaintextType, ShorthandOperations, valueIsEncrypted} from "./common";

interface FunctionMetadata {
    functionName: string;
    inputCount: number;
    hasDifferentInputTypes: boolean;
    returnValueType?: string;
    inputs: AllTypes[];
}

const generateMetadataPayload = async (): Promise<FunctionMetadata[]> => {
    let result = await getFunctionsFromGo('../precompiles/contracts.go');

    return result.map((value) => {
        return {
            functionName: value.name,
            hasDifferentInputTypes: !value.needsSameType,
            inputCount: value.paramsCount,
            returnValueType: value.returnType,
            inputs: value.inputTypes,
        }
    })
}

// Function to generate all combinations of parameters.
function generateCombinations(arr: string[][], current: string[] = [], index: number = 0): string[][] {
    if (index === arr.length) {
        // Only add to the result if there are elements in the current combination
        return current.length === 0 ? [] : [current];
    }

    let result: string[][] = [];
    // Add each element of the current sub-array to the combination
    for (let item of arr[index]) {
        result = result.concat(generateCombinations(arr, current.concat(item), index + 1));
    }
    return result;
}

const getReturnType = (inputs: string[], returnType?: string) => {
    if (returnType === 'plaintext') {
        if (inputs.length != 1) {
            throw new Error("expecting exactly one input for functions returning plaintext");
        }

        let inputType = inputs[0].split(' ')[1];
        if (inputType[0] !== 'e') {
            throw new Error("expecting encrypted input for plaintext output");
        }

        return inputType.slice(1);
    }


    if (returnType && returnType !== "encrypted") {
        return returnType;
    }

    if (inputs.includes("bytes") || inputs.includes("bytes32")) {
        return "bytes";
    }

    let maxRank = 0;
    for (let input of inputs) {
        let inputType = input.split(' ')[1];
        maxRank = Math.max(EInputType.indexOf(inputType), EPlaintextType.indexOf(inputType), maxRank);
    }

    return EInputType[maxRank]
}

function getAllFunctionDeclarations(functionName: string, functions: string[][], returnValueType?: string): string[] {
    let functionDecl = `function ${functionName}`;

    // Generate all combinations of input parameters.
    let allCombinations = generateCombinations(functions);

    // Create function declarations for each combination.
    return allCombinations.map(combination => {
        let returnType = getReturnType(combination, returnValueType);
        let returnStr =  `internal pure returns (${returnType});`;

        return `${functionDecl}(${combination.join(', ')}) ${returnStr}`;
    });
}

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
        inputs
    } = metadata;

    let functions: string[][] = [];

    inputs.forEach((input, idx) => {
        let inputVariants = [];
        switch (input) {
            case "encrypted":
                for (let inputType of EInputType) {
                    inputVariants.push(`input${idx} ${inputType}`)
                }
                break;
            case "plaintext":
                for (let inputType of EPlaintextType) {
                    inputVariants.push(`input${idx} ${inputType}`)
                }
                break;
            default:
                inputVariants.push(`input${idx} ${input}`)
        }
        functions.push(inputVariants);
    });

    return getAllFunctionDeclarations(functionName, functions, returnValueType);
};

type ParsedFunction = {
    funcName: string;
    inputs: AllTypes[];
    returnType: AllTypes;
    inputPlaintext: string;
};

// Regular expression to match the Solidity function signature pattern
const functionPattern = /function (\w+)\((.*?)\) internal pure returns \((.*?)\);/;

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
    const inputTypes = inputs.split(',').map(input => {
        return input.trim().split(/\s+/).pop();
    });

    const inputPlaintext = inputTypes[0]!.startsWith('e') ? "none" : "all";

    return <ParsedFunction>{
        funcName: functionName,
        inputs: inputTypes,
        returnType: outputType,
        inputPlaintext
    };
};


// Helper function to capitalize type name for asEuintX function call.

// This will generate the Solidity function body based on the function definition provided.
const generateSolidityFunction = (
    parsedFunction: ParsedFunction,
): string => {

    const {funcName, inputs, returnType} = parsedFunction;

    switch (inputs.length) {
        case 1:
            return SolTemplate1Arg(funcName, inputs[0], returnType);
        case 2:
            return SolTemplate2Arg(funcName, inputs[0], inputs[1], returnType);
        case 3:
            return SolTemplate3Arg(funcName, inputs[0], inputs[1], inputs[2], returnType);
        default:
            throw new Error("Unknown number of inputs");
    }
}

const main = async () => {

    let metadata = await generateMetadataPayload();
    let solidityHeaders: string[] = [];
    for (let func of metadata) {
        // this generates solidity header functions for all the different possible types
        solidityHeaders = solidityHeaders.concat(genSolidityFunctionHeaders(func));
    }

    //console.log(solidityHeaders.filter(name => name.includes('cmux')).map(item => parseFunctionDefinition(item)));

    let outputFile = preamble();
    for (let fn of solidityHeaders) {
        // this generates the function body from the header
        const funcDefinition = generateSolidityFunction(parseFunctionDefinition(fn));
        outputFile += funcDefinition;
    }
    outputFile += `\n// ********** TYPE CASTING ************* //\n`

    // generate casting functions
    for (let fromType of EInputType.concat('uint256', 'bytes memory')) {
        for (let toType of EInputType) {

            // casting from bool is annoying - if we really need this we can add it later
            if (fromType === "bool") {
                continue;
            }
            outputFile += AsTypeFunction(fromType, toType);
        }
    }

    outputFile += PostFix();

    outputFile += `\n// ********** OPERATOR OVERLOADING ************* //\n`

    // generate operator overloading
    ShorthandOperations.forEach((value) =>  {
        for (let encType of EInputType) {
            if (!valueIsEncrypted(encType)) {
                throw new Error("InputType mismatch");
            }
            outputFile += OperatorOverloadDecl(value.func, value.operator, encType, value.unary)
        }
    });

    outputFile += `\n// ********** BINDING DEFS ************* //\n`

    EInputType.forEach(encryptedType => {

        BindMathOperators.forEach(bindMathOp => {

            if (ShorthandOperations.filter(value => value.func === bindMathOp).length === 0) {
                // console.log(`${bindMathOp}`)
                outputFile += BindingsWithoutOperator(bindMathOp, encryptedType);
            }
        });

        outputFile += BindingLibraryType(encryptedType);
        BindMathOperators.forEach(fnToBind => {
            let foundFnDef = solidityHeaders.find((funcHeader) => {
                const fnDef = parseFunctionDefinition(funcHeader);
                const input = fnDef.inputs[0];

                if (!EInputType.includes(input)) {
                    return false;
                }

                return (fnDef.funcName === fnToBind && fnDef.inputs.every(item => item === input))
            });

            if (foundFnDef) {
                const fnDef = parseFunctionDefinition(foundFnDef);
                outputFile += OperatorBinding(fnDef.funcName, encryptedType, fnDef.inputs.length === 1)
            }
        });
        outputFile += PostFix();
    })



    await fs.promises.writeFile('FHE.sol', outputFile);
}

main();


