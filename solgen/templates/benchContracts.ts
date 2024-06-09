import {
  EInputType,
  SEALING_FUNCTION_NAME,
  AllowedOperations,
  capitalize,
} from "../common";

const toVarSuffix = (inputType: string) => capitalize(inputType.slice(1).replace("uint", ""));
const toInType = (inputType: string) => "inE" + inputType.slice(1);
const toInTypeFunction = (inputType: string) => toInType(inputType) + " calldata";
const toAsType = (inputType: string) => "asE" + inputType.slice(1);

const IsOperationAllowed = (
  functionName: string,
  dataType: string,
): boolean => {
  const inputIdx= EInputType.indexOf(dataType)
  const regexes = AllowedOperations[inputIdx];
  for (let regex of regexes) {
    if (!new RegExp(regex).test(functionName.toLowerCase())) {
      return false;
    }
  }

  return true;
};

export function generateBenchContract(
  name: string,
  testFunc: string,
  importStatement: string = `import {ebool, euint8} from "../../../FHE.sol";`,
) {
  return `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
${importStatement}

contract ${capitalize(name)}Bench {
${testFunc}
}
`;
}

export function benchContract1Arg(name: string) {
  let importTypes = ``;
  let privateVarsA = ``;
  let funcLoad = "";
  let funcBench = "";

  for (let inputType of EInputType) {
    if (IsOperationAllowed(name, inputType)) {
      importTypes += "\n\t" + inputType + ", " + toInType(inputType) + ",";
      privateVarsA += `\t${inputType} internal a${toVarSuffix(inputType)};\n`;
      funcLoad += `
    function load${toVarSuffix(inputType)}(${toInTypeFunction(inputType)} _a) public {
        a${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_a);
    }`;

      // todo: should this return something? should we verify the decrypted result of the operation?
      funcBench += `
    function bench${capitalize(name)}${toVarSuffix(inputType)}() public view {
        FHE.${name}(a${toVarSuffix(inputType)});
    }`;
    }
  }

  const func = privateVarsA + funcLoad + "\n" + funcBench;
  importTypes = importTypes.slice(0, -1); // remove last comma
  const importStatement = `import {${importTypes}
} from "../../../FHE.sol";`;

  return [generateBenchContract(name, func, importStatement)];
}

export function benchContract2Arg(name: string) {
  let importTypes = ``;
  let privateVarsA = ``;
  let privateVarsB = ``;
  let funcLoad = "";
  let funcBench = "";

  for (let inputType of EInputType) {
    if (IsOperationAllowed(name, inputType)) {
      importTypes += "\n\t" + inputType + ", " + toInType(inputType) + ",";
      privateVarsA += `\t${inputType} internal a${toVarSuffix(inputType)};\n`;
      privateVarsB += `\t${inputType} internal b${toVarSuffix(inputType)};\n`;
      funcLoad += `
    function load${toVarSuffix(inputType)}(${toInTypeFunction(inputType)} _a, ${toInTypeFunction(inputType)} _b) public {
        a${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_a);
        b${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_b);
    }`;

      // todo: should this return something? should we verify the decrypted result of the operation?
      funcBench += `
    function bench${capitalize(name)}${toVarSuffix(inputType)}() public view {
        FHE.${name}(a${toVarSuffix(inputType)}, b${toVarSuffix(inputType)});
    }`;
    }
  }

  const func = privateVarsA + "\n" + privateVarsB + funcLoad + "\n" + funcBench;
  importTypes = importTypes.slice(0, -1); // remove last comma
  const importStatement = `import {${importTypes}
} from "../../../FHE.sol";`;

  return [generateBenchContract(name, func, importStatement)];
}

export function benchContract3Arg(name: string) {
  let importTypes = ``;
  let privateVarsA = `\tebool internal control;\n\n`;
  let privateVarsB = ``;
  let funcLoad = "";
  let funcBench = "";

  for (let inputType of EInputType) {
    if (IsOperationAllowed(name, inputType)) {
      importTypes += "\n\t" + inputType + ", " + toInType(inputType) + ",";
      privateVarsA += `\t${inputType} internal a${toVarSuffix(inputType)};\n`;
      privateVarsB += `\t${inputType} internal b${toVarSuffix(inputType)};\n`;
      funcLoad += `
    function load${toVarSuffix(inputType)}(inEbool calldata _control, ${toInTypeFunction(inputType)} _a, ${toInTypeFunction(inputType)} _b) public {
        control = FHE.asEbool(_control);
        a${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_a);
        b${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_b);
    }`;

      // todo: should this return something? should we verify the decrypted result of the operation?
      funcBench += `
    function bench${capitalize(name)}${toVarSuffix(inputType)}() public view {
        FHE.${name}(control, a${toVarSuffix(inputType)}, b${toVarSuffix(inputType)});
    }`;
    }
  }

  const func = privateVarsA + "\n" + privateVarsB + funcLoad + "\n" + funcBench;
  importTypes = importTypes.slice(0, -1); // remove last comma
  const importStatement = `import {${importTypes}
} from "../../../FHE.sol";`;

  return [generateBenchContract(name, func, importStatement)];
}

export function benchContractReencrypt() {
  let importTypes = ``;
  let privateVarsA = `\tbytes32 internal pubkey;\n`;
  let funcLoad = "";
  let funcBench = "";

  for (let inputType of EInputType) {
    if (IsOperationAllowed(SEALING_FUNCTION_NAME, inputType)) {
      importTypes += "\n\t" + inputType + ", " + toInType(inputType) + ",";
      privateVarsA += `\t${inputType} internal a${toVarSuffix(inputType)};\n`;
      funcLoad += `
    function load${toVarSuffix(inputType)}(${toInTypeFunction(inputType)} _a, bytes32 _pubkey) public {
        a${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_a);
        pubkey = _pubkey;
    }`;

      // todo: should this return something? should we verify the decrypted result of the operation?
      funcBench += `
    function bench${capitalize(SEALING_FUNCTION_NAME)}${toVarSuffix(inputType)}() public view {
        FHE.${SEALING_FUNCTION_NAME}(a${toVarSuffix(inputType)}, pubkey);
    }`;
    }
  }

  const func = privateVarsA + funcLoad + "\n" + funcBench;
  importTypes = importTypes.slice(0, -1); // remove last comma
  const importStatement = `import {${importTypes}
} from "../../../FHE.sol";`;

  return [generateBenchContract(SEALING_FUNCTION_NAME, func, importStatement)];
}

export function AsTypeBenchmarkContract(type: string) {
  let funcs = "";
  // Although casts from eaddress to types with < 256 bits are possible, we don't want to bench them.
  let eaddressAllowedTypes = ["euint256", "uint256", "bytes memory"];
  let fromTypeCollection = type === "eaddress" ? eaddressAllowedTypes : EInputType.concat("uint256", "bytes memory");

  for (const fromType of fromTypeCollection) {
    if (type === fromType || (fromType === "eaddress" && !eaddressAllowedTypes.includes(type))) {
      continue;
    }

    const fromTypeTs = fromType === "bytes memory" ? "Uint8Array" : `bigint`;
    const fromTypeSol = fromType === "bytes memory" ? fromType : `uint256`;
    const fromTypeEncrypted = EInputType.includes(fromType)
      ? fromType
      : undefined;
    const contractInfo = TypeCastBenchmarkFunction(
      fromTypeSol,
      fromTypeTs,
      type,
      fromTypeEncrypted
    );
    funcs += contractInfo[0];
    abi += contractInfo[1];
  }

  funcs = funcs.slice(1);
  abi += `}\n`;

  return [generateBenchContract(`As${capitalize(type)}`, funcs), abi];
}

function TypeCastBenchmarkFunction(
  fromType: string,
  fromTypeForTs: string,
  toType: string,
  fromTypeEncrypted?: string
) {
  let to = capitalize(toType);
  const retType = to.slice(1);
  let testType = fromTypeEncrypted ? fromTypeEncrypted : fromType;
  testType =
    testType === "bytes memory" ? "PreEncrypted" : capitalize(testType);
  testType = testType === "Uint256" ? "Plaintext" : testType;
  const encryptedVal = fromTypeEncrypted
    ? `FHE.as${capitalize(fromTypeEncrypted)}(val)`
    : "val";
  let retTypeTs = retType === "bool" ? "boolean" : retType;
  retTypeTs = retTypeTs.includes("uint") || retTypeTs.includes("address") ? "bigint" : retTypeTs;

  let abi: string;
  let func = "\n\n    ";

  if (testType === "PreEncrypted" || testType === "Plaintext") {
    func += `function castFrom${testType}To${to}(${fromType} val) public pure returns (${retType}) {
        return FHE.decrypt(FHE.as${to}(${encryptedVal}));
    }`;
    abi = `    castFrom${testType}To${to}: (val: ${fromTypeForTs}) => Promise<${retTypeTs}>;\n`;
  } else {
    func += `function castFrom${testType}To${to}(${fromType} val, string calldata test) public pure returns (${retType}) {
        if (Utils.cmp(test, "bound")) {
            return ${encryptedVal}.to${shortenType(toType)}().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.as${to}(${encryptedVal}));
        }
        revert TestNotFound(test);
    }`;
    abi = `    castFrom${testType}To${to}: (val: ${fromTypeForTs}, test: string) => Promise<${retTypeTs}>;\n`;
  }

  return [func, abi];
}
