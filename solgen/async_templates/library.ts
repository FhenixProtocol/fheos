import {
  AllTypes,
  EInputType,
  EPlaintextType,
  EUintType,
  SEAL_RETURN_TYPE,
  UnderlyingTypes,
  UintTypes,
  valueIsEncrypted,
  valueIsPlaintext,
  LOCAL_SEAL_FUNCTION_NAME,
  LOCAL_DECRYPT_FUNCTION_NAME,
  AllowedOperations,
  AllowedTypesOnCastToEaddress,
  toPlaintextType,
  capitalize,
  shortenType,
  SealedOutputStructs,
  SEALING_FUNCTION_NAME,
  SEALING_TYPED_FUNCTION_NAME,
  LOCAL_SEAL_TYPED_FUNCTION_NAME,
} from "../common";

export const preamble = () => {
  return `// SPDX-License-Identifier: BSD-3-Clause-Clear
// solhint-disable one-contract-per-file

pragma solidity >=0.8.19 <0.9.0;

// import {Console} from "@fhenixprotocol/contracts/utils/debug/Console.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import {Precompiles, FheOps} from "./FheOS.sol";

${EInputType.map((type) => {
  return `type ${type} is uint256;`;
}).join("\n")}

${EInputType.map((type) => {
  return `struct in${capitalize(type)} {
    bytes data;
    int32 securityZone;
}`;
}).join("\n")}

struct SealedArray {
  bytes[] data;
}
${SealedOutputStructs.map((struct) => {
  let docstring = `
/// @dev Utility structure providing clients with type context of a sealed output string.
/// Return type of \`FHE.sealoutputTyped\` and \`sealTyped\` within the binding libraries.`
  if (struct === `SealedBool`) {
    docstring +=  `
/// \`utype\` representing Bool is 13. See \`FHE.sol\` for more.`;
  }
  if (struct === `SealedUint`) {
    docstring += `
/// \`utype\` representing Uints is 0-5. See \`FHE.sol\` for more.`;
    docstring += `
/// \`utype\` map: {uint8: 0} {uint16: 1} {uint32: 2} {uint64: 3} {uint128: 4} {uint256: 5}.`;
  }
  if (struct === `SealedAddress`) {
    docstring +=  `
/// \`utype\` representing Address is 12. See \`FHE.sol\` for more.`;
  }
  return `${docstring}
struct ${struct} {
    string data;
    uint8 utype;
}`;
}).join("\n")}

library TaskManager {
	//solhint-disable const-name-snakecase
	address public constant TASK_MANAGER_ADDRESS = address(129);
}

library Common {
    // Values used to communicate types to the runtime.
    // Must match values defined in warp-drive protobufs for everything to 
    // make sense
    uint8 internal constant EUINT8_TFHE = 0;
    uint8 internal constant EUINT16_TFHE = 1;
    uint8 internal constant EUINT32_TFHE = 2;
    uint8 internal constant EUINT64_TFHE = 3;
    uint8 internal constant EUINT128_TFHE = 4;
    uint8 internal constant EUINT256_TFHE = 5;
    uint8 internal constant EADDRESS_TFHE = 12;
    // uint8 internal constant INT_BGV = 12;
    uint8 internal constant EBOOL_TFHE = 13;
    
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

    function bigIntToAddress(uint256 i) internal pure returns (address) {
        return address(uint160(i));
    }
    
    function toBytes(uint256 x) internal pure returns (bytes memory b) {
        b = new bytes(32);
        assembly { mstore(add(b, 32), x) }
    }
    
    function bytesToUint256(bytes memory b) internal pure returns (uint256) {
        require(b.length == 32, string(abi.encodePacked("Input bytes length must be 32, but got ", Strings.toString(b.length))));

        uint256 result;
        assembly {
            result := mload(add(b, 32))
        }
        return result;
    }

    function hexCharToUint8(bytes1 char) internal pure returns (uint8) {
        if (char >= "0" && char <= "9") {
            return uint8(char) - uint8(bytes1("0"));
        } else if (char >= "a" && char <= "f") {
            return uint8(char) - uint8(bytes1("a")) + 10;
        } else if (char >= "A" && char <= "F") {
            return uint8(char) - uint8(bytes1("A")) + 10;
        } else {
            revert("Invalid hex character");
        }
    }

    function hexStringToUint(string memory hexString) internal pure returns (uint8) {
        require(bytes(hexString).length == 2, "Invalid hex string length");

        uint8 value = 0;
        for (uint8 i = 0; i < 2; i++) {
            value = value * 16 + hexCharToUint8(bytes(hexString)[i]);
        }

        return value;
    }

    function hexStringToBytes32(string memory hexString) internal pure returns (bytes memory) {
        bytes memory hexBytes = bytes(hexString);
        // Ensure the string has the correct length (64 characters for 32 bytes)
        require(hexBytes.length == 64, "Invalid hex string length");

        // Iterate every 2 bytes in string, consider them as 1 byte
        bytes memory bb = new bytes(32);
        string memory l = "";
        for (uint i = 0; i < 32; i++) {
            l = string(abi.encodePacked("", hexBytes[i * 2], hexBytes[i * 2 + 1]));
            bb[i] = bytes1(hexStringToUint(l));
        }

        return bb;
    }

    function bytesArrayToString(bytes memory a) internal pure returns (string memory) {
        string memory b = "[";
        for (uint i = 0; i < a.length; i++) {
            b = string(abi.encodePacked(b, Strings.toHexString(uint8(a[i])), " "));
        }

        b = string(abi.encodePacked(b, "]"));
        return b;
    }

    function functionCodeToBytes1(string memory functionCode) internal pure returns (bytes memory) {
        // Convert the hex string to bytes
        bytes memory result = new bytes(1);
        assembly {
            result := mload(add(functionCode, 1)) // Load the bytes directly from memory
        }

        return result;
    }

    function bytesToHexString(bytes memory buffer) internal pure returns (string memory) {
        // Each byte takes 2 characters
        bytes memory hexChars = new bytes(buffer.length * 2);
        
        for(uint i = 0; i < buffer.length; i++) {
            uint8 value = uint8(buffer[i]);
            hexChars[i * 2] = byteToChar(value / 16);
            hexChars[i * 2 + 1] = byteToChar(value % 16);
        }
        
        return string(hexChars);
    }

    // Helper function for bytesToHexString
    function byteToChar(uint8 value) internal pure returns (bytes1) {
        if (value < 10) {
            return bytes1(uint8(48 + value)); // 0-9
        } else {
            return bytes1(uint8(87 + value)); // a-f
        }
    }

    function uint256ToBytes32(uint256 value) internal pure returns (bytes memory) {
        bytes memory result = new bytes(32);
        assembly {
            mstore(add(result, 32), value)
        }
        return result;
    }

    function euint32ToUint256(euint32 value) internal pure returns (uint256) {
        return uint256(euint32.unwrap(value));
    }

    function euint32ToHexString(euint32 value) internal pure returns (string memory) {
        return bytesToHexString(uint256ToBytes32(euint32ToUint256(value)));
    }
}

library Impl {
    function verify(bytes memory _ciphertextBytes, uint8 _toType, int32 securityZone) internal pure returns (uint256 result) {
        bytes memory output;

        // Call the verify precompile.
        output = FheOps(Precompiles.Fheos).verify(_toType, _ciphertextBytes, securityZone);
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

    function trivialEncrypt(uint256 value, uint8 toType, int32 securityZone) internal pure returns (uint256 result) {
        bytes memory output;

        // Call the trivialEncrypt precompile.
        output = FheOps(Precompiles.Fheos).trivialEncrypt(Common.toBytes(value), toType, securityZone);

        result = getValue(output);
    }

    function select(uint8 utype, uint256 control, uint256 ifTrue, uint256 ifFalse) internal pure returns (uint256 result) {
        bytes memory output;

        // Call the trivialEncrypt precompile.
        output = FheOps(Precompiles.Fheos).select(utype, Common.toBytes(control), Common.toBytes(ifTrue), Common.toBytes(ifFalse));

        result = getValue(output);
    }
}

/// @title Interface for consumers that wanna receive result of decrypt tasks
/// @notice Implement the callback function in your contract to handle the decrypt result
interface IAsyncFHEReceiver {
    function handleDecryptResult(bytes memory ctHash, uint256 result) external;
    function handleSealOutputResult(bytes memory ctHash, ${SEAL_RETURN_TYPE} memory result) external;
}

interface ITaskManager {
    function createTask(uint256 ctHash, string memory operation, uint256 input1, uint256 input2) external;
    function handleOpResult(bytes memory tempKey, bytes memory result) external;
    function createDecryptTask(uint256 ctHash) external;
    function createSealOutputTask(uint256 ctHash, bytes32 publicKey) external;
}

library FHE {
    bytes private constant CT_HASH_MAGIC_BYTES = hex"deedbeaf";
    euint8 public constant NIL8 = euint8.wrap(0);
    euint16 public constant NIL16 = euint16.wrap(0);
    euint32 public constant NIL32 = euint32.wrap(0);
    
    // Default value for temp hash calculation in unary operations
    string private constant DEFAULT_VALUE = "0";

    // Order is set as in fheos/precompiles/types/types.go
    enum FunctionId {
        _0,         // 0 - GetNetworkKey
        _1,         // 1 - Verify
        _2,         // 2 - Cast
        sealoutput, // 3
        select,     // 4
        req,        // 5
        decrypt,    // 6
        sub,        // 7
        add,        // 8
        xor,        // 9
        and,        // 10
        or,         // 11
        not,        // 12
        div,        // 13
        rem,        // 14
        mul,        // 15
        shl,        // 16
        shr,        // 17
        gte,        // 18
        lte,        // 19
        lt,         // 20
        gt,         // 21
        min,        // 22
        max,        // 23
        eq,         // 24
        ne,         // 25
        _26,        // 26 - TrivialEncrypt
        random,     // 27
        rol,        // 28
        ror,        // 29
        square      // 30
    }

    /// @notice Calculates the temporary hash for unary operations
    /// @param value - The value to hash
    /// @param functionId - The function id
    /// @return The calculated temporary hash
    function calcUnaryPlaceholderValueHash(uint256 value, FunctionId functionId) private pure returns (uint256) {
        return calcBinaryPlaceholderValueHash(0, value, functionId);
    }

    /// @notice Calculates the temporary hash for async operations
    /// @dev Must result the same temp hash as calculated by warp-drive/fhe-driver/CalcBinaryPlaceholderValueHash
    /// @param lhsHash - Left hand side operand hash
    /// @param rhsHash - Right hand side operand hash
    /// @param functionId - The function id
    /// @return The calculated temporary hash
    function calcBinaryPlaceholderValueHash(
        uint256 lhsHash,
        uint256 rhsHash,
        FunctionId functionId
    ) private pure returns (uint256) {
        bytes memory lhsBytes = Common.uint256ToBytes32(lhsHash);
        bytes memory rhsBytes = Common.uint256ToBytes32(rhsHash);
        bytes1 functionIdByte = bytes1(uint8(functionId));
        bytes memory combined = bytes.concat(lhsBytes, rhsBytes, functionIdByte);

        // Calculate Keccak256 hash
        bytes memory hash = abi.encodePacked(keccak256(combined));

        // Copy magic bytes to the first four bytes of the hash
        for (uint i = 0; i < 4; i++) {
            hash[i] = CT_HASH_MAGIC_BYTES[i];
        }

        return uint256(bytes32(hash));
    }

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
    
    // Return true if the encrypted integer is initialized and false otherwise.
    function isInitialized(euint64 v) internal pure returns (bool) {
        return euint64.unwrap(v) != 0;
    }
    
        // Return true if the encrypted integer is initialized and false otherwise.
    function isInitialized(euint128 v) internal pure returns (bool) {
        return euint128.unwrap(v) != 0;
    }
    
        // Return true if the encrypted integer is initialized and false otherwise.
    function isInitialized(euint256 v) internal pure returns (bool) {
        return euint256.unwrap(v) != 0;
    }

    function isInitialized(eaddress v) internal pure returns (bool) {
        return eaddress.unwrap(v) != 0;
    }

    function getValue(bytes memory a) private pure returns (uint256 value) {
        assembly {
            value := mload(add(a, 0x20))
        }
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
  if (!valueIsEncrypted(fromType)) {
    return ""; // casting from plaintext type is handled elsewhere
  }
  return `Impl.cast(${
    UintTypes[fromType]
  }, ${fromType}.unwrap(${name}), Common.${toType.toUpperCase()}_TFHE)`;
};

const castFromPlaintext = (
  name: string,
  toType: string,
  addSecurityZone: boolean = false
): string => {
  return `Impl.trivialEncrypt(${name}, Common.${toType.toUpperCase()}_TFHE, ${
    addSecurityZone ? "securityZone" : "0"
  })`;
};

const castFromPlaintextAddress = (
  name: string,
  toType: string,
  addSecurityZone: boolean = false
): string => {
  return `Impl.trivialEncrypt(uint256(uint160(${name})), Common.${toType.toUpperCase()}_TFHE, ${
    addSecurityZone ? "securityZone" : "0"
  })`;
};

const castFromBytes = (name: string, toType: string): string => {
  return `Impl.verify(${name}, Common.${toType.toUpperCase()}_TFHE, securityZone)`;
};

const castFromInputType = (name: string, toType: string): string => {
  return `FHE.as${capitalize(toType)}(${name}.data, ${name}.securityZone)`;
};

const castToEbool = (name: string, fromType: string): string => {
  // todo (eshel): this should not work for non-default security zones because the second operand of the 'ne' is always from security zone 0.
  // check if the precompiles support the EBOOL_TFHE constant as with the other casts
  return `
    \n    /// @notice Converts a ${fromType} to an ebool
    function asEbool(${fromType} value) internal returns (ebool) {
        return ne(${name}, as${capitalize(fromType)}(0));
    }`;
};

export const AsTypeFunction = (
  fromType: string,
  toType: string,
  addSecurityZone: boolean = false
) => {
  if (
    toType === "eaddress" &&
    !AllowedTypesOnCastToEaddress.includes(fromType)
  ) {
    return ""; // skip unsupported cast
  }

  let castString = castFromEncrypted(fromType, toType, "value");
  let overrideFuncs = "";

  let docString = `
    /// @notice Converts a ${fromType} to an ${toType}`;

  if (fromType === "bool" && toType === "ebool") {
    return `
    /// @notice Converts a plaintext boolean value to a ciphertext ebool
    /// @dev Privacy: The input value is public, therefore the resulting ciphertext should be considered public until involved in an fhe operation
    /// @return A ciphertext representation of the input
    function asEbool(bool value) internal returns (ebool) {
        uint256 sVal = 0;
        if (value) {
            sVal = 1;
        }
        return asEbool(sVal);
    }
    /// @notice Converts a plaintext boolean value to a ciphertext ebool, specifying security zone
    /// @dev Privacy: The input value is public, therefore the resulting ciphertext should be considered public until involved in an fhe operation
    /// @return A ciphertext representation of the input
    function asEbool(bool value, int32 securityZone) internal returns (ebool) {
        uint256 sVal = 0;
        if (value) {
          sVal = 1;
        }
        return asEbool(sVal, securityZone);
    }`;
  } else if (fromType.startsWith("inE")) {
    docString = `
    /// @notice Parses input ciphertexts from the user. Converts from encrypted raw bytes to an ${toType}
    /// @dev Also performs validation that the ciphertext is valid and has been encrypted using the network encryption key
    /// @return a ciphertext representation of the input`;
    castString = castFromInputType("value", toType);

    return `${docString}
    function as${capitalize(
      toType
    )}(${fromType} memory value) internal returns (${toType}) {
        return ${castString};
    }`;
  } else if (fromType === "bytes memory") {
    docString = `
    /// @notice Parses input ciphertexts from the user. Converts from encrypted raw bytes to an ${toType}
    /// @dev Also performs validation that the ciphertext is valid and has been encrypted using the network encryption key
    /// @return a ciphertext representation of the input`;
    addSecurityZone = true;
    castString = castFromBytes("value", toType);
  } else if (EPlaintextType.includes(fromType)) {
    if (!addSecurityZone) {
      // recursive call to add the asType override with the security zone
      overrideFuncs += AsTypeFunction(fromType, toType, true);
    } else {
      docString += `, specifying security zone`;
    }

    docString += `
    /// @dev Privacy: The input value is public, therefore the resulting ciphertext should be considered public until involved in an fhe operation`;

    if (fromType === "address" && toType == "eaddress") {
      docString += `
    /// Allows for a better user experience when working with eaddresses`;
      castString = castFromPlaintextAddress("value", toType, addSecurityZone);
    } else {
      castString = castFromPlaintext("value", toType, addSecurityZone);
    }
  } else if (toType === "ebool") {
    return castToEbool("value", fromType);
  } else if (!EInputType.includes(fromType)) {
    throw new Error(`Unsupported type for casting: ${fromType}`);
  }

  let func = `${docString}
    function as${capitalize(toType)}(${fromType} value${
    addSecurityZone ? ", int32 securityZone" : ""
  }) internal returns (${toType}) {
        return ${toType}.wrap(${castString});
    }`;

  return func + overrideFuncs;
};

const unwrapType = (typeName: EUintType, inputName: string): string =>
  `${typeName}.unwrap(${inputName})`;
const wrapType = (resultType: EUintType, inputName: string): string =>
  `${resultType}.wrap(${inputName})`;
const asEuintFuncName = (typeName: EUintType): string =>
  `as${capitalize(typeName)}`;
const asDecryptedType = (typeName: EUintType): string =>
  `${typeName.toLowerCase().slice(1, typeName.length)}`;

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
    /// @notice This function performs the ${name} async operation
    /// @dev If any of the inputs are expected to be a ciphertext, it verifies that the value matches a valid ciphertext
    /// @param lhs The first input 
    /// @param rhs The second input
    /// @return The result of the operation
    `;

  // reencrypt (seal)
  if (name === SEALING_FUNCTION_NAME || name === SEALING_TYPED_FUNCTION_NAME) {
    docString = `
    /// @notice performs the ${name} async function on a ${input1} ciphertext. This operation returns the plaintext value, sealed for the public key provided 
    /// @param value Ciphertext to decrypt and seal
    /// @param publicKey Public Key that will receive the sealed plaintext
    `;

    // Diff return types of the sealing functions
    if (name === SEALING_FUNCTION_NAME) {
      docString += `/// @return Plaintext input, sealed for the owner of \`publicKey\`
    `;
    }
    if (name === SEALING_TYPED_FUNCTION_NAME) {
      // Example Output: "/// @return SealedBool({ data: Plaintext input, sealed for the owner of `publicKey`, utype: Common.EBOOL_TFHE })"
      const returnTypeClean = returnType.replace(" memory", "");
      docString += `/// @return ${returnTypeClean}({ data: Plaintext input, sealed for the owner of \`publicKey\`, utype: ${UintTypes[input1 as EUintType]} })
    `;
    }
  }

  let funcBody = docString;

  funcBody += `function ${name}(${input1} ${variableName1}, ${input2} ${variableName2}) internal returns (${returnType}) {`;

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
        ${UnderlyingTypes[returnType]} result = calcBinaryPlaceholderValueHash(unwrappedInput1, unwrappedInput2, FunctionId.${name});
        ITaskManager(TaskManager.TASK_MANAGER_ADDRESS).createTask(result, "${name}", unwrappedInput1, unwrappedInput2);
        return ${wrapType(returnType, "result")};`;
      } else {
        // TODO : What's the use case for this? still not removing till I figure this out completely
        funcBody += `
        ${returnType} result = calcBinaryPlaceholderValueHash(unwrappedInput1, unwrappedInput2, FunctionId.${name});
        ITaskManager(TaskManager.TASK_MANAGER_ADDRESS).createTask(result, "${name}", unwrappedInput1, unwrappedInput2);
        return result;
        }`;
      }
    } else if (name === SEALING_FUNCTION_NAME) {
      // **** Value 1 is encrypted, value 2 is bytes32 - this is basically reencrypt/wrapForUser
      funcBody += `
        if (!isInitialized(${variableName1})) {
            ${variableName1} = ${asEuintFuncName(input1)}(0);
        }
        ${UnderlyingTypes[input1]} unwrapped = ${unwrapType(input1,variableName1)};
        ITaskManager(TaskManager.TASK_MANAGER_ADDRESS).createSealOutputTask(unwrapped, publicKey);
        return "";`;
    } else if (name === SEALING_TYPED_FUNCTION_NAME) {
      const returnTypeClean = returnType.replace(" memory", "");
      funcBody += `
        return ${returnTypeClean}({ data: sealoutput(${variableName1}, ${variableName2}), utype: ${UintTypes[input1 as EUintType]} });`;
    }
  } else {
    // don't support input 1 is plaintext
    throw new Error("Unsupported plaintext input1");
  }

  funcBody += `\n    }`;

  return funcBody;
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

export function genAbiFile(abi: string) {
  return `
import { BaseContract } from 'ethers';

export interface EncryptedNumber {
  data: Uint8Array;
  securityZone: number;
}
${abi}
`;
}

export function SolTemplateDecrypt(input1: AllTypes, returnType: AllTypes) {
  // const decryptedType = asDecryptedType(returnType as EUintType);

  // TODO : Consider returning default value with optional overloading of "${getHalfValueFrom(decryptedType)}"
  // @param defaultValue default value to be returned on gas estimation
  let funcBody = `
    /// @notice Performs the async decrypt operation on a ciphertext
    /// @dev The decrypted output should be asynchronously handled by the IAsyncFHEReceiver implementation
    /// @param input1 the input ciphertext
    /// @return the input ciphertext
    function decrypt(${input1} input1) internal returns (${returnType}) {`;

  if (valueIsEncrypted(input1)) {
    // Get the proper function
    funcBody += `
        if (!isInitialized(input1)) {
            input1 = ${asEuintFuncName(input1)}(0);
        }
        ${UnderlyingTypes[input1]} unwrappedInput1 = ${unwrapType(
          input1,
          "input1"
        )};
        ITaskManager(TaskManager.TASK_MANAGER_ADDRESS).createDecryptTask(unwrappedInput1);
        return input1;
    }`;
  } else {
    throw new Error("unsupported decrypt function of 1 input that is not encrypted");
  }
  return funcBody;
}

export function SolTemplate1Arg(
  name: string,
  input1: AllTypes,
  returnType: AllTypes
) {
  let docString = `
    /// @notice Performs the ${name} operation on a ciphertext
    /// @dev Verifies that the input value matches a valid ciphertext.
    /// @param input1 the input ciphertext
    `;

  let returnStr = returnType === "none" ? `` : `returns (${returnType}) `;

  let funcBody = docString;

  funcBody += `function ${name}(${input1} input1) internal ${returnStr}{`;

  // TODO : Implement for FHE.req support
  if (!valueIsEncrypted(returnType)) {
    funcBody += `
      // TODO : Not Implemented
    }`;
    return funcBody;
  } else if (valueIsEncrypted(input1)) {
    funcBody += `
      if (!isInitialized(input1)) {
          input1 = ${asEuintFuncName(input1)}(0);
      }
      ${UnderlyingTypes[input1]} unwrappedInput1 = ${unwrapType(
        input1,
        "input1"
      )};
      uint256 ctHash = calcUnaryPlaceholderValueHash(unwrappedInput1, FunctionId.${name});
      ITaskManager(TaskManager.TASK_MANAGER_ADDRESS).createTask(ctHash, "${name}", unwrappedInput1, 0);
      return ${wrapType(returnType, "ctHash")};
    }`;
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
    function ${name}(${input1} input1, ${input2} input2, ${input3} input3) internal returns (${returnType}) {
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

function operatorFunctionName(funcName: string, forType: EUintType) {
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

  if (!unary) {
    docString += `/// @param rhs second input of type ${forType}\n`;
  }

  docString += `/// @return the result of the ${funcName}`;

  return `
    ${docString}
    function ${funcName}(${funcParams}) internal returns (${returnType}) {
        return FHE.${funcName}(${unaryParameters});
    }`;
};

export const CastBinding = (thisType: string, targetType: string) => {
  if (
    targetType === "eaddress" &&
    !AllowedTypesOnCastToEaddress.includes(thisType)
  ) {
    return ""; // skip unsupported cast
  }

  return `
    function to${shortenType(
      targetType
    )}(${thisType} value) internal returns (${targetType}) {
        return FHE.as${capitalize(targetType)}(value);
    }`;
};

export const SealFromType = (thisType: string) => {
  return `
    function ${LOCAL_SEAL_FUNCTION_NAME}(${thisType} value, bytes32 publicKey) internal returns (${SEAL_RETURN_TYPE} memory) {
        return FHE.sealoutput(value, publicKey);
    }`;
};

export const SealTypedFromType = (thisType: string) => {
  let returnType: string
  if (thisType === 'ebool') returnType = 'SealedBool'
  else if (thisType === 'eaddress') returnType = 'SealedAddress'
  else returnType = 'SealedUint'
  return `
    function ${LOCAL_SEAL_TYPED_FUNCTION_NAME}(${thisType} value, bytes32 publicKey) internal returns (${returnType} memory) {
        return FHE.sealoutputTyped(value, publicKey);
    }`;
};

export const getHalfValueFrom = (thisType: string) => {
  if (thisType === "bool") {
    return "false";
  }
  if (thisType === "address") {
    return "address(0)";
  }
  const match = thisType.match(/\d+/);
  const bits = match ? parseInt(match[0], 10) : 0;
  if (bits === 0) {
    return "0";
  }

  return `(2 ** ${bits}) / 2`;
};

export const DecryptBinding = (thisType: string) => {
  return `
    function ${LOCAL_DECRYPT_FUNCTION_NAME}(${thisType} value) internal returns (${thisType}) {
        return FHE.decrypt(value);
    }`;

    // TODO : Consider returning the defaultValue here
    // let plaintextType = toPlaintextType(thisType);
    // `function ${LOCAL_DECRYPT_FUNCTION_NAME}(${thisType} value, ${plaintextType} defaultValue) internal pure returns (${thisType}) {
        // return FHE.decrypt(value, defaultValue);
    // }`;
};


export const RandomGenericFunction = () => {
  return `
    /// @notice Generates a random value of a given type with the given seed, for the provided securityZone
    /// @dev Calls the desired precompile and returns the hash of the ciphertext
    /// @param uintType the type of the random value to generate
    /// @param seed the seed to use to create a random value from
    /// @param securityZone the security zone to use for the random value
    function random(uint8 uintType, uint64 seed, int32 securityZone) internal pure returns (uint256) {
        bytes memory b = FheOps(Precompiles.Fheos).random(uintType, seed, securityZone);
        return Impl.getValue(b);
    }
    /// @notice Generates a random value of a given type with the given seed
    /// @dev Calls the desired precompile and returns the hash of the ciphertext
    /// @param uintType the type of the random value to generate
    /// @param seed the seed to use to create a random value from
    function random(uint8 uintType, uint32 seed) internal pure returns (uint256) {
        return random(uintType, seed, 0);
    }
    /// @notice Generates a random value of a given type
    /// @dev Calls the desired precompile and returns the hash of the ciphertext
    /// @param uintType the type of the random value to generate
    function random(uint8 uintType) internal pure returns (uint256) {
        return random(uintType, 0, 0);
    }
    `;
};

const RandomFunctionForType = (type: string) => {
  if (type === "ebool" || type === "eaddress") {
    return "";
  }
  return `/// @notice Generates a random value of a ${type} type for provided securityZone
    /// @dev Calls the desired precompile and returns the hash of the ciphertext
    /// @param securityZone the security zone to use for the random value
    function random${capitalize(
    type
  )}(int32 securityZone) internal pure returns (${type}) {
        uint256 result = random(Common.${type.toUpperCase()}_TFHE, 0, securityZone);
        return ${type}.wrap(result);
    }
    /// @notice Generates a random value of a ${type} type
    /// @dev Calls the desired precompile and returns the hash of the ciphertext
    function random${capitalize(
    type
  )}() internal pure returns (${type}) {
        return random${capitalize(type)}(0);
    }
    `;
};

export const RandomFunctions = () => {
  let outputFile = "";
  for (let type of EInputType) {
    outputFile += RandomFunctionForType(type);
  }

  return outputFile;
};
