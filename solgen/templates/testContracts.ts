import {
  EInputType,
  SEALING_FUNCTION_NAME,
  SEAL_RETURN_TYPE,
  LOCAL_SEAL_FUNCTION_NAME,
  AllowedOperations,
  AllowedTypesOnCastToEaddress,
  capitalize,
  shortenType,
  toInType,
  toInTypeParam,
} from "../common";

function TypeCastTestingFunction(
  fromType: string,
  fromTypeForTs: string,
  toType: string,
  fromTypeEncrypted?: string
) {
  let to = capitalize(toType);
  const retType = to.slice(1);

  let testType = fromTypeEncrypted ? fromTypeEncrypted : fromType;
  testType = testType === "uint256" ? "Plaintext" : testType;
  testType = testType === "address" ? "PlaintextAddress" : testType;
  testType = testType.startsWith("inE") ? "PreEncrypted" : testType;
  testType = capitalize(testType)

  const encryptedVal = fromTypeEncrypted
    ? `FHE.as${capitalize(fromTypeEncrypted)}(val)`
    : "val";
  let retTypeTs = retType === "bool" ? "boolean" : retType;
  retTypeTs = retTypeTs.includes("uint") || retTypeTs.includes("address") ? "bigint" : retTypeTs;

  let abi: string;
  let func = "\n\n    ";

  if (testType === "PreEncrypted" || testType === "Plaintext" || testType === "PlaintextAddress") {
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

export function AsTypeTestingContract(type: string) {
  let funcs = "";
  let abi = `export interface As${capitalize(
    type
  )}TestType extends BaseContract {\n`;

  let typesToEaddres = AllowedTypesOnCastToEaddress
    .filter(t => t !== "inEaddress")    // added explicitly later
    .filter(t => t !== "bytes memory"); // tested indirectly via "inEaddress"
  let fromTypeCollection = type === "eaddress" ? typesToEaddres : EInputType.concat("uint256");

  // add inE(type) calldata
  fromTypeCollection = fromTypeCollection.concat(toInTypeParam(type));

  for (const fromType of fromTypeCollection) {
    if (type === fromType || (fromType === "eaddress" && !AllowedTypesOnCastToEaddress.includes(type))) {
      continue;
    }

    const fromTypeTs = fromType.startsWith("inE") ? "EncryptedNumber" : `bigint`;
    const fromTypeSol = fromType.startsWith("e") ? `uint256` : fromType;
    const fromTypeEncrypted = EInputType.includes(fromType)
      ? fromType
      : undefined;
    const contractInfo = TypeCastTestingFunction(
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

  return [generateTestContract(`As${capitalize(type)}`, funcs, [toInType(type)]), abi];
}

export function testContract2ArgBoolRes(name: string, isBoolean: boolean) {
  const isEuint64Allowed = IsOperationAllowed(
    name,
    EInputType.indexOf("euint64")
  );
  const isEuint128Allowed = IsOperationAllowed(
    name,
    EInputType.indexOf("euint128")
  );
  const isEuint256Allowed = IsOperationAllowed(
    name,
    EInputType.indexOf("euint256")
  );
  const isEaddressAllowed = IsOperationAllowed(
    name,
    EInputType.indexOf("eaddress")
  )
  let func = `function ${name}(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "${name}(euint8,euint8)")) {
            if (FHE.decrypt(FHE.${name}(FHE.asEuint8(a), FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "${name}(euint16,euint16)")) {
            if (FHE.decrypt(FHE.${name}(FHE.asEuint16(a), FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "${name}(euint32,euint32)")) {
            if (FHE.decrypt(FHE.${name}(FHE.asEuint32(a), FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        }`;
  if (isEuint64Allowed) {
    func += ` else if (Utils.cmp(test, "${name}(euint64,euint64)")) {
            if (FHE.decrypt(FHE.${name}(FHE.asEuint64(a), FHE.asEuint64(b)))) {
                return 1;
            }

            return 0;
        }`;
  }
  if (isEuint128Allowed) {
    func += ` else if (Utils.cmp(test, "${name}(euint128,euint128)")) {
            if (FHE.decrypt(FHE.${name}(FHE.asEuint128(a), FHE.asEuint128(b)))) {
                return 1;
            }

            return 0;
        }`;
  }
  if (isEuint256Allowed) {
    func += ` else if (Utils.cmp(test, "${name}(euint256,euint256)")) {
            if (FHE.decrypt(FHE.${name}(FHE.asEuint256(a), FHE.asEuint256(b)))) {
                return 1;
            }

            return 0;
        }`;
  }
  if (isEaddressAllowed) {
    func += ` else if (Utils.cmp(test, "${name}(eaddress,eaddress)")) {
            if (FHE.decrypt(FHE.${name}(FHE.asEaddress(a), FHE.asEaddress(b)))) {
                return 1;
            }

            return 0;
        }`;
  }
  func += ` else if (Utils.cmp(test, "euint8.${name}(euint8)")) {
            if (FHE.asEuint8(a).${name}(FHE.asEuint8(b)).decrypt()) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.${name}(euint16)")) {
            if (FHE.asEuint16(a).${name}(FHE.asEuint16(b)).decrypt()) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.${name}(euint32)")) {
            if (FHE.asEuint32(a).${name}(FHE.asEuint32(b)).decrypt()) {
                return 1;
            }

            return 0;
        }`;

  if (isEuint64Allowed) {
    func += ` else if (Utils.cmp(test, "euint64.${name}(euint64)")) {
            if (FHE.asEuint64(a).${name}(FHE.asEuint64(b)).decrypt()) {
                return 1;
            }
            return 0;
        }`;
  }
  if (isEuint128Allowed) {
    func += ` else if (Utils.cmp(test, "euint128.${name}(euint128)")) {
            if (FHE.asEuint128(a).${name}(FHE.asEuint128(b)).decrypt()) {
                return 1;
            }
            return 0;
        }`;
  }
  if (isEuint256Allowed) {
    func += ` else if (Utils.cmp(test, "euint256.${name}(euint256)")) {
            if (FHE.asEuint256(a).${name}(FHE.asEuint256(b)).decrypt()) {
                return 1;
            }
            return 0;
        }`;
  }
  if (isEaddressAllowed) {
    func += ` else if (Utils.cmp(test, "eaddress.${name}(eaddress)")) {
            if (FHE.asEaddress(a).${name}(FHE.asEaddress(b)).decrypt()) {
                return 1;
            }
            return 0;
        }`;
  }
  if (isBoolean) {
    func += ` else if (Utils.cmp(test, "${name}(ebool,ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.${name}(FHE.asEbool(aBool), FHE.asEbool(bBool)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ebool.${name}(ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.asEbool(aBool).${name}(FHE.asEbool(bBool)).decrypt()) {
                return 1;
            }
            return 0;
        }`;
  }
  func += `
        revert TestNotFound(test);
    }`;

  const abi = `export interface ${capitalize(
    name
  )}TestType extends BaseContract {
    ${name}: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}\n`;
  return [generateTestContract(name, func), abi];
}

export function testContract1Arg(name: string) {
  const isEuint64Allowed = IsOperationAllowed(
    name,
    EInputType.indexOf("euint64")
  );
  const isEuint128Allowed = IsOperationAllowed(
    name,
    EInputType.indexOf("euint128")
  );
  const isEuint256Allowed = IsOperationAllowed(
    name,
    EInputType.indexOf("euint256")
  );
  let func = `function ${name}(string calldata test, uint256 a, int32 securityZone) public pure returns (uint256 output) {
        if (Utils.cmp(test, "${name}(euint8)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint8(a, securityZone)));
        } else if (Utils.cmp(test, "${name}(euint16)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint16(a, securityZone)));
        } else if (Utils.cmp(test, "${name}(euint32)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint32(a, securityZone)));
        }`;
  if (isEuint64Allowed) {
    func += ` else if (Utils.cmp(test, "${name}(euint64)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint64(a, securityZone)));
        }`;
  }
  if (isEuint128Allowed) {
    func += ` else if (Utils.cmp(test, "${name}(euint128)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint128(a, securityZone)));
        }`;
  }
  if (isEuint256Allowed) {
    func += ` else if (Utils.cmp(test, "${name}(euint256)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint256(a, securityZone)));
        }`;
  }
  func += ` else if (Utils.cmp(test, "${name}(ebool)")) {
            bool aBool = true;
            if (a == 0) {
                aBool = false;
            }

            if (FHE.decrypt(FHE.${name}(FHE.asEbool(aBool, securityZone)))) {
                return 1;
            }

            return 0;
        }
        
        revert TestNotFound(test);
    }`;
  const abi = `export interface ${capitalize(
    name
  )}TestType extends BaseContract {
    ${name}: (test: string, a: bigint) => Promise<bigint>;
}\n`;
  return [generateTestContract(name, func), abi];
}

export function generateTestContract(
  name: string,
  testFunc: string,
  importTypes: Array<string> = []
) {
  const importStatement = importTypes.length > 0
    ? `\nimport {${importTypes.join(", ")}} from "../../FHE.sol";`
    : "";
  return `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";${importStatement}
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract ${capitalize(name)}Test {
    using Utils for *;
    
    ${testFunc}
}
`;
}

export function testContractReq() {
  // Req is failing on EthCall so we need to make it as tx for now
  let func = `function req(string calldata test, uint256 a) public {
        if (Utils.cmp(test, "req(euint8)")) {
            FHE.req(FHE.asEuint8(a));
        } else if (Utils.cmp(test, "req(euint16)")) {
            FHE.req(FHE.asEuint16(a));
        } else if (Utils.cmp(test, "req(euint32)")) {
            FHE.req(FHE.asEuint32(a));
        } else if (Utils.cmp(test, "req(euint64)")) {
            FHE.req(FHE.asEuint64(a));
        } else if (Utils.cmp(test, "req(euint128)")) {
            FHE.req(FHE.asEuint128(a));
        } else if (Utils.cmp(test, "req(euint256)")) {
            FHE.req(FHE.asEuint256(a));
        } else if (Utils.cmp(test, "req(ebool)")) {
            bool b = true;
            if (a == 0) {
                b = false;
            }
            FHE.req(FHE.asEbool(b));
        } else {
            revert TestNotFound(test);
        }
    }`;
  const abi = `export interface ReqTestType extends BaseContract {
    req: (test: string, a: bigint) => Promise<{}>;
}\n`;
  return [generateTestContract("req", func), abi];
}

export function testContractReencrypt() {
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
  const abi = `export interface SealoutputTestType extends BaseContract {
    ${SEALING_FUNCTION_NAME}: (test: string, a: bigint, pubkey: Uint8Array) => Promise<string>;
}\n`;
  return [generateTestContract(SEALING_FUNCTION_NAME, func, ["ebool", "euint8"]), abi];
}

export function testContract3Arg(name: string) {
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
  )}TestType extends BaseContract {
    ${name}: (test: string, c: boolean, a: bigint, b: bigint) => Promise<bigint>;
}\n`;
  return [generateTestContract(name, func, ["ebool"]), abi];
}

export const IsOperationAllowed = (
  functionName: string,
  inputIdx: number
): boolean => {
  const regexes = AllowedOperations[inputIdx];
  for (let regex of regexes) {
    if (!new RegExp(regex).test(functionName.toLowerCase())) {
      return false;
    }
  }

  return true;
};

export function testContract2Arg(
  name: string,
  isBoolean: boolean,
  op?: string
) {
  const isEuint64Allowed = IsOperationAllowed(
    name,
    EInputType.indexOf("euint64")
  );
  const isEuint128Allowed = IsOperationAllowed(
    name,
    EInputType.indexOf("euint128")
  );
  const isEuint256Allowed = IsOperationAllowed(
    name,
    EInputType.indexOf("euint256")
  );
  let func = `function ${name}(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "${name}(euint8,euint8)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "${name}(euint16,euint16)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "${name}(euint32,euint32)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint32(a), FHE.asEuint32(b)));
        }`;
  if (isEuint64Allowed) {
    func += ` else if (Utils.cmp(test, "${name}(euint64,euint64)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint64(a), FHE.asEuint64(b)));
        }`;
  }
  if (isEuint128Allowed) {
    func += ` else if (Utils.cmp(test, "${name}(euint128,euint128)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint128(a), FHE.asEuint128(b)));
        }`;
  }
  if (isEuint256Allowed) {
    func += ` else if (Utils.cmp(test, "${name}(euint256,euint256)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint256(a), FHE.asEuint256(b))); 
        }`;
  }
  func += ` else if (Utils.cmp(test, "euint8.${name}(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).${name}(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.${name}(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).${name}(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.${name}(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).${name}(FHE.asEuint32(b)));
        }`;

  if (isEuint64Allowed) {
    func += ` else if (Utils.cmp(test, "euint64.${name}(euint64)")) {
            return FHE.decrypt(FHE.asEuint64(a).${name}(FHE.asEuint64(b)));
        }`;
  }
  if (isEuint128Allowed) {
    func += ` else if (Utils.cmp(test, "euint128.${name}(euint128)")) {
            return FHE.decrypt(FHE.asEuint128(a).${name}(FHE.asEuint128(b)));
        }`;
  }
  if (isEuint256Allowed) {
    func += ` else if (Utils.cmp(test, "euint256.${name}(euint256)")) {
            return FHE.decrypt(FHE.asEuint256(a).${name}(FHE.asEuint256(b)));
        }`;
  }
  if (op) {
    func += ` else if (Utils.cmp(test, "euint8 ${op} euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) ${op} FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 ${op} euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) ${op} FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 ${op} euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) ${op} FHE.asEuint32(b));
        }`;
    if (isEuint64Allowed) {
      func += ` else if (Utils.cmp(test, "euint64 ${op} euint64")) {
            return FHE.decrypt(FHE.asEuint64(a) ${op} FHE.asEuint64(b));
        }`;
    }
    if (isEuint128Allowed) {
      func += ` else if (Utils.cmp(test, "euint128 ${op} euint128")) {
            return FHE.decrypt(FHE.asEuint128(a) ${op} FHE.asEuint128(b));
        }`;
    }
    if (isEuint256Allowed) {
      func += ` else if (Utils.cmp(test, "euint256 ${op} euint256")) {
            return FHE.decrypt(FHE.asEuint256(a) ${op} FHE.asEuint256(b));
        }`;
    }
  }
  if (isBoolean) {
    func += ` else if (Utils.cmp(test, "${name}(ebool,ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.${name}(FHE.asEbool(aBool), FHE.asEbool(bBool)))) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "ebool.${name}(ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.asEbool(aBool).${name}(FHE.asEbool(bBool)).decrypt()) {
                return 1;
            }
            return 0;
        }`;
    if (op) {
      func += ` else if (Utils.cmp(test, "ebool ${op} ebool")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.asEbool(aBool) ${op} FHE.asEbool(bBool))) {
                return 1;
            }
            return 0;
        }`;
    }
  }
  func += `
    
        revert TestNotFound(test);
    }`;
  const abi = `export interface ${capitalize(
    name
  )}TestType extends BaseContract {
    ${name}: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}\n`;
  return [generateTestContract(name, func), abi];
}