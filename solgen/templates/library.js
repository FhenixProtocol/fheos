"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DecryptBinding = exports.SealFromType = exports.CastBinding = exports.OperatorBinding = exports.BindingLibraryType = exports.OperatorOverloadDecl = exports.SolTemplate3Arg = exports.SolTemplate1Arg = exports.genAbiFile = exports.testContract2Arg = exports.IsOperationAllowed = exports.testContract3Arg = exports.testContractReencrypt = exports.testContractReq = exports.generateTestContract = exports.testContract1Arg = exports.testContract2ArgBoolRes = exports.SolTemplate2Arg = exports.capitalize = exports.AsTypeTestingContract = exports.AsTypeFunction = exports.PostFix = exports.preamble = void 0;
var common_1 = require("../common");
var preamble = function () {
    return "// SPDX-License-Identifier: BSD-3-Clause-Clear\n// solhint-disable one-contract-per-file\n\npragma solidity >=0.8.19 <=0.8.25;\n\nimport {Precompiles, FheOps} from \"./FheOS.sol\";\n\n".concat(common_1.EInputType.map(function (type) {
        return "type ".concat(type, " is uint256;");
    }).join("\n"), "\n\n").concat(common_1.EInputType.map(function (type) {
        return "struct in".concat((0, exports.capitalize)(type), " {\n    bytes data;\n}");
    }).join("\n"), "\n\nstruct SealedArray {\n  bytes[] data;\n}\n\nlibrary Common {\n    // Values used to communicate types to the runtime.\n    // Must match values defined in warp-drive protobufs for everything to \n    // make sense\n    uint8 internal constant EUINT8_TFHE = 0;\n    uint8 internal constant EUINT16_TFHE = 1;\n    uint8 internal constant EUINT32_TFHE = 2;\n    uint8 internal constant EUINT64_TFHE = 3;\n    uint8 internal constant EUINT128_TFHE = 4;\n    uint8 internal constant EUINT256_TFHE = 5;\n    uint8 internal constant EADDRESS_TFHE = 12;\n    // uint8 internal constant INT_BGV = 12;\n    uint8 internal constant EBOOL_TFHE = 13;\n    \n    function bigIntToBool(uint256 i) internal pure returns (bool) {\n        return (i > 0);\n    }\n\n    function bigIntToUint8(uint256 i) internal pure returns (uint8) {\n        return uint8(i);\n    }\n\n    function bigIntToUint16(uint256 i) internal pure returns (uint16) {\n        return uint16(i);\n    }\n\n    function bigIntToUint32(uint256 i) internal pure returns (uint32) {\n        return uint32(i);\n    }\n\n    function bigIntToUint64(uint256 i) internal pure returns (uint64) {\n        return uint64(i);\n    }\n\n    function bigIntToUint128(uint256 i) internal pure returns (uint128) {\n        return uint128(i);\n    }\n\n    function bigIntToUint256(uint256 i) internal pure returns (uint256) {\n        return i;\n    }\n\n    function bigIntToAddress(uint256 i) internal pure returns (address) {\n      return address(uint160(i));\n    }\n    \n    function toBytes(uint256 x) internal pure returns (bytes memory b) {\n        b = new bytes(32);\n        assembly { mstore(add(b, 32), x) }\n    }\n}\n\nlibrary Impl {\n    function sealoutput(uint8 utype, uint256 ciphertext, bytes32 publicKey) internal pure returns (").concat(common_1.SEAL_RETURN_TYPE, " memory reencrypted) {\n        // Call the sealoutput precompile.\n        reencrypted = FheOps(Precompiles.Fheos).sealOutput(utype, Common.toBytes(ciphertext), bytes.concat(publicKey));\n\n        return reencrypted;\n    }\n\n    function verify(bytes memory _ciphertextBytes, uint8 _toType) internal pure returns (uint256 result) {\n        bytes memory output;\n\n        // Call the verify precompile.\n        output = FheOps(Precompiles.Fheos).verify(_toType, _ciphertextBytes);\n        result = getValue(output);\n    }\n\n    function cast(uint8 utype, uint256 ciphertext, uint8 toType) internal pure returns (uint256 result) {\n        bytes memory output;\n\n        // Call the cast precompile.\n        output = FheOps(Precompiles.Fheos).cast(utype, Common.toBytes(ciphertext), toType);\n        result = getValue(output);\n    }\n\n    function getValue(bytes memory a) internal pure returns (uint256 value) {\n        assembly {\n            value := mload(add(a, 0x20))\n        }\n    }\n\n    function trivialEncrypt(uint256 value, uint8 toType) internal pure returns (uint256 result) {\n        bytes memory output;\n\n        // Call the trivialEncrypt precompile.\n        output = FheOps(Precompiles.Fheos).trivialEncrypt(Common.toBytes(value), toType);\n\n        result = getValue(output);\n    }\n\n    function select(uint8 utype, uint256 control, uint256 ifTrue, uint256 ifFalse) internal pure returns (uint256 result) {\n        bytes memory output;\n\n        // Call the trivialEncrypt precompile.\n        output = FheOps(Precompiles.Fheos).select(utype, Common.toBytes(control), Common.toBytes(ifTrue), Common.toBytes(ifFalse));\n\n        result = getValue(output);\n    }\n}\n\nlibrary FHE {\n    euint8 public constant NIL8 = euint8.wrap(0);\n    euint16 public constant NIL16 = euint16.wrap(0);\n    euint32 public constant NIL32 = euint32.wrap(0);\n\n    // Return true if the encrypted integer is initialized and false otherwise.\n    function isInitialized(ebool v) internal pure returns (bool) {\n        return ebool.unwrap(v) != 0;\n    }\n\n    // Return true if the encrypted integer is initialized and false otherwise.\n    function isInitialized(euint8 v) internal pure returns (bool) {\n        return euint8.unwrap(v) != 0;\n    }\n\n    // Return true if the encrypted integer is initialized and false otherwise.\n    function isInitialized(euint16 v) internal pure returns (bool) {\n        return euint16.unwrap(v) != 0;\n    }\n\n    // Return true if the encrypted integer is initialized and false otherwise.\n    function isInitialized(euint32 v) internal pure returns (bool) {\n        return euint32.unwrap(v) != 0;\n    }\n    \n    // Return true if the encrypted integer is initialized and false otherwise.\n    function isInitialized(euint64 v) internal pure returns (bool) {\n        return euint64.unwrap(v) != 0;\n    }\n    \n        // Return true if the encrypted integer is initialized and false otherwise.\n    function isInitialized(euint128 v) internal pure returns (bool) {\n        return euint128.unwrap(v) != 0;\n    }\n    \n        // Return true if the encrypted integer is initialized and false otherwise.\n    function isInitialized(euint256 v) internal pure returns (bool) {\n        return euint256.unwrap(v) != 0;\n    }\n\n    function isInitialized(eaddress v) internal pure returns (bool) {\n        return eaddress.unwrap(v) != 0;\n    }\n\n    function getValue(bytes memory a) private pure returns (uint256 value) {\n        assembly {\n            value := mload(add(a, 0x20))\n        }\n    }\n    \n    function mathHelper(\n        uint8 utype,\n        uint256 lhs,\n        uint256 rhs,\n        function(uint8, bytes memory, bytes memory) external pure returns (bytes memory) impl\n    ) internal pure returns (uint256 result) {\n        bytes memory output;\n        output = impl(utype, Common.toBytes(lhs), Common.toBytes(rhs));\n        result = getValue(output);\n    }\n    ");
};
exports.preamble = preamble;
var PostFix = function () {
    return "\n}";
};
exports.PostFix = PostFix;
var castFromEncrypted = function (fromType, toType, name) {
    if (!(0, common_1.valueIsEncrypted)(toType)) {
        console.log("Unsupported type for casting: ".concat(toType));
        process.exit(1);
    }
    if (!(0, common_1.valueIsEncrypted)(fromType)) {
        return ""; // casting from plaintext type is handled elsewhere
    }
    return "Impl.cast(".concat(common_1.UintTypes[fromType], ", ").concat(fromType, ".unwrap(").concat(name, "), Common.").concat(toType.toUpperCase(), "_TFHE)");
};
var castFromPlaintext = function (name, toType) {
    return "Impl.trivialEncrypt(".concat(name, ", Common.").concat(toType.toUpperCase(), "_TFHE)");
};
var castFromAddress = function (name, toType) {
    return "Impl.trivialEncrypt(uint256(uint160(".concat(name, ")), Common.").concat(toType.toUpperCase(), "_TFHE)");
};
var castFromBytes = function (name, toType) {
    return "Impl.verify(".concat(name, ", Common.").concat(toType.toUpperCase(), "_TFHE)");
};
var castFromInputType = function (name, toType) {
    return "FHE.as".concat((0, exports.capitalize)(toType), "(").concat(name, ".data)");
};
var castToEbool = function (name, fromType) {
    return "\n    \n    /// @notice Converts a ".concat(fromType, " to an ebool\n    function asEbool(").concat(fromType, " value) internal pure returns (ebool) {\n        return ne(").concat(name, ", as").concat((0, exports.capitalize)(fromType), "(0));\n    }");
};
var AsTypeFunction = function (fromType, toType) {
    var castString = castFromEncrypted(fromType, toType, "value");
    var docString = "\n    /// @notice Converts a ".concat(fromType, " to an ").concat(toType);
    if (fromType === "bool" && toType === "ebool") {
        return "\n    /// @notice Converts a plaintext boolean value to a ciphertext ebool\n    /// @dev Privacy: The input value is public, therefore the ciphertext should be considered public and should be used\n    ///only for mathematical operations, not to represent data that should be private\n    /// @return A ciphertext representation of the input \n    function asEbool(bool value) internal pure returns (ebool) {\n        uint256 sVal = 0;\n        if (value) {\n            sVal = 1;\n        }\n\n        return asEbool(sVal);\n    }";
    }
    else if (fromType.startsWith("in")) {
        docString = "\n    /// @notice Parses input ciphertexts from the user. Converts from encrypted raw bytes to an ".concat(toType, "\n    /// @dev Also performs validation that the ciphertext is valid and has been encrypted using the network encryption key\n    /// @return a ciphertext representation of the input");
        castString = castFromInputType("value", toType);
        return "".concat(docString, "\n    function as").concat((0, exports.capitalize)(toType), "(").concat(fromType, " memory value) internal pure returns (").concat(toType, ") {\n        return ").concat(castString, ";\n    }");
    }
    else if (fromType === "bytes memory") {
        docString = "\n    /// @notice Parses input ciphertexts from the user. Converts from encrypted raw bytes to an ".concat(toType, "\n    /// @dev Also performs validation that the ciphertext is valid and has been encrypted using the network encryption key\n    /// @return a ciphertext representation of the input");
        castString = castFromBytes("value", toType);
    }
    else if (fromType == "address" && toType == "eaddress") {
        docString += "\n    /// Allows for a better user experience when working with eaddresses";
        castString = castFromAddress("value", toType);
    }
    else if (common_1.EPlaintextType.includes(fromType)) {
        castString = castFromPlaintext("value", toType);
    }
    else if (toType === "ebool") {
        return castToEbool("value", fromType);
    }
    else if (!common_1.EInputType.includes(fromType)) {
        throw new Error("Unsupported type for casting: ".concat(fromType));
    }
    return "".concat(docString, "\n    function as").concat((0, exports.capitalize)(toType), "(").concat(fromType, " value) internal pure returns (").concat(toType, ") {\n        return ").concat(toType, ".wrap(").concat(castString, ");\n    }");
};
exports.AsTypeFunction = AsTypeFunction;
function TypeCastTestingFunction(fromType, fromTypeForTs, toType, fromTypeEncrypted) {
    var to = (0, exports.capitalize)(toType);
    var retType = to.slice(1);
    var testType = fromTypeEncrypted ? fromTypeEncrypted : fromType;
    testType =
        testType === "bytes memory" ? "PreEncrypted" : (0, exports.capitalize)(testType);
    testType = testType === "Uint256" ? "Plaintext" : testType;
    var encryptedVal = fromTypeEncrypted
        ? "FHE.as".concat((0, exports.capitalize)(fromTypeEncrypted), "(val)")
        : "val";
    var retTypeTs = retType === "bool" ? "boolean" : retType;
    retTypeTs = retTypeTs.includes("uint") || retTypeTs.includes("address") ? "bigint" : retTypeTs;
    var abi;
    var func = "\n\n    ";
    if (testType === "PreEncrypted" || testType === "Plaintext") {
        func += "function castFrom".concat(testType, "To").concat(to, "(").concat(fromType, " val) public pure returns (").concat(retType, ") {\n        return FHE.decrypt(FHE.as").concat(to, "(").concat(encryptedVal, "));\n    }");
        abi = "    castFrom".concat(testType, "To").concat(to, ": (val: ").concat(fromTypeForTs, ") => Promise<").concat(retTypeTs, ">;\n");
    }
    else {
        func += "function castFrom".concat(testType, "To").concat(to, "(").concat(fromType, " val, string calldata test) public pure returns (").concat(retType, ") {\n        if (Utils.cmp(test, \"bound\")) {\n            return ").concat(encryptedVal, ".to").concat(shortenType(toType), "().decrypt();\n        } else if (Utils.cmp(test, \"regular\")) {\n            return FHE.decrypt(FHE.as").concat(to, "(").concat(encryptedVal, "));\n        }\n        revert TestNotFound(test);\n    }");
        abi = "    castFrom".concat(testType, "To").concat(to, ": (val: ").concat(fromTypeForTs, ", test: string) => Promise<").concat(retTypeTs, ">;\n");
    }
    return [func, abi];
}
function AsTypeTestingContract(type) {
    var funcs = "";
    var abi = "export interface As".concat((0, exports.capitalize)(type), "TestType extends BaseContract {\n");
    // Although casts from eaddress to types with < 256 bits are possible, we don't want to test them.
    var eaddressAllowedTypes = ["euint256", "uint256", "bytes memory"];
    var fromTypeCollection = type === "eaddress" ? eaddressAllowedTypes : common_1.EInputType.concat("uint256", "bytes memory");
    for (var _i = 0, fromTypeCollection_1 = fromTypeCollection; _i < fromTypeCollection_1.length; _i++) {
        var fromType = fromTypeCollection_1[_i];
        if (type === fromType || (fromType === "eaddress" && !eaddressAllowedTypes.includes(type))) {
            continue;
        }
        var fromTypeTs = fromType === "bytes memory" ? "Uint8Array" : "bigint";
        var fromTypeSol = fromType === "bytes memory" ? fromType : "uint256";
        var fromTypeEncrypted = common_1.EInputType.includes(fromType)
            ? fromType
            : undefined;
        var contractInfo = TypeCastTestingFunction(fromTypeSol, fromTypeTs, type, fromTypeEncrypted);
        funcs += contractInfo[0];
        abi += contractInfo[1];
    }
    funcs = funcs.slice(1);
    abi += "}\n";
    return [generateTestContract("As".concat((0, exports.capitalize)(type)), funcs), abi];
}
exports.AsTypeTestingContract = AsTypeTestingContract;
var unwrapType = function (typeName, inputName) {
    return "".concat(typeName, ".unwrap(").concat(inputName, ")");
};
var wrapType = function (resultType, inputName) {
    return "".concat(resultType, ".wrap(").concat(inputName, ")");
};
var asEuintFuncName = function (typeName) {
    return "as".concat((0, exports.capitalize)(typeName));
};
var capitalize = function (s) { return s.charAt(0).toUpperCase() + s.slice(1); };
exports.capitalize = capitalize;
function SolTemplate2Arg(name, input1, input2, returnType) {
    // special names for reencrypt function (don't check name === reencrypt) because the name could change
    var variableName1 = input2 === "bytes32" ? "value" : "lhs";
    var variableName2 = input2 === "bytes32" ? "publicKey" : "rhs";
    var docString = "\n    /// @notice This functions performs the ".concat(name, " operation\n    /// @dev If any of the inputs are expected to be a ciphertext, it verifies that the value matches a valid ciphertext\n    ///Pure in this function is marked as a hack/workaround - note that this function is NOT pure as fetches of ciphertexts require state access\n    /// @param lhs The first input \n    /// @param rhs The second input\n    /// @return The result of the operation\n    ");
    // reencrypt
    if (variableName2 === "publicKey") {
        docString = "\n    /// @notice performs the ".concat(name, " function on a ").concat(input1, " ciphertext. This operation returns the plaintext value, sealed for the public key provided \n    /// @dev Pure in this function is marked as a hack/workaround - note that this function is NOT pure as fetches of ciphertexts require state access\n    /// @param value Ciphertext to decrypt and seal\n    /// @param publicKey Public Key that will receive the sealed plaintext\n    /// @return Plaintext input, sealed for the owner of `publicKey`\n    ");
    }
    var funcBody = docString;
    funcBody += "function ".concat(name, "(").concat(input1, " ").concat(variableName1, ", ").concat(input2, " ").concat(variableName2, ") internal pure returns (").concat(returnType, ") {");
    if ((0, common_1.valueIsEncrypted)(input1)) {
        // both inputs encrypted - this is a generic math function. i.e. div, mul, eq, etc.
        // 1. possibly cast input1
        // 2. possibly cast input2
        // 3. possibly cast return type
        if ((0, common_1.valueIsEncrypted)(input2) && input1 !== input2) {
            return "";
        }
        if ((0, common_1.valueIsEncrypted)(input2) || common_1.EPlaintextType.includes(input2)) {
            var input2Cast = input1 === input2
                ? variableName2
                : "".concat(asEuintFuncName(input1), "(").concat(variableName2, ")");
            //
            funcBody += "\n        if (!isInitialized(".concat(variableName1, ")) {\n            ").concat(variableName1, " = ").concat(asEuintFuncName(input1), "(0);\n        }\n        if (!isInitialized(").concat(variableName2, ")) {\n            ").concat(variableName2, " = ").concat(asEuintFuncName(input1), "(0);\n        }\n        ").concat(common_1.UnderlyingTypes[input1], " unwrappedInput1 = ").concat(unwrapType(input1, variableName1), ";\n        ").concat(common_1.UnderlyingTypes[input1], " unwrappedInput2 = ").concat(unwrapType(input1, "".concat(input2Cast)), ";\n");
            if ((0, common_1.valueIsEncrypted)(returnType)) {
                funcBody += "\n        ".concat(common_1.UnderlyingTypes[returnType], " result = mathHelper(").concat(common_1.UintTypes[input1], ", unwrappedInput1, unwrappedInput2, FheOps(Precompiles.Fheos).").concat(name, ");\n        return ").concat(wrapType(returnType, "result"), ";");
            }
            else {
                funcBody += "\n        ".concat(returnType, " result = mathHelper(").concat(common_1.UintTypes[input1], ", unwrappedInput1, unwrappedInput2, FheOps(Precompiles.Fheos).").concat(name, ");\n        return result;");
            }
        }
        else if (input2 === "bytes32") {
            // **** Value 1 is encrypted, value 2 is bytes32 - this is basically reencrypt/wrapForUser
            funcBody += "\n        if (!isInitialized(".concat(variableName1, ")) {\n            ").concat(variableName1, " = ").concat(asEuintFuncName(input1), "(0);\n        }\n        ").concat(common_1.UnderlyingTypes[input1], " unwrapped = ").concat(unwrapType(input1, variableName1), ";\n\n        return Impl.").concat(name, "(").concat(common_1.UintTypes[input1], ", unwrapped, ").concat(variableName2, ");");
        }
    }
    else {
        // don't support input 1 is plaintext
        throw new Error("Unsupported plaintext input1");
    }
    funcBody += "\n    }";
    return funcBody;
}
exports.SolTemplate2Arg = SolTemplate2Arg;
function testContract2ArgBoolRes(name, isBoolean) {
    var isEuint64Allowed = (0, exports.IsOperationAllowed)(name, common_1.EInputType.indexOf("euint64"));
    var isEuint128Allowed = (0, exports.IsOperationAllowed)(name, common_1.EInputType.indexOf("euint128"));
    var isEuint256Allowed = (0, exports.IsOperationAllowed)(name, common_1.EInputType.indexOf("euint256"));
    var isEaddressAllowed = (0, exports.IsOperationAllowed)(name, common_1.EInputType.indexOf("eaddress"));
    var func = "function ".concat(name, "(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {\n        if (Utils.cmp(test, \"").concat(name, "(euint8,euint8)\")) {\n            if (FHE.decrypt(FHE.").concat(name, "(FHE.asEuint8(a), FHE.asEuint8(b)))) {\n                return 1;\n            }\n\n            return 0;\n        } else if (Utils.cmp(test, \"").concat(name, "(euint16,euint16)\")) {\n            if (FHE.decrypt(FHE.").concat(name, "(FHE.asEuint16(a), FHE.asEuint16(b)))) {\n                return 1;\n            }\n\n            return 0;\n        } else if (Utils.cmp(test, \"").concat(name, "(euint32,euint32)\")) {\n            if (FHE.decrypt(FHE.").concat(name, "(FHE.asEuint32(a), FHE.asEuint32(b)))) {\n                return 1;\n            }\n\n            return 0;\n        }");
    if (isEuint64Allowed) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(euint64,euint64)\")) {\n            if (FHE.decrypt(FHE.").concat(name, "(FHE.asEuint64(a), FHE.asEuint64(b)))) {\n                return 1;\n            }\n\n            return 0;\n        }");
    }
    if (isEuint128Allowed) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(euint128,euint128)\")) {\n            if (FHE.decrypt(FHE.").concat(name, "(FHE.asEuint128(a), FHE.asEuint128(b)))) {\n                return 1;\n            }\n\n            return 0;\n        }");
    }
    if (isEuint256Allowed) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(euint256,euint256)\")) {\n            if (FHE.decrypt(FHE.").concat(name, "(FHE.asEuint256(a), FHE.asEuint256(b)))) {\n                return 1;\n            }\n\n            return 0;\n        }");
    }
    if (isEaddressAllowed) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(eaddress,eaddress)\")) {\n            if (FHE.decrypt(FHE.").concat(name, "(FHE.asEaddress(a), FHE.asEaddress(b)))) {\n                return 1;\n            }\n\n            return 0;\n        }");
    }
    func += " else if (Utils.cmp(test, \"euint8.".concat(name, "(euint8)\")) {\n            if (FHE.asEuint8(a).").concat(name, "(FHE.asEuint8(b)).decrypt()) {\n                return 1;\n            }\n\n            return 0;\n        } else if (Utils.cmp(test, \"euint16.").concat(name, "(euint16)\")) {\n            if (FHE.asEuint16(a).").concat(name, "(FHE.asEuint16(b)).decrypt()) {\n                return 1;\n            }\n\n            return 0;\n        } else if (Utils.cmp(test, \"euint32.").concat(name, "(euint32)\")) {\n            if (FHE.asEuint32(a).").concat(name, "(FHE.asEuint32(b)).decrypt()) {\n                return 1;\n            }\n\n            return 0;\n        }");
    if (isEuint64Allowed) {
        func += " else if (Utils.cmp(test, \"euint64.".concat(name, "(euint64)\")) {\n            if (FHE.asEuint64(a).").concat(name, "(FHE.asEuint64(b)).decrypt()) {\n                return 1;\n            }\n            return 0;\n        }");
    }
    if (isEuint128Allowed) {
        func += " else if (Utils.cmp(test, \"euint128.".concat(name, "(euint128)\")) {\n            if (FHE.asEuint128(a).").concat(name, "(FHE.asEuint128(b)).decrypt()) {\n                return 1;\n            }\n            return 0;\n        }");
    }
    if (isEuint256Allowed) {
        func += " else if (Utils.cmp(test, \"euint256.".concat(name, "(euint256)\")) {\n            if (FHE.asEuint256(a).").concat(name, "(FHE.asEuint256(b)).decrypt()) {\n                return 1;\n            }\n            return 0;\n        }");
    }
    if (isEaddressAllowed) {
        func += " else if (Utils.cmp(test, \"eaddress.".concat(name, "(eaddress)\")) {\n            if (FHE.asEaddress(a).").concat(name, "(FHE.asEaddress(b)).decrypt()) {\n                return 1;\n            }\n            return 0;\n        }");
    }
    if (isBoolean) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(ebool,ebool)\")) {\n            bool aBool = true;\n            bool bBool = true;\n            if (a == 0) {\n                aBool = false;\n            }\n            if (b == 0) {\n                bBool = false;\n            }\n            if (FHE.decrypt(FHE.").concat(name, "(FHE.asEbool(aBool), FHE.asEbool(bBool)))) {\n                return 1;\n            }\n\n            return 0;\n        } else if (Utils.cmp(test, \"ebool.").concat(name, "(ebool)\")) {\n            bool aBool = true;\n            bool bBool = true;\n            if (a == 0) {\n                aBool = false;\n            }\n            if (b == 0) {\n                bBool = false;\n            }\n            if (FHE.asEbool(aBool).").concat(name, "(FHE.asEbool(bBool)).decrypt()) {\n                return 1;\n            }\n            return 0;\n        }");
    }
    func += "\n        revert TestNotFound(test);\n    }";
    var abi = "export interface ".concat((0, exports.capitalize)(name), "TestType extends BaseContract {\n    ").concat(name, ": (test: string, a: bigint, b: bigint) => Promise<bigint>;\n}\n");
    return [generateTestContract(name, func), abi];
}
exports.testContract2ArgBoolRes = testContract2ArgBoolRes;
function testContract1Arg(name) {
    var isEuint64Allowed = (0, exports.IsOperationAllowed)(name, common_1.EInputType.indexOf("euint64"));
    var isEuint128Allowed = (0, exports.IsOperationAllowed)(name, common_1.EInputType.indexOf("euint128"));
    var isEuint256Allowed = (0, exports.IsOperationAllowed)(name, common_1.EInputType.indexOf("euint256"));
    var func = "function ".concat(name, "(string calldata test, uint256 a) public pure returns (uint256 output) {\n        if (Utils.cmp(test, \"").concat(name, "(euint8)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint8(a)));\n        } else if (Utils.cmp(test, \"").concat(name, "(euint16)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint16(a)));\n        } else if (Utils.cmp(test, \"").concat(name, "(euint32)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint32(a)));\n        }");
    if (isEuint64Allowed) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(euint64)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint64(a)));\n        }");
    }
    if (isEuint128Allowed) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(euint128)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint128(a)));\n        }");
    }
    if (isEuint256Allowed) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(euint256)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint256(a)));\n        }");
    }
    func += " else if (Utils.cmp(test, \"".concat(name, "(ebool)\")) {\n            bool aBool = true;\n            if (a == 0) {\n                aBool = false;\n            }\n\n            if (FHE.decrypt(FHE.").concat(name, "(FHE.asEbool(aBool)))) {\n                return 1;\n            }\n\n            return 0;\n        }\n        \n        revert TestNotFound(test);\n    }");
    var abi = "export interface ".concat((0, exports.capitalize)(name), "TestType extends BaseContract {\n    ").concat(name, ": (test: string, a: bigint) => Promise<bigint>;\n}\n");
    return [generateTestContract(name, func), abi];
}
exports.testContract1Arg = testContract1Arg;
function generateTestContract(name, testFunc, importTypes) {
    if (importTypes === void 0) { importTypes = false; }
    var importStatement = importTypes
        ? "\nimport {ebool, euint8} from \"../../FHE.sol\";"
        : "";
    return "// SPDX-License-Identifier: MIT\npragma solidity ^0.8.17;\n\nimport {FHE} from \"../../FHE.sol\";".concat(importStatement, "\nimport {Utils} from \"./utils/Utils.sol\";\n\nerror TestNotFound(string test);\n\ncontract ").concat((0, exports.capitalize)(name), "Test {\n    using Utils for *;\n    \n    ").concat(testFunc, "\n}\n");
}
exports.generateTestContract = generateTestContract;
function testContractReq() {
    // Req is failing on EthCall so we need to make it as tx for now
    var func = "function req(string calldata test, uint256 a) public {\n        if (Utils.cmp(test, \"req(euint8)\")) {\n            FHE.req(FHE.asEuint8(a));\n        } else if (Utils.cmp(test, \"req(euint16)\")) {\n            FHE.req(FHE.asEuint16(a));\n        } else if (Utils.cmp(test, \"req(euint32)\")) {\n            FHE.req(FHE.asEuint32(a));\n        } else if (Utils.cmp(test, \"req(euint64)\")) {\n            FHE.req(FHE.asEuint64(a));\n        } else if (Utils.cmp(test, \"req(euint128)\")) {\n            FHE.req(FHE.asEuint128(a));\n        } else if (Utils.cmp(test, \"req(euint256)\")) {\n            FHE.req(FHE.asEuint256(a));\n        } else if (Utils.cmp(test, \"req(ebool)\")) {\n            bool b = true;\n            if (a == 0) {\n                b = false;\n            }\n            FHE.req(FHE.asEbool(b));\n        } else {\n            revert TestNotFound(test);\n        }\n    }";
    var abi = "export interface ReqTestType extends BaseContract {\n    req: (test: string, a: bigint) => Promise<{}>;\n}\n";
    return [generateTestContract("req", func), abi];
}
exports.testContractReq = testContractReq;
function testContractReencrypt() {
    var func = "function ".concat(common_1.SEALING_FUNCTION_NAME, "(string calldata test, uint256 a, bytes32 pubkey) public pure returns (").concat(common_1.SEAL_RETURN_TYPE, " memory reencrypted) {\n        if (Utils.cmp(test, \"").concat(common_1.SEALING_FUNCTION_NAME, "(euint8)\")) {\n            return FHE.").concat(common_1.SEALING_FUNCTION_NAME, "(FHE.asEuint8(a), pubkey);\n        } else if (Utils.cmp(test, \"").concat(common_1.SEALING_FUNCTION_NAME, "(euint16)\")) {\n            return FHE.").concat(common_1.SEALING_FUNCTION_NAME, "(FHE.asEuint16(a), pubkey);\n        } else if (Utils.cmp(test, \"").concat(common_1.SEALING_FUNCTION_NAME, "(euint32)\")) {\n            return FHE.").concat(common_1.SEALING_FUNCTION_NAME, "(FHE.asEuint32(a), pubkey);\n        } else if (Utils.cmp(test, \"").concat(common_1.SEALING_FUNCTION_NAME, "(euint64)\")) {\n            return FHE.").concat(common_1.SEALING_FUNCTION_NAME, "(FHE.asEuint64(a), pubkey);\n        } else if (Utils.cmp(test, \"").concat(common_1.SEALING_FUNCTION_NAME, "(euint128)\")) {\n            return FHE.").concat(common_1.SEALING_FUNCTION_NAME, "(FHE.asEuint128(a), pubkey);\n        } else if (Utils.cmp(test, \"").concat(common_1.SEALING_FUNCTION_NAME, "(euint256)\")) {\n            return FHE.").concat(common_1.SEALING_FUNCTION_NAME, "(FHE.asEuint256(a), pubkey);\n        } else if (Utils.cmp(test, \"").concat(common_1.SEALING_FUNCTION_NAME, "(ebool)\")) {\n            bool b = true;\n            if (a == 0) {\n                b = false;\n            }\n\n            return FHE.").concat(common_1.SEALING_FUNCTION_NAME, "(FHE.asEbool(b), pubkey);\n        } else if (Utils.cmp(test, \"").concat(common_1.LOCAL_SEAL_FUNCTION_NAME, "(euint8)\")) {\n            euint8 aEnc = FHE.asEuint8(a);\n            return aEnc.").concat(common_1.LOCAL_SEAL_FUNCTION_NAME, "(pubkey);\n        }\n        revert TestNotFound(test);\n    }");
    var abi = "export interface SealoutputTestType extends BaseContract {\n    ".concat(common_1.SEALING_FUNCTION_NAME, ": (test: string, a: bigint, pubkey: Uint8Array) => Promise<string>;\n}\n");
    return [generateTestContract(common_1.SEALING_FUNCTION_NAME, func, true), abi];
}
exports.testContractReencrypt = testContractReencrypt;
function testContract3Arg(name) {
    var func = "function ".concat(name, "(string calldata test, bool c, uint256 a, uint256 b) public pure returns (uint256 output) {\n        ebool condition = FHE.asEbool(c);\n        if (Utils.cmp(test, \"").concat(name, ": euint8\")) {\n            return FHE.decrypt(FHE.").concat(name, "(condition, FHE.asEuint8(a), FHE.asEuint8(b)));\n        } else if (Utils.cmp(test, \"").concat(name, ": euint16\")) {\n            return FHE.decrypt(FHE.").concat(name, "(condition, FHE.asEuint16(a), FHE.asEuint16(b)));\n        } else if (Utils.cmp(test, \"").concat(name, ": euint32\")) {\n            return FHE.decrypt(FHE.").concat(name, "(condition, FHE.asEuint32(a), FHE.asEuint32(b)));\n        } else if (Utils.cmp(test, \"").concat(name, ": euint64\")) {\n            return FHE.decrypt(FHE.").concat(name, "(condition, FHE.asEuint64(a), FHE.asEuint64(b)));\n        } else if (Utils.cmp(test, \"").concat(name, ": euint128\")) {\n            return FHE.decrypt(FHE.").concat(name, "(condition, FHE.asEuint128(a), FHE.asEuint128(b)));\n        } else if (Utils.cmp(test, \"").concat(name, ": ebool\")) {\n            bool aBool = true;\n            bool bBool = true;\n            if (a == 0) {\n                aBool = false;\n            }\n            if (b == 0) {\n                bBool = false;\n            }\n\n            if(FHE.decrypt(FHE.").concat(name, "(condition, FHE.asEbool(aBool), FHE.asEbool(bBool)))) {\n                return 1;\n            }\n            return 0;\n        } \n        \n        revert TestNotFound(test);\n    }");
    var abi = "export interface ".concat((0, exports.capitalize)(name), "TestType extends BaseContract {\n    ").concat(name, ": (test: string, c: boolean, a: bigint, b: bigint) => Promise<bigint>;\n}\n");
    return [generateTestContract(name, func, true), abi];
}
exports.testContract3Arg = testContract3Arg;
var IsOperationAllowed = function (functionName, inputIdx) {
    var regexes = common_1.AllowedOperations[inputIdx];
    for (var _i = 0, regexes_1 = regexes; _i < regexes_1.length; _i++) {
        var regex = regexes_1[_i];
        if (!new RegExp(regex).test(functionName.toLowerCase())) {
            return false;
        }
    }
    return true;
};
exports.IsOperationAllowed = IsOperationAllowed;
function testContract2Arg(name, isBoolean, op) {
    var isEuint64Allowed = (0, exports.IsOperationAllowed)(name, common_1.EInputType.indexOf("euint64"));
    var isEuint128Allowed = (0, exports.IsOperationAllowed)(name, common_1.EInputType.indexOf("euint128"));
    var isEuint256Allowed = (0, exports.IsOperationAllowed)(name, common_1.EInputType.indexOf("euint256"));
    var func = "function ".concat(name, "(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {\n        if (Utils.cmp(test, \"").concat(name, "(euint8,euint8)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint8(a), FHE.asEuint8(b)));\n        } else if (Utils.cmp(test, \"").concat(name, "(euint16,euint16)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint16(a), FHE.asEuint16(b)));\n        } else if (Utils.cmp(test, \"").concat(name, "(euint32,euint32)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint32(a), FHE.asEuint32(b)));\n        }");
    if (isEuint64Allowed) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(euint64,euint64)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint64(a), FHE.asEuint64(b)));\n        }");
    }
    if (isEuint128Allowed) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(euint128,euint128)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint128(a), FHE.asEuint128(b)));\n        }");
    }
    if (isEuint256Allowed) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(euint256,euint256)\")) {\n            return FHE.decrypt(FHE.").concat(name, "(FHE.asEuint256(a), FHE.asEuint256(b))); \n        }");
    }
    func += " else if (Utils.cmp(test, \"euint8.".concat(name, "(euint8)\")) {\n            return FHE.decrypt(FHE.asEuint8(a).").concat(name, "(FHE.asEuint8(b)));\n        } else if (Utils.cmp(test, \"euint16.").concat(name, "(euint16)\")) {\n            return FHE.decrypt(FHE.asEuint16(a).").concat(name, "(FHE.asEuint16(b)));\n        } else if (Utils.cmp(test, \"euint32.").concat(name, "(euint32)\")) {\n            return FHE.decrypt(FHE.asEuint32(a).").concat(name, "(FHE.asEuint32(b)));\n        }");
    if (isEuint64Allowed) {
        func += " else if (Utils.cmp(test, \"euint64.".concat(name, "(euint64)\")) {\n            return FHE.decrypt(FHE.asEuint64(a).").concat(name, "(FHE.asEuint64(b)));\n        }");
    }
    if (isEuint128Allowed) {
        func += " else if (Utils.cmp(test, \"euint128.".concat(name, "(euint128)\")) {\n            return FHE.decrypt(FHE.asEuint128(a).").concat(name, "(FHE.asEuint128(b)));\n        }");
    }
    if (isEuint256Allowed) {
        func += " else if (Utils.cmp(test, \"euint256.".concat(name, "(euint256)\")) {\n            return FHE.decrypt(FHE.asEuint256(a).").concat(name, "(FHE.asEuint256(b)));\n        }");
    }
    if (op) {
        func += " else if (Utils.cmp(test, \"euint8 ".concat(op, " euint8\")) {\n            return FHE.decrypt(FHE.asEuint8(a) ").concat(op, " FHE.asEuint8(b));\n        } else if (Utils.cmp(test, \"euint16 ").concat(op, " euint16\")) {\n            return FHE.decrypt(FHE.asEuint16(a) ").concat(op, " FHE.asEuint16(b));\n        } else if (Utils.cmp(test, \"euint32 ").concat(op, " euint32\")) {\n            return FHE.decrypt(FHE.asEuint32(a) ").concat(op, " FHE.asEuint32(b));\n        }");
        if (isEuint64Allowed) {
            func += " else if (Utils.cmp(test, \"euint64 ".concat(op, " euint64\")) {\n            return FHE.decrypt(FHE.asEuint64(a) ").concat(op, " FHE.asEuint64(b));\n        }");
        }
        if (isEuint128Allowed) {
            func += " else if (Utils.cmp(test, \"euint128 ".concat(op, " euint128\")) {\n            return FHE.decrypt(FHE.asEuint128(a) ").concat(op, " FHE.asEuint128(b));\n        }");
        }
        if (isEuint256Allowed) {
            func += " else if (Utils.cmp(test, \"euint256 ".concat(op, " euint256\")) {\n            return FHE.decrypt(FHE.asEuint256(a) ").concat(op, " FHE.asEuint256(b));\n        }");
        }
    }
    if (isBoolean) {
        func += " else if (Utils.cmp(test, \"".concat(name, "(ebool,ebool)\")) {\n            bool aBool = true;\n            bool bBool = true;\n            if (a == 0) {\n                aBool = false;\n            }\n            if (b == 0) {\n                bBool = false;\n            }\n            if (FHE.decrypt(FHE.").concat(name, "(FHE.asEbool(aBool), FHE.asEbool(bBool)))) {\n                return 1;\n            }\n            return 0;\n        } else if (Utils.cmp(test, \"ebool.").concat(name, "(ebool)\")) {\n            bool aBool = true;\n            bool bBool = true;\n            if (a == 0) {\n                aBool = false;\n            }\n            if (b == 0) {\n                bBool = false;\n            }\n            if (FHE.asEbool(aBool).").concat(name, "(FHE.asEbool(bBool)).decrypt()) {\n                return 1;\n            }\n            return 0;\n        }");
        if (op) {
            func += " else if (Utils.cmp(test, \"ebool ".concat(op, " ebool\")) {\n            bool aBool = true;\n            bool bBool = true;\n            if (a == 0) {\n                aBool = false;\n            }\n            if (b == 0) {\n                bBool = false;\n            }\n            if (FHE.decrypt(FHE.asEbool(aBool) ").concat(op, " FHE.asEbool(bBool))) {\n                return 1;\n            }\n            return 0;\n        }");
        }
    }
    func += "\n    \n        revert TestNotFound(test);\n    }";
    var abi = "export interface ".concat((0, exports.capitalize)(name), "TestType extends BaseContract {\n    ").concat(name, ": (test: string, a: bigint, b: bigint) => Promise<bigint>;\n}\n");
    return [generateTestContract(name, func), abi];
}
exports.testContract2Arg = testContract2Arg;
function genAbiFile(abi) {
    return "import { BaseContract } from 'ethers';\n".concat(abi, "\n\n");
}
exports.genAbiFile = genAbiFile;
function SolTemplate1Arg(name, input1, returnType) {
    var docString = "\n    /// @notice Performs the ".concat(name, " operation on a ciphertext\n    /// @dev Verifies that the input value matches a valid ciphertext. Pure in this function is marked as a hack/workaround - note that this function is NOT pure as fetches of ciphertexts require state access\n    /// @param input1 the input ciphertext\n    ");
    if (name === "not" && input1 === "ebool") {
        return "\n\n    /// @notice Performs the \"not\" for the ebool type\n    /// @dev Implemented by a workaround due to ebool being a euint8 type behind the scenes, therefore xor is needed to assure that not(true) = false and vise-versa\n    /// @param value input ebool ciphertext\n    /// @return Result of the not operation on `value` \n    function not(ebool value) internal pure returns (ebool) {\n        return xor(value, asEbool(true));\n    }";
    }
    var returnStr = returnType === "none" ? "" : "returns (".concat(returnType, ")");
    var funcBody = docString;
    funcBody += "function ".concat(name, "(").concat(input1, " input1) internal pure ").concat(returnStr, " {");
    if ((0, common_1.valueIsEncrypted)(input1)) {
        // Get the proper function
        funcBody += "\n        if (!isInitialized(input1)) {\n            input1 = ".concat(asEuintFuncName(input1), "(0);\n        }");
        var unwrap = "".concat(common_1.UnderlyingTypes[input1], " unwrappedInput1 = ").concat(unwrapType(input1, "input1"), ";");
        var getResult = function (inputName) {
            return "FheOps(Precompiles.Fheos).".concat(name, "(").concat(common_1.UintTypes[input1], ", ").concat(inputName, ");");
        };
        if ((0, common_1.valueIsEncrypted)(returnType)) {
            // input and return type are encrypted - not/neg other unary functions
            funcBody += "\n        ".concat(unwrap, "\n        bytes memory inputAsBytes = Common.toBytes(unwrappedInput1);\n        bytes memory b = ").concat(getResult("inputAsBytes"), "\n        uint256 result = Impl.getValue(b);\n        return ").concat(wrapType(returnType, "result"), ";\n    }");
        }
        else if (returnType === "none") {
            // this is essentially req
            funcBody += "\n        ".concat(unwrap, "\n        bytes memory inputAsBytes = Common.toBytes(unwrappedInput1);\n        ").concat(getResult("inputAsBytes"), "\n    }");
        }
        else if ((0, common_1.valueIsPlaintext)(returnType)) {
            var returnTypeCamelCase = returnType.charAt(0).toUpperCase() + returnType.slice(1);
            var outputConvertor = "Common.bigIntTo".concat(returnTypeCamelCase, "(result);");
            funcBody += "\n        ".concat(unwrap, "\n        bytes memory inputAsBytes = Common.toBytes(unwrappedInput1);\n        uint256 result = ").concat(getResult("inputAsBytes"), "\n        return ").concat(outputConvertor, "\n    }");
        }
    }
    else {
        throw new Error("unsupported function of 1 input that is not encrypted");
    }
    return funcBody;
}
exports.SolTemplate1Arg = SolTemplate1Arg;
function SolTemplate3Arg(name, input1, input2, input3, returnType) {
    if ((0, common_1.valueIsEncrypted)(returnType)) {
        if ((0, common_1.valueIsEncrypted)(input1) &&
            (0, common_1.valueIsEncrypted)(input2) &&
            (0, common_1.valueIsEncrypted)(input3) &&
            input1 === "ebool") {
            if (input2 !== input3) {
                return "";
            }
            return "\n\n    function ".concat(name, "(").concat(input1, " input1, ").concat(input2, " input2, ").concat(input3, " input3) internal pure returns (").concat(returnType, ") {\n        if (!isInitialized(input1)) {\n            input1 = ").concat(asEuintFuncName(input1), "(0);\n        }\n        if (!isInitialized(input2)) {\n            input2 = ").concat(asEuintFuncName(input2), "(0);\n        }\n        if (!isInitialized(input3)) {\n            input3 = ").concat(asEuintFuncName(input3), "(0);\n        }\n\n        ").concat(common_1.UnderlyingTypes[input1], " unwrappedInput1 = ").concat(unwrapType(input1, "input1"), ";\n        ").concat(common_1.UnderlyingTypes[input2], " unwrappedInput2 = ").concat(unwrapType(input2, "input2"), ";\n        ").concat(common_1.UnderlyingTypes[input3], " unwrappedInput3 = ").concat(unwrapType(input3, "input3"), ";\n\n        ").concat(common_1.UnderlyingTypes[returnType], " result = Impl.").concat(name, "(").concat(common_1.UintTypes[input2], ", unwrappedInput1, unwrappedInput2, unwrappedInput3);\n        return ").concat(wrapType(returnType, "result"), ";\n    }");
        }
        else {
            return "";
        }
    }
    else {
        throw new Error("Unsupported return type ".concat(returnType, " for 3 args"));
    }
}
exports.SolTemplate3Arg = SolTemplate3Arg;
function operatorFunctionName(funcName, forType) {
    return "operator".concat((0, exports.capitalize)(funcName)).concat((0, exports.capitalize)(forType));
}
var OperatorOverloadDecl = function (funcName, op, forType, unary, returnsBool) {
    var opOverloadName = operatorFunctionName(funcName, forType);
    var unaryParameters = unary ? "lhs" : "lhs, rhs";
    var funcParams = unaryParameters
        .split(", ")
        .map(function (key) {
        return "".concat(forType, " ").concat(key);
    })
        .join(", ");
    var returnType = returnsBool ? "ebool" : forType;
    return "\nusing {".concat(opOverloadName, " as ").concat(op, "} for ").concat(forType, " global;\n/// @notice Performs the ").concat(funcName, " operation\nfunction ").concat(opOverloadName, "(").concat(funcParams, ") pure returns (").concat(returnType, ") {\n    return FHE.").concat(funcName, "(").concat(unaryParameters, ");\n}\n");
};
exports.OperatorOverloadDecl = OperatorOverloadDecl;
var BindingLibraryType = function (type) {
    var typeCap = (0, exports.capitalize)(type);
    return "\n\nusing Bindings".concat(typeCap, " for ").concat(type, " global;\nlibrary Bindings").concat(typeCap, " {");
};
exports.BindingLibraryType = BindingLibraryType;
var OperatorBinding = function (funcName, forType, unary, returnsBool) {
    var unaryParameters = unary ? "lhs" : "lhs, rhs";
    var funcParams = unaryParameters
        .split(", ")
        .map(function (key) {
        return "".concat(forType, " ").concat(key);
    })
        .join(", ");
    var returnType = returnsBool ? "ebool" : forType;
    var docString = "\n    /// @notice Performs the ".concat(funcName, " operation\n    /// @dev Pure in this function is marked as a hack/workaround - note that this function is NOT pure as fetches of ciphertexts require state access\n    /// @param lhs input of type ").concat(forType, "\n    ");
    if (unary) {
        docString += "/// @param rhs second input of type ".concat(forType, "\n");
    }
    docString += "/// @return the result of the ".concat(funcName);
    return "\n    ".concat(docString, "\n    function ").concat(funcName, "(").concat(funcParams, ") internal pure returns (").concat(returnType, ") {\n        return FHE.").concat(funcName, "(").concat(unaryParameters, ");\n    }");
};
exports.OperatorBinding = OperatorBinding;
var shortenType = function (type) {
    if (type === "eaddress") {
        return "Eaddress";
    }
    return type === "ebool" ? "Bool" : "U" + type.slice(5); // get only number at the end
};
var CastBinding = function (thisType, targetType) {
    return "\n    function to".concat(shortenType(targetType), "(").concat(thisType, " value) internal pure returns (").concat(targetType, ") {\n        return FHE.as").concat((0, exports.capitalize)(targetType), "(value);\n    }");
};
exports.CastBinding = CastBinding;
var SealFromType = function (thisType) {
    return "\n    function ".concat(common_1.LOCAL_SEAL_FUNCTION_NAME, "(").concat(thisType, " value, bytes32 publicKey) internal pure returns (").concat(common_1.SEAL_RETURN_TYPE, " memory) {\n        return FHE.sealoutput(value, publicKey);\n    }");
};
exports.SealFromType = SealFromType;
var DecryptBinding = function (thisType) {
    return "\n    function ".concat(common_1.LOCAL_DECRYPT_FUNCTION_NAME, "(").concat(thisType, " value) internal pure returns (").concat((0, common_1.toPlaintextType)(thisType), ") {\n        return FHE.decrypt(value);\n    }");
};
exports.DecryptBinding = DecryptBinding;
