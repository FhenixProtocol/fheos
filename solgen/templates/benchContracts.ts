import {
  EInputType,
  SEALING_FUNCTION_NAME,
  SEAL_RETURN_TYPE,
  LOCAL_SEAL_FUNCTION_NAME,
  AllowedOperations,
  capitalize,
} from "../common";

export const toVarSuffix = (inputType: string) => capitalize(inputType.slice(1).replace("uint", ""));
export const toInType = (inputType: string) => "inE" + inputType.slice(1);
export const toInTypeFunction = (inputType: string) => toInType(inputType) + " calldata";
export const toAsType = (inputType: string) => "asE" + inputType.slice(1);

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

  // todo: verify that the ts input should be bytes for inEuints
  // todo: add all abi functions
  const abi = `export interface ${capitalize(
    name
  )}BenchType extends BaseContract {
    ${name}: (_a: bytes, _b: bytes) => Promise<bigint>;
}\n`;
  return [generateBenchContract(name, func, importStatement), abi];
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

  // todo: verify that the ts input should be bytes for inEuints
  // todo: add all abi functions
  const abi = `export interface ${capitalize(
    name
  )}BenchType extends BaseContract {
    ${name}: (_a: bytes, _b: bytes) => Promise<bigint>;
}\n`;
  return [generateBenchContract(name, func, importStatement), abi];
}

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

export function benchContractReencrypt() {
  // let func = `function ${SEALING_FUNCTION_NAME}(string calldata test, uint256 a, bytes32 pubkey) public pure returns (${SEAL_RETURN_TYPE} memory reencrypted) {
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

  // todo: verify that the ts input should be bytes for inEuints
  // todo: add all abi functions
  const abi = `export interface SealoutputBenchType extends BaseContract {
    ${SEALING_FUNCTION_NAME}: (test: string, a: bigint, pubkey: Uint8Array) => Promise<string>;`;

  return [generateBenchContract(SEALING_FUNCTION_NAME, func, importStatement), abi];
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

  // todo: verify that the ts input should be bytes for inEuints
  // todo: add all abi functions
  const abi = `export interface ${capitalize(
    name
  )}BenchType extends BaseContract {
    ${name}: (_a: bytes, _b: bytes) => Promise<bigint>;
}\n`;
  return [generateBenchContract(name, func, importStatement), abi];
}

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