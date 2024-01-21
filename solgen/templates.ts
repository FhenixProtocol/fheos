import {
  AllTypes,
  EInputType,
  EPlaintextType,
  EUintType,
  SEALING_FUNCTION_NAME,
  UnderlyingTypes,
  UintTypes,
  valueIsEncrypted,
  valueIsPlaintext,
} from "./common";

export const preamble = () => {
  return `// SPDX-License-Identifier: BSD-3-Clause-Clear
// solhint-disable one-contract-per-file

pragma solidity >=0.8.19 <0.9.0;

import {Precompiles, FheOps} from "./FheOS.sol";

${EInputType.map((type) => {
  return `type ${type} is uint256;`;
}).join("\n")}

${EInputType.map((type) => {
  return `struct in${capitalize(type)} {
    bytes data;
}`;
}).join("\n")}

library Common {
    // Values used to communicate types to the runtime.
    uint8 internal constant EBOOL_TFHE_GO = 0;
    uint8 internal constant EUINT8_TFHE_GO = 0;
    uint8 internal constant EUINT16_TFHE_GO = 1;
    uint8 internal constant EUINT32_TFHE_GO = 2;

    function bigIntToBool(uint256 i) internal pure returns (bool) {
        return (i > 0);
    }

    function bigIntToUint8(uint256 i) internal pure returns (uint8) {
        return uint8(i);
    }

    function bigIntToUint16(uint256 i) internal pure returns (uint16) {
        return uint16(i);
    }

    function bigIntToUint32(uint256 i) internal pure returns (uint32) {
        return uint32(i);
    }

    function bigIntToUint64(uint256 i) internal pure returns (uint64) {
        return uint64(i);
    }

    function bigIntToUint128(uint256 i) internal pure returns (uint128) {
        return uint128(i);
    }

    function bigIntToUint256(uint256 i) internal pure returns (uint256) {
        return i;
    }
    
    function toBytes(uint256 x) internal pure returns (bytes memory b) {
        b = new bytes(32);
        assembly { mstore(add(b, 32), x) }
    }
    
}

library Impl {
    function sealoutput(uint8 utype, uint256 ciphertext, bytes32 publicKey) internal pure returns (bytes memory reencrypted) {
        // Call the sealoutput precompile.
        reencrypted = FheOps(Precompiles.Fheos).sealOutput(utype, Common.toBytes(ciphertext), bytes.concat(publicKey));

        return reencrypted;
    }

    function verify(bytes memory _ciphertextBytes, uint8 _toType) internal pure returns (uint256 result) {
        bytes memory output;

        // Call the verify precompile.
        output = FheOps(Precompiles.Fheos).verify(_toType, _ciphertextBytes);
        result = getValue(output);
    }

    function cast(uint8 utype, uint256 ciphertext, uint8 toType) internal pure returns (uint256 result) {
        bytes memory output;

        // Call the cast precompile.
        output = FheOps(Precompiles.Fheos).cast(utype, Common.toBytes(ciphertext), toType);
        result = getValue(output);
    }

    function getValue(bytes memory a) internal pure returns (uint256 value) {
        assembly {
            value := mload(add(a, 0x20))
        }
    }

    function trivialEncrypt(uint256 value, uint8 toType) internal pure returns (uint256 result) {
        bytes memory output;

        // Call the trivialEncrypt precompile.
        output = FheOps(Precompiles.Fheos).trivialEncrypt(Common.toBytes(value), toType);

        result = getValue(output);
    }

    function select(uint8 utype, uint256 control, uint256 ifTrue, uint256 ifFalse) internal pure returns (uint256 result) {
        bytes memory output;

        // Call the trivialEncrypt precompile.
        output = FheOps(Precompiles.Fheos).select(utype, Common.toBytes(control), Common.toBytes(ifTrue), Common.toBytes(ifFalse));

        result = getValue(output);
    }
}

library FHE {
    euint8 public constant NIL8 = euint8.wrap(0);
    euint16 public constant NIL16 = euint16.wrap(0);
    euint32 public constant NIL32 = euint32.wrap(0);

    // Return true if the encrypted integer is initialized and false otherwise.
    function isInitialized(ebool v) internal pure returns (bool) {
        return ebool.unwrap(v) != 0;
    }

    // Return true if the encrypted integer is initialized and false otherwise.
    function isInitialized(euint8 v) internal pure returns (bool) {
        return euint8.unwrap(v) != 0;
    }

    // Return true if the encrypted integer is initialized and false otherwise.
    function isInitialized(euint16 v) internal pure returns (bool) {
        return euint16.unwrap(v) != 0;
    }

    // Return true if the encrypted integer is initialized and false otherwise.
    function isInitialized(euint32 v) internal pure returns (bool) {
        return euint32.unwrap(v) != 0;
    }

    function getValue(bytes memory a) private pure returns (uint256 value) {
        assembly {
            value := mload(add(a, 0x20))
        }
    }
    
    function mathHelper(
        uint8 utype,
        uint256 lhs,
        uint256 rhs,
        function(uint8, bytes memory, bytes memory) external pure returns (bytes memory) impl
    ) internal pure returns (uint256 result) {
        bytes memory output;
        output = impl(utype, Common.toBytes(lhs), Common.toBytes(rhs));
        result = getValue(output);
    }
    `;
};

export const PostFix = () => {
  return `\n}`;
};

const castFromEncrypted = (
  fromType: string,
  toType: string,
  name: string
): string => {
  if (!valueIsEncrypted(toType)) {
    console.log(`Unsupported type for casting: ${toType}`);
    process.exit(1);
  }
  return `Impl.cast(${
    UintTypes[toType]
  }, ${fromType}.unwrap(${name}), Common.${toType.toUpperCase()}_TFHE_GO)`;
};

const castFromPlaintext = (name: string, toType: string): string => {
  return `Impl.trivialEncrypt(${name}, Common.${toType.toUpperCase()}_TFHE_GO)`;
};

const castFromBytes = (name: string, toType: string): string => {
  return `Impl.verify(${name}, Common.${toType.toUpperCase()}_TFHE_GO)`;
};

const castFromInputType = (name: string, toType: string): string => {
  return `FHE.as${capitalize(toType)}(${name}.data)`;
};

const castToEbool = (name: string, fromType: string): string => {
  return `
    \n    /// @notice Converts a ${fromType} to an ebool
    function asEbool(${fromType} value) internal pure returns (ebool) {
        return ne(${name}, as${capitalize(fromType)}(0));
    }`;
};

export const AsTypeFunction = (fromType: string, toType: string) => {
  let castString = castFromEncrypted(fromType, toType, "value");

  let docString = `
    /// @notice Converts a ${fromType} to an ${toType}`;

  if (fromType === "bool" && toType === "ebool") {
    return `
    /// @notice Converts a plaintext boolean value to a ciphertext ebool
    /// @dev Privacy: The input value is public, therefore the ciphertext should be considered public and should be used
    ///only for mathematical operations, not to represent data that should be private
    /// @return A ciphertext representation of the input 
    function asEbool(bool value) internal pure returns (ebool) {
        uint256 sVal = 0;
        if (value) {
            sVal = 1;
        }

        return asEbool(sVal);
    }`;
  } else if (fromType.startsWith("in")) {
    docString = `
    /// @notice Parses input ciphertexts from the user. Converts from encrypted raw bytes to an ${toType}
    /// @dev Also performs validation that the ciphertext is valid and has been encrypted using the network encryption key
    /// @return a ciphertext representation of the input`;
    castString = castFromInputType("value", toType);

    return `${docString}
    function as${capitalize(
      toType
    )}(${fromType} memory value) internal pure returns (${toType}) {
        return ${castString};
    }`;
  } else if (fromType === "bytes memory") {
    docString = `
    /// @notice Parses input ciphertexts from the user. Converts from encrypted raw bytes to an ${toType}
    /// @dev Also performs validation that the ciphertext is valid and has been encrypted using the network encryption key
    /// @return a ciphertext representation of the input`;
    castString = castFromBytes("value", toType);
  } else if (EPlaintextType.includes(fromType)) {
    castString = castFromPlaintext("value", toType);
  } else if (toType === "ebool") {
    return castToEbool("value", fromType);
  } else if (!EInputType.includes(fromType)) {
    throw new Error(`Unsupported type for casting: ${fromType}`);
  }

  return `${docString}
    function as${capitalize(
      toType
    )}(${fromType} value) internal pure returns (${toType}) {
        return ${toType}.wrap(${castString});
    }`;
};

function TypeCastTestingFunction(
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
  retTypeTs = retTypeTs.includes("uint") ? "bigint" : retTypeTs;

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
            return FHE.decrypt(${encryptedVal}.to${shortenType(toType)}());
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
  for (const fromType of EInputType.concat("uint256", "bytes memory")) {
    if (type === fromType) {
      continue;
    }

    const fromTypeTs = fromType === "bytes memory" ? "Uint8Array" : `bigint`;
    const fromTypeSol = fromType === "bytes memory" ? fromType : `uint256`;
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

  return [generateTestContract(`As${capitalize(type)}`, funcs), abi];
}

const unwrapType = (typeName: EUintType, inputName: string): string =>
  `${typeName}.unwrap(${inputName})`;
const wrapType = (resultType: EUintType, inputName: string): string =>
  `${resultType}.wrap(${inputName})`;
const asEuintFuncName = (typeName: EUintType): string =>
  `as${capitalize(typeName)}`;
export const capitalize = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);

export function SolTemplate2Arg(
  name: string,
  input1: AllTypes,
  input2: AllTypes,
  returnType: AllTypes
) {
  // special names for reencrypt function (don't check name === reencrypt) because the name could change
  let variableName1 = input2 === "bytes32" ? "value" : "lhs";
  let variableName2 = input2 === "bytes32" ? "publicKey" : "rhs";

  let docString = `
    /// @notice This functions performs the ${name} operation
    /// @dev If any of the inputs are expected to be a ciphertext, it verifies that the value matches a valid ciphertext
    ///Pure in this function is marked as a hack/workaround - note that this function is NOT pure as fetches of ciphertexts require state access
    /// @param lhs The first input 
    /// @param rhs The second input
    /// @return The result of the operation
    `;

  // reencrypt
  if (variableName2 === "publicKey") {
    docString = `
    /// @notice performs the ${name} function on a ${input1} ciphertext. This operation returns the plaintext value, sealed for the public key provided 
    /// @dev Pure in this function is marked as a hack/workaround - note that this function is NOT pure as fetches of ciphertexts require state access
    /// @param value Ciphertext to decrypt and seal
    /// @param publicKey Public Key that will receive the sealed plaintext
    /// @return Plaintext input, sealed for the owner of \`publicKey\`
    `;
  }

  let funcBody = docString;

  funcBody += `function ${name}(${input1} ${variableName1}, ${input2} ${variableName2}) internal pure returns (${returnType}) {`;

  if (valueIsEncrypted(input1)) {
    // both inputs encrypted - this is a generic math function. i.e. div, mul, eq, etc.
    // 1. possibly cast input1
    // 2. possibly cast input2
    // 3. possibly cast return type
    if (valueIsEncrypted(input2) && input1 !== input2) {
      return "";
    }
    if (valueIsEncrypted(input2) || EPlaintextType.includes(input2)) {
      let input2Cast =
        input1 === input2
          ? variableName2
          : `${asEuintFuncName(input1)}(${variableName2})`;
      //
      funcBody += `
        if (!isInitialized(${variableName1})) {
            ${variableName1} = ${asEuintFuncName(input1)}(0);
        }
        if (!isInitialized(${variableName2})) {
            ${variableName2} = ${asEuintFuncName(input1)}(0);
        }
        ${UnderlyingTypes[input1]} unwrappedInput1 = ${unwrapType(
        input1,
        variableName1
      )};
        ${UnderlyingTypes[input1]} unwrappedInput2 = ${unwrapType(
        input1,
        `${input2Cast}`
      )};
`;
      if (valueIsEncrypted(returnType)) {
        funcBody += `
        ${UnderlyingTypes[returnType]} result = mathHelper(${
          UintTypes[input1]
        }, unwrappedInput1, unwrappedInput2, FheOps(Precompiles.Fheos).${name});
        return ${wrapType(returnType, "result")};`;
      } else {
        funcBody += `
        ${returnType} result = mathHelper(${UintTypes[input1]}, unwrappedInput1, unwrappedInput2, FheOps(Precompiles.Fheos).${name});
        return result;`;
      }
    } else if (input2 === "bytes32") {
      // **** Value 1 is encrypted, value 2 is bytes32 - this is basically reencrypt/wrapForUser
      funcBody += `
        if (!isInitialized(${variableName1})) {
            ${variableName1} = ${asEuintFuncName(input1)}(0);
        }
        ${UnderlyingTypes[input1]} unwrapped = ${unwrapType(
        input1,
        variableName1
      )};

        return Impl.${name}(${
        UintTypes[input1]
      }, unwrapped, ${variableName2});`;
    }
  } else {
    // don't support input 1 is plaintext
    throw new Error("Unsupported plaintext input1");
  }

  funcBody += `\n    }`;

  return funcBody;
}

export function testContract2ArgBoolRes(name: string, isBoolean: boolean) {
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
        } else if (Utils.cmp(test, "euint8.${name}(euint8)")) {
            if (FHE.decrypt(FHE.asEuint8(a).${name}(FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.${name}(euint16)")) {
            if (FHE.decrypt(FHE.asEuint16(a).${name}(FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.${name}(euint32)")) {
            if (FHE.decrypt(FHE.asEuint32(a).${name}(FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        }`;
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
            if (FHE.decrypt(FHE.asEbool(aBool).${name}(FHE.asEbool(bBool)))) {
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
  let func = `function ${name}(string calldata test, uint256 a) public pure returns (uint256 output) {
        if (Utils.cmp(test, "${name}(euint8)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint8(a)));
        } else if (Utils.cmp(test, "${name}(euint16)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint16(a)));
        } else if (Utils.cmp(test, "${name}(euint32)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint32(a)));
        } else if (Utils.cmp(test, "${name}(ebool)")) {
            bool aBool = true;
            if (a == 0) {
                aBool = false;
            }

            if (FHE.decrypt(FHE.${name}(FHE.asEbool(aBool)))) {
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
  importTypes: boolean = false
) {
  const importStatement = importTypes
    ? `\nimport {ebool} from "../../FHE.sol";`
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
  let func = `function ${SEALING_FUNCTION_NAME}(string calldata test, uint256 a, bytes32 pubkey) public pure returns (bytes memory reencrypted) {
        if (Utils.cmp(test, "${SEALING_FUNCTION_NAME}(euint8)")) {
            return FHE.${SEALING_FUNCTION_NAME}(FHE.asEuint8(a), pubkey);
        } else if (Utils.cmp(test, "${SEALING_FUNCTION_NAME}(euint16)")) {
            return FHE.${SEALING_FUNCTION_NAME}(FHE.asEuint16(a), pubkey);
        } else if (Utils.cmp(test, "${SEALING_FUNCTION_NAME}(euint32)")) {
            return FHE.${SEALING_FUNCTION_NAME}(FHE.asEuint32(a), pubkey);
        } else if (Utils.cmp(test, "${SEALING_FUNCTION_NAME}(ebool)")) {
            bool b = true;
            if (a == 0) {
                b = false;
            }

            return FHE.${SEALING_FUNCTION_NAME}(FHE.asEbool(b), pubkey);
        } 
        
        revert TestNotFound(test);
    }`;
  const abi = `export interface SealoutputTestType extends BaseContract {
    ${SEALING_FUNCTION_NAME}: (test: string, a: bigint, pubkey: Uint8Array) => Promise<Uint8Array>;
}\n`;
  return [generateTestContract(SEALING_FUNCTION_NAME, func), abi];
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
  return [generateTestContract(name, func, true), abi];
}

export function testContract2Arg(
  name: string,
  isBoolean: boolean,
  op?: string
) {
  let func = `function ${name}(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "${name}(euint8,euint8)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "${name}(euint16,euint16)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "${name}(euint32,euint32)")) {
            return FHE.decrypt(FHE.${name}(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.${name}(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).${name}(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.${name}(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).${name}(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.${name}(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).${name}(FHE.asEuint32(b)));
        }`;
  if (op) {
    func += ` else if (Utils.cmp(test, "euint8 ${op} euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) ${op} FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 ${op} euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) ${op} FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 ${op} euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) ${op} FHE.asEuint32(b));
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
            if (FHE.decrypt(FHE.asEbool(aBool).${name}(FHE.asEbool(bBool)))) {
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

export function genAbiFile(abi: string) {
  return `import { BaseContract } from 'ethers';\n${abi}\n\n`;
}

export function SolTemplate1Arg(
  name: string,
  input1: AllTypes,
  returnType: AllTypes
) {
  let docString = `
    /// @notice Performs the ${name} operation on a ciphertext
    /// @dev Verifies that the input value matches a valid ciphertext. Pure in this function is marked as a hack/workaround - note that this function is NOT pure as fetches of ciphertexts require state access
    /// @param input1 the input ciphertext
    `;

  if (name === "not" && input1 === "ebool") {
    return `\n
    /// @notice Performs the "not" for the ebool type
    /// @dev Implemented by a workaround due to ebool being a euint8 type behind the scenes, therefore xor is needed to assure that not(true) = false and vise-versa
    /// @param value input ebool ciphertext
    /// @return Result of the not operation on \`value\` 
    function not(ebool value) internal pure returns (ebool) {
        return xor(value, asEbool(true));
    }`;
  }

  let returnStr = returnType === "none" ? `` : `returns (${returnType})`;

  let funcBody = docString;

  funcBody += `function ${name}(${input1} input1) internal pure ${returnStr} {`;

  if (valueIsEncrypted(input1)) {
    // Get the proper function
    funcBody += `
        if (!isInitialized(input1)) {
            input1 = ${asEuintFuncName(input1)}(0);
        }`;
    let unwrap = `${UnderlyingTypes[input1]} unwrappedInput1 = ${unwrapType(
      input1,
      "input1"
    )};`;
    let getResult = (inputName: string) =>
      `FheOps(Precompiles.Fheos).${name}(${UintTypes[input1]}, ${inputName});`;

    if (valueIsEncrypted(returnType)) {
      // input and return type are encrypted - not/neg other unary functions
      funcBody += `
        ${unwrap}
        bytes memory inputAsBytes = Common.toBytes(unwrappedInput1);
        bytes memory b = ${getResult("inputAsBytes")}
        uint256 result = Impl.getValue(b);
        return ${wrapType(returnType, "result")};
    }`;
    } else if (returnType === "none") {
      // this is essentially req
      funcBody += `
        ${unwrap}
        bytes memory inputAsBytes = Common.toBytes(unwrappedInput1);
        ${getResult("inputAsBytes")}
    }`;
    } else if (valueIsPlaintext(returnType)) {
      let returnTypeCamelCase =
        returnType.charAt(0).toUpperCase() + returnType.slice(1);
      let outputConvertor = `Common.bigIntTo${returnTypeCamelCase}(result);`;
      funcBody += `
        ${unwrap}
        bytes memory inputAsBytes = Common.toBytes(unwrappedInput1);
        uint256 result = ${getResult("inputAsBytes")}
        return ${outputConvertor}
    }`;
    }
  } else {
    throw new Error("unsupported function of 1 input that is not encrypted");
  }
  return funcBody;
}

export function SolTemplate3Arg(
  name: string,
  input1: AllTypes,
  input2: AllTypes,
  input3: AllTypes,
  returnType: AllTypes
) {
  if (valueIsEncrypted(returnType)) {
    if (
      valueIsEncrypted(input1) &&
      valueIsEncrypted(input2) &&
      valueIsEncrypted(input3) &&
      input1 === "ebool"
    ) {
      if (input2 !== input3) {
        return "";
      }

      return `\n
    function ${name}(${input1} input1, ${input2} input2, ${input3} input3) internal pure returns (${returnType}) {
        if (!isInitialized(input1)) {
            input1 = ${asEuintFuncName(input1)}(0);
        }
        if (!isInitialized(input2)) {
            input2 = ${asEuintFuncName(input2)}(0);
        }
        if (!isInitialized(input3)) {
            input3 = ${asEuintFuncName(input3)}(0);
        }

        ${UnderlyingTypes[input1]} unwrappedInput1 = ${unwrapType(
        input1,
        "input1"
      )};
        ${UnderlyingTypes[input2]} unwrappedInput2 = ${unwrapType(
        input2,
        "input2"
      )};
        ${UnderlyingTypes[input3]} unwrappedInput3 = ${unwrapType(
        input3,
        `input3`
      )};

        ${UnderlyingTypes[returnType]} result = Impl.${name}(${
        UintTypes[input2]
      }, unwrappedInput1, unwrappedInput2, unwrappedInput3);
        return ${wrapType(returnType, "result")};
    }`;
    } else {
      return "";
    }
  } else {
    throw new Error(`Unsupported return type ${returnType} for 3 args`);
  }
}

function operatorFunctionName(
  funcName: string,
  forType: "ebool" | "euint8" | "euint16" | "euint32"
) {
  return `operator${capitalize(funcName)}${capitalize(forType)}`;
}

export const OperatorOverloadDecl = (
  funcName: string,
  op: string,
  forType: EUintType,
  unary: boolean,
  returnsBool: boolean
) => {
  let opOverloadName = operatorFunctionName(funcName, forType);
  let unaryParameters = unary ? "lhs" : "lhs, rhs";
  let funcParams = unaryParameters
    .split(", ")
    .map((key) => {
      return `${forType} ${key}`;
    })
    .join(", ");
  let returnType = returnsBool ? "ebool" : forType;

  return `\nusing {${opOverloadName} as ${op}} for ${forType} global;
/// @notice Performs the ${funcName} operation
function ${opOverloadName}(${funcParams}) pure returns (${returnType}) {
    return FHE.${funcName}(${unaryParameters});
}\n`;
};

export const BindingLibraryType = (type: string) => {
  let typeCap = capitalize(type);
  return `\n\nusing Bindings${typeCap} for ${type} global;
library Bindings${typeCap} {`;
};

export const OperatorBinding = (
  funcName: string,
  forType: string,
  unary: boolean,
  returnsBool: boolean
) => {
  let unaryParameters = unary ? "lhs" : "lhs, rhs";
  let funcParams = unaryParameters
    .split(", ")
    .map((key) => {
      return `${forType} ${key}`;
    })
    .join(", ");
  let returnType = returnsBool ? "ebool" : forType;

  let docString = `
    /// @notice Performs the ${funcName} operation
    /// @dev Pure in this function is marked as a hack/workaround - note that this function is NOT pure as fetches of ciphertexts require state access
    /// @param lhs input of type ${forType}
    `;

  if (unary) {
    docString += `/// @param rhs second input of type ${forType}\n`;
  }

  docString += `/// @return the result of the ${funcName}`;

  return `
    ${docString}
    function ${funcName}(${funcParams}) internal pure returns (${returnType}) {
        return FHE.${funcName}(${unaryParameters});
    }`;
};

const shortenType = (type: string) => {
  return type === "ebool" ? "Bool" : "U" + type.slice(5); // get only number at the end
};

export const CastBinding = (thisType: string, targetType: string) => {
  return `
    function to${shortenType(
      targetType
    )}(${thisType} value) internal pure returns (${targetType}) {
        return FHE.as${capitalize(targetType)}(value);
    }`;
};
