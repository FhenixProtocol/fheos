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
  toPlaintextType,
  capitalize,
  shortenType,
} from "../common";

export const preamble = () => {
  return `// SPDX-License-Identifier: BSD-3-Clause-Clear
// solhint-disable one-contract-per-file

pragma solidity >=0.8.19 <=0.8.25;

import {Precompiles, FheOps} from "./FheOS.sol";

${EInputType.map((type) => {
  return `type ${type} is uint256;`;
}).join("\n")}

${EInputType.map((type) => {
  return `struct in${capitalize(type)} {
    bytes data;
}`;
}).join("\n")}

struct SealedArray {
  bytes[] data;
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
}

library Impl {
    function sealoutput(uint8 utype, uint256 ciphertext, bytes32 publicKey) internal pure returns (${SEAL_RETURN_TYPE} memory reencrypted) {
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
  if (!valueIsEncrypted(fromType)) {
    return ""; // casting from plaintext type is handled elsewhere
  }
  return `Impl.cast(${
    UintTypes[fromType]
  }, ${fromType}.unwrap(${name}), Common.${toType.toUpperCase()}_TFHE)`;
};

const castFromPlaintext = (name: string, toType: string): string => {
  return `Impl.trivialEncrypt(${name}, Common.${toType.toUpperCase()}_TFHE)`;
};

const castFromAddress = (name: string, toType: string): string => {
  return `Impl.trivialEncrypt(uint256(uint160(${name})), Common.${toType.toUpperCase()}_TFHE)`;
};

const castFromBytes = (name: string, toType: string): string => {
  return `Impl.verify(${name}, Common.${toType.toUpperCase()}_TFHE)`;
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
  } else if (fromType == "address" && toType == "eaddress") {
    docString += `
    /// Allows for a better user experience when working with eaddresses`;
    castString = castFromAddress("value", toType);
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

const unwrapType = (typeName: EUintType, inputName: string): string =>
  `${typeName}.unwrap(${inputName})`;
const wrapType = (resultType: EUintType, inputName: string): string =>
  `${resultType}.wrap(${inputName})`;
const asEuintFuncName = (typeName: EUintType): string =>
  `as${capitalize(typeName)}`;

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

export const CastBinding = (thisType: string, targetType: string) => {
  return `
    function to${shortenType(
      targetType
    )}(${thisType} value) internal pure returns (${targetType}) {
        return FHE.as${capitalize(targetType)}(value);
    }`;
};

export const SealFromType = (thisType: string) => {
  return `
    function ${LOCAL_SEAL_FUNCTION_NAME}(${thisType} value, bytes32 publicKey) internal pure returns (${SEAL_RETURN_TYPE} memory) {
        return FHE.sealoutput(value, publicKey);
    }`;
};

export const DecryptBinding = (thisType: string) => {
  return `
    function ${LOCAL_DECRYPT_FUNCTION_NAME}(${thisType} value) internal pure returns (${toPlaintextType(
    thisType
  )}) {
        return FHE.decrypt(value);
    }`;
};
