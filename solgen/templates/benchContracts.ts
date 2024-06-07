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
    function load${toVarSuffix(inputType)}(${toInType(inputType)} _a, ${toInType(inputType)} _b) public {
        a32 = FHE.${toAsType(inputType)}(_a);
        b32 = FHE.${toAsType(inputType)}(_b);
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
    function load${toVarSuffix(inputType)}(${toInType(inputType)} _a) public {
        a32 = FHE.${toAsType(inputType)}(_a);
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

export function benchContractReq() {
  // Req is failing on EthCall so we need to make it as tx for now
  // todo: check the claim that req failing on EthCall
  let func = `
    euint8 internal a8;
    euint16 internal a16;
    euint32 internal a32;
    euint64 internal a64;
    euint128 internal a128;
    euint256 internal a256;
  
    function load32(inEuint32 _a) public {
        a32 = FHE.asEuint32(_a);
    }
    
    function benchReq32() public view {
        FHE.req(a32);
    }`;
  // todo verify that the ts input should be bytes for inEuints
  const abi = `export interface ReqBenchType extends BaseContract {
    req: (a: bytes) => Promise<{}>;
}\n`;
  return [generateBenchContract("req", func), abi];
}

export function benchContractReencrypt() {
  let func = `function ${SEALING_FUNCTION_NAME}(string calldata test, uint256 a, bytes32 pubkey) public pure returns (${SEAL_RETURN_TYPE} memory reencrypted) {
        if (Utils.cmp(test, "${SEALING_FUNCTION_NAME}(euint8)")) {
            return FHE.${SEALING_FUNCTION_NAME}(FHE.asEuint8(a), pubkey);
        } else if (Utils.cmp(test, "${SEALING_FUNCTION_NAME}(euint16)")) {
            return FHE.${SEALING_FUNCTION_NAME}(FHE.asEuint16(a), pubkey);
        } else if (Utils.cmp(test, "${SEALING_FUNCTION_NAME}(euint32)")) {
            return FHE.${SEALING_FUNCTION_NAME}(FHE.asEuint32(a), pubkey);
        } else if (Utils.cmp(test, "${SEALING_FUNCTION_NAME}(euint64)")) {
            return FHE.${SEALING_FUNCTION_NAME}(FHE.asEuint64(a), pubkey);
        } else if (Utils.cmp(test, "${SEALING_FUNCTION_NAME}(euint128)")) {
            return FHE.${SEALING_FUNCTION_NAME}(FHE.asEuint128(a), pubkey);
        } else if (Utils.cmp(test, "${SEALING_FUNCTION_NAME}(euint256)")) {
            return FHE.${SEALING_FUNCTION_NAME}(FHE.asEuint256(a), pubkey);
        } else if (Utils.cmp(test, "${SEALING_FUNCTION_NAME}(ebool)")) {
            bool b = true;
            if (a == 0) {
                b = false;
            }

            return FHE.${SEALING_FUNCTION_NAME}(FHE.asEbool(b), pubkey);
        } else if (Utils.cmp(test, "${LOCAL_SEAL_FUNCTION_NAME}(euint8)")) {
            euint8 aEnc = FHE.asEuint8(a);
            return aEnc.${LOCAL_SEAL_FUNCTION_NAME}(pubkey);
        }
        revert TestNotFound(test);
    }`;
  const abi = `export interface SealoutputBenchType extends BaseContract {
    ${SEALING_FUNCTION_NAME}: (test: string, a: bigint, pubkey: Uint8Array) => Promise<string>;
}\n`;
  return [generateBenchContract(SEALING_FUNCTION_NAME, func), abi];
}

export function benchContract3Arg(name: string) {
  let func = `function ${name}(string calldata test, bool c, uint256 a, uint256 b) public pure returns (uint256 output) {
        ebool condition = FHE.asEbool(c);
        if (Utils.cmp(test, "${name}: euint8")) {
            return FHE.decrypt(FHE.${name}(condition, FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "${name}: euint16")) {
            return FHE.decrypt(FHE.${name}(condition, FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "${name}: euint32")) {
            return FHE.decrypt(FHE.${name}(condition, FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "${name}: euint64")) {
            return FHE.decrypt(FHE.${name}(condition, FHE.asEuint64(a), FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "${name}: euint128")) {
            return FHE.decrypt(FHE.${name}(condition, FHE.asEuint128(a), FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "${name}: ebool")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }

            if(FHE.decrypt(FHE.${name}(condition, FHE.asEbool(aBool), FHE.asEbool(bBool)))) {
                return 1;
            }
            return 0;
        } 
        
        revert TestNotFound(test);
    }`;
  const abi = `export interface ${capitalize(
    name
  )}BenchType extends BaseContract {
    ${name}: (test: string, c: boolean, a: bigint, b: bigint) => Promise<bigint>;
}\n`;
  return [generateBenchContract(name, func), abi];
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