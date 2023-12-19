import {
    AllTypes,
    EInputType,
    EPlaintextType,
    EUintType,
    UnderlyingTypes,
    valueIsEncrypted,
    valueIsPlaintext
} from "./common";

export const preamble = () => {
    return `// SPDX-License-Identifier: BSD-3-Clause-Clear

pragma solidity >=0.8.13 <0.9.0;

import "./FheOS.sol";

type ebool is uint256;
type euint8 is uint256;
type euint16 is uint256;
type euint32 is uint256;

library Common {
    // Values used to communicate types to the runtime.
    uint8 internal constant ebool_tfhe_go = 0;
    uint8 internal constant euint8_tfhe_go = 0;
    uint8 internal constant euint16_tfhe_go = 1;
    uint8 internal constant euint32_tfhe_go = 2;
    
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
}

library Impl {
    function reencrypt(uint256 ciphertext, bytes32 publicKey) internal pure returns (bytes memory reencrypted) {
        bytes32[2] memory input;
        input[0] = bytes32(ciphertext);
        input[1] = publicKey;

        // Call the reencrypt precompile.
        reencrypted = FheOps(Precompiles.Fheos).reencrypt(bytes.concat(input[0], input[1]));

        return reencrypted;
    }

    function verify(bytes memory _ciphertextBytes, uint8 _toType) internal pure returns (uint256 result) {
        bytes memory input = bytes.concat(_ciphertextBytes, bytes1(_toType));

        bytes memory output;

        // Call the verify precompile.
        output = FheOps(Precompiles.Fheos).verify(input);
        result = getValue(output);
    }

    function cast(uint256 ciphertext, uint8 toType) internal pure returns (uint256 result) {
        bytes memory input = bytes.concat(bytes32(ciphertext), bytes1(toType));

        bytes memory output;

        // Call the cast precompile.
        output = FheOps(Precompiles.Fheos).cast(input);
        result = getValue(output);
    }
    
    function getValue(bytes memory a) internal pure returns (uint256 value) {
        assembly {
            value := mload(add(a, 0x20))
        }
    }
    
    function trivialEncrypt(uint256 value, uint8 toType) internal pure returns (uint256 result) {
        bytes memory input = bytes.concat(bytes32(value), bytes1(toType));

        bytes memory output;

        // Call the trivialEncrypt precompile.
        output = FheOps(Precompiles.Fheos).trivialEncrypt(input);

        result = getValue(output);
    }
    
    function select(uint256 control, uint256 ifTrue, uint256 ifFalse) internal pure returns (uint256 result) {
        bytes memory input = bytes.concat(bytes32(control), bytes32(ifTrue), bytes32(ifFalse));

        bytes memory output;

        // Call the trivialEncrypt precompile.
        output = FheOps(Precompiles.Fheos).select(input);

        result = getValue(output);
    }

}

library TFHE {
    euint8 constant NIL8 = euint8.wrap(0);
    euint16 constant NIL16 = euint16.wrap(0);
    euint32 constant NIL32 = euint32.wrap(0);

    // Return true if the enrypted integer is initialized and false otherwise.
    function isInitialized(ebool v) internal pure returns (bool) {
        return ebool.unwrap(v) != 0;
    }

    // Return true if the enrypted integer is initialized and false otherwise.
    function isInitialized(euint8 v) internal pure returns (bool) {
        return euint8.unwrap(v) != 0;
    }

    // Return true if the enrypted integer is initialized and false otherwise.
    function isInitialized(euint16 v) internal pure returns (bool) {
        return euint16.unwrap(v) != 0;
    }

    // Return true if the enrypted integer is initialized and false otherwise.
    function isInitialized(euint32 v) internal pure returns (bool) {
        return euint32.unwrap(v) != 0;
    }
    
    function getValue(bytes memory a) internal pure returns (uint256 value) {
        assembly {
            value := mload(add(a, 0x20))
        }
    }

    function mathHelper(
        uint256 lhs,
        uint256 rhs,
        function(bytes memory) external pure returns (bytes memory) impl
    ) internal pure returns (uint256 result) {
        bytes memory input = bytes.concat(bytes32(lhs), bytes32(rhs));

        bytes memory output;
        // Call the add precompile.

        output = impl(input);
        result = getValue(output);
    }`;
}

export const PostFix = () => {
    return `\n}`;
}

const castFromEncrypted = (fromType: string, toType: string, name: string): string => {
    return `Impl.cast(${fromType}.unwrap(${name}), Common.${toType}_tfhe_go)`;
}

const castFromPlaintext = (name: string, toType: string): string => {
    return `Impl.trivialEncrypt(${name}, Common.${toType}_tfhe_go)`;
}

const castFromBytes = (name: string, toType: string): string => {
    return `Impl.verify(${name}, Common.${toType}_tfhe_go)`;
}

const castToEbool = (name: string, fromType: string): string => {
    return `\n\tfunction asEbool(${fromType} value) internal pure returns (ebool) {
        return ne(${name},  as${capitalize(fromType)}(0));
    }`;
}

export const AsTypeFunction = (fromType: string, toType: string) => {
    let castString = castFromEncrypted(fromType, toType, "value");
    if (fromType === 'bool' && toType === 'ebool') {
        return `\n\tfunction asEbool(bool value) internal pure returns (ebool) {
        uint256 sVal = 0;
        if (value) {
            sVal = 1;
        }

        return asEbool(sVal);
    }`;
    }
    else if (fromType === 'bytes memory') {
        castString = castFromBytes("value", toType)
    } else if (EPlaintextType.includes(fromType)) {
        castString = castFromPlaintext("value", toType);
    } else if (toType === "ebool") {
        return castToEbool("value", fromType);
    } else if (!EInputType.includes(fromType)) {
        throw new Error(`Unsupported type for casting: ${fromType}`)
    }

    return `\n\tfunction as${capitalize(toType)}(${fromType} value) internal pure returns (${toType}) {
        return ${toType}.wrap(${castString});
    }`;
}

function TypeCastTestingFunction(fromType: string, fromTypeForTs: string, toType: string, fromTypeEncrypted?: string) {
    let to = capitalize(toType);
    const retType = to.slice(1);
    let testType = fromTypeEncrypted? fromTypeEncrypted : fromType;
    testType = testType === 'bytes memory' ? 'PreEncrypted' : capitalize(testType);
    testType = testType === 'Uint256' ? 'Plaintext' : testType;
    const encryptedVal = fromTypeEncrypted ? `TFHE.as${capitalize(fromTypeEncrypted)}(val)` : "val";
    const func = `\n\tfunction castFrom${testType}To${to}(${fromType} val) public pure returns (${retType}) {
        return TFHE.decrypt(TFHE.as${to}(${encryptedVal}));
    }`;

    let retTypeTs = (retType === 'bool') ? 'boolean' : retType;
    retTypeTs = retTypeTs.includes("uint") ? 'bigint' : retTypeTs;
    const abi =  `\tcastFrom${testType}To${to}: (val: ${fromTypeForTs}) => Promise<${retTypeTs}>;\n`;
    return [func, abi];
}

export function AsTypeTestingContract(type: string) {
    let funcs = "";
    let abi = `export interface As${capitalize(type)}TestType extends Contract {\n`;
    for (const fromType of EInputType.concat('uint256', 'bytes memory')) {
        if (type === fromType) {
            continue;
        }

        const fromTypeTs = (fromType === 'bytes memory') ? 'Uint8Array' : `bigint`;
        const fromTypeSol = (fromType === 'bytes memory') ? fromType : `uint256`;
        const fromTypeEncrypted = (EInputType.includes(fromType)) ? fromType : undefined;
        const contractInfo = TypeCastTestingFunction(fromTypeSol, fromTypeTs, type, fromTypeEncrypted);
        funcs += contractInfo[0];
        abi += contractInfo[1];

    }

    funcs = funcs.slice(1);

    abi += `}\n`;

    return [generateTestContract(`As${capitalize(type)}`, funcs), abi];
}


const unwrapType = (typeName: EUintType, inputName: string): string => `${typeName}.unwrap(${inputName})`;
const wrapType = (resultType: EUintType, inputName: string): string => `${resultType}.wrap(${inputName})`;
const asEuintFuncName = (typeName: EUintType): string => `as${capitalize(typeName)}`;
export const capitalize = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);

export function SolTemplate2Arg(name: string, input1: AllTypes, input2: AllTypes, returnType: AllTypes) {

    // special names for reencrypt function (don't check name === reencrypt) because the name could change
    let variableName1 = input2 === "bytes32" ? "value" : "lhs";
    let variableName2 = input2 === "bytes32" ? "publicKey" : "rhs";

    let funcBody = `\n
    function ${name}(${input1} ${variableName1}, ${input2} ${variableName2}) internal pure returns (${returnType}) {`;


    if (valueIsEncrypted(input1)) {
        // both inputs encrypted - this is a generic math function. i.e. div, mul, eq, etc.
        // 1. possibly cast input1
        // 2. possibly cast input2
        // 3. possibly cast return type
        if (valueIsEncrypted(input2) && input1 !== input2) {
           return "";
        }
        if (valueIsEncrypted(input2) || EPlaintextType.includes(input2)) {
            let input2Cast = input1 === input2 ? variableName2 : `${asEuintFuncName(input1)}(${variableName2})`;
            //
            funcBody += `
        if (!isInitialized(${variableName1}) || !isInitialized(${variableName2})) {
            revert("One or more inputs are not initialized.");
        }
        ${UnderlyingTypes[input1]} unwrappedInput1 = ${unwrapType(input1, variableName1)};
        ${UnderlyingTypes[input1]} unwrappedInput2 = ${unwrapType(input1, `${input2Cast}`)};
`;
            if (valueIsEncrypted(returnType)) {
                funcBody += `
        ${UnderlyingTypes[returnType]} result = mathHelper(unwrappedInput1, unwrappedInput2, FheOps(Precompiles.Fheos).${name});
        return ${wrapType(returnType, "result")};`;
            }
            else {
                funcBody += `
        ${returnType} result = mathHelper(unwrappedInput1, unwrappedInput2, FheOps(Precompiles.Fheos).${name});
        return result;`;
            }
        }

        else if (input2 === "bytes32") {
            // **** Value 1 is encrypted, value 2 is bytes32 - this is basically reencrypt/wrapForUser
            funcBody += `
        ${UnderlyingTypes[input1]} unwrapped = ${unwrapType(input1, variableName1)};

        return Impl.${name}(unwrapped, ${variableName2});`;
        }
    } else {
        // don't support input 1 is plaintext
        throw new Error("Unsupported plaintext input1");
    }

    funcBody += `\n\t}`;

    return funcBody;
}

export function testContract2ArgBoolRes(name: string, isBoolean: boolean) {
    let func = `function ${name}(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "${name}(euint8,euint8)")) {
            if (TFHE.decrypt(TFHE.${name}(TFHE.asEuint8(a), TFHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "${name}(euint16,euint16)")) {
            if (TFHE.decrypt(TFHE.${name}(TFHE.asEuint16(a), TFHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "${name}(euint32,euint32)")) {
            if (TFHE.decrypt(TFHE.${name}(TFHE.asEuint32(a), TFHE.asEuint32(b)))) {
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
            if (TFHE.decrypt(TFHE.${name}(TFHE.asEbool(aBool), TFHE.asEbool(bBool)))) {
                return 1;
            }

            return 0;
        }`;
    }
    func += ` else {
            require(false, string(abi.encodePacked("test '", test, "' not found")));
        }
    }`;

    const abi = `export interface ${capitalize(name)}TestType extends Contract {
    ${name}: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}\n`;
    return [generateTestContract(name, func), abi];
}

export function testContract1Arg(name: string) {
    let func = `function ${name}(string calldata test, uint256 a) public pure returns (uint256 output) {
        if (Utils.cmp(test, "${name}(euint8)")) {
            return TFHE.decrypt(TFHE.${name}(TFHE.asEuint8(a)));
        } else if (Utils.cmp(test, "${name}(euint16)")) {
            return TFHE.decrypt(TFHE.${name}(TFHE.asEuint16(a)));
        } else if (Utils.cmp(test, "${name}(euint32)")) {
            return TFHE.decrypt(TFHE.${name}(TFHE.asEuint32(a)));
        } else if (Utils.cmp(test, "${name}(ebool)")) {
            bool aBool = true;
            if (a == 0) {
                aBool = false;
            }
            
            if (TFHE.decrypt(TFHE.${name}(TFHE.asEbool(aBool)))) {
                return 1;
            }
            
            return 0;
        } else {
            require(false, string(abi.encodePacked("test '", test, "' not found")));
        }
    }`;
    const abi = `export interface ${capitalize(name)}TestType extends Contract {
    ${name}: (test: string, a: bigint) => Promise<bigint>;
}\n`
    return [generateTestContract(name, func), abi];
}

export function generateTestContract(name: string, testFunc: string) {
    return `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "../../FHE.sol";
import "./utils/Utils.sol";

contract ${capitalize(name)}Test {
    using Utils for *;

    ${testFunc}\n
}`;
}
export function testContractReq() {
    // Req is failing on EthCall so we need to make it as tx for now
    let func = `function req(string calldata test, uint256 a) public {
        if (Utils.cmp(test, "req(euint8)")) {
            TFHE.req(TFHE.asEuint8(a));
        } else if (Utils.cmp(test, "req(euint16)")) {
            TFHE.req(TFHE.asEuint16(a));
        } else if (Utils.cmp(test, "req(euint32)")) {
            TFHE.req(TFHE.asEuint32(a));
        } else if (Utils.cmp(test, "req(ebool)")) {
            bool b = true;
            if (a == 0) {
                b = false;
            }
            TFHE.req(TFHE.asEbool(b));
        } else {
            require(false, string(abi.encodePacked("test '", test, "' not found")));
        }
    }`;
    const abi = `export interface ReqTestType extends Contract {
    req: (test: string, a: bigint) => Promise<()>;
}\n`
    return [generateTestContract("req", func), abi];
}

export function testContractReencrypt() {
    let func = `function reencrypt(string calldata test, uint256 a, bytes32 pubkey) public pure returns (bytes memory reencrypted) {
        if (Utils.cmp(test, "reencrypt(euint8)")) {
            return TFHE.reencrypt(TFHE.asEuint8(a), pubkey);
        } else if (Utils.cmp(test, "reencrypt(euint16)")) {
            return TFHE.reencrypt(TFHE.asEuint16(a), pubkey);
        } else if (Utils.cmp(test, "reencrypt(euint32)")) {
            return TFHE.reencrypt(TFHE.asEuint32(a), pubkey);
        } else if (Utils.cmp(test, "reencrypt(ebool)")) {
            bool b = true;
            if (a == 0) {
                b = false;
            }
            
            return TFHE.reencrypt(TFHE.asEbool(b), pubkey);
        } else {
            require(false, string(abi.encodePacked("test '", test, "' not found")));
        }
    }`;
    const abi = `export interface ReencryptTestType extends Contract {
    reencrypt: (test: string, a: bigint, pubkey: Uint8Array) => Promise<Uint8Array>;
}\n`
    return [generateTestContract("reencrypt", func), abi];
}

export function testContract3Arg(name: string) {
    let func = `function ${name}(string calldata test, bool c, uint256 a, uint256 b) public pure returns (uint256 output) {
        ebool condition = TFHE.asEbool(c);
        if (Utils.cmp(test, "${name}: euint8")) {
            return TFHE.decrypt(TFHE.${name}(condition, TFHE.asEuint8(a), TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "${name}: euint16")) {
            return TFHE.decrypt(TFHE.${name}(condition, TFHE.asEuint16(a), TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "${name}: euint32")) {
            return TFHE.decrypt(TFHE.${name}(condition, TFHE.asEuint32(a), TFHE.asEuint32(b)));
        } else if (Utils.cmp(test, "${name}: ebool")) {
            bool aBool = true;
            bool bBool = true;
            
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            
            if(TFHE.decrypt(TFHE.${name}(condition, TFHE.asEbool(aBool), TFHE.asEbool(bBool)))) {
                return 1;
            }
            
            return 0;
        } else {
            require(false, string(abi.encodePacked("test '", test, "' not found")));
        }
    }`;
    const abi = `export interface ${capitalize(name)}TestType extends Contract {
    ${name}: (test: string, c: boolean, a: bigint, b: bigint) => Promise<bigint>;
}\n`
    return [generateTestContract(name, func), abi];
}

export function testContract2Arg(name: string, isBoolean: boolean) {
    let func = ` function ${name}(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "${name}(euint8,euint8)")) {
            return TFHE.decrypt(TFHE.${name}(TFHE.asEuint8(a), TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "${name}(euint16,euint16)")) {
            return TFHE.decrypt(TFHE.${name}(TFHE.asEuint16(a), TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "${name}(euint32,euint32)")) {
            return TFHE.decrypt(TFHE.${name}(TFHE.asEuint32(a), TFHE.asEuint32(b)));
        }`
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
            if (TFHE.decrypt(TFHE.${name}(TFHE.asEbool(aBool), TFHE.asEbool(bBool)))) {
                return 1;
            }

            return 0;
        }`
    }
    func += ` else {
            require(false, string(abi.encodePacked("test '", test, "' not found")));
        }
    }`;
    const abi = `export interface ${capitalize(name)}TestType extends Contract {
    ${name}: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}\n`
    return [generateTestContract(name,func), abi];
}

export function genAbiFile(abi: string) {
    return `import { Contract } from 'ethers';\n${abi}\n\n`;
}


export function SolTemplate1Arg(name: string, input1: AllTypes, returnType: AllTypes) {

    if (name === "not" && input1 === "ebool") {
        return `\n\n\t// "not" for ebool not working in the traditional way as it is converting ebool to euint8
    // Ebool(true) is Euint8(1) so !Ebool(true) is !Euint8(1) which is Euint8(254) which is still Ebool(true)
    function not(ebool value) internal pure returns (ebool) {
        return xor(value,  asEbool(true));
    }`;
    }

    let returnStr = returnType === "none" ? `` : `returns (${returnType})`;

    let funcBody = `\n
    function ${name}(${input1} input1) internal pure ${returnStr} {`;

    if (valueIsEncrypted(input1)) {
        funcBody += `
        if (!isInitialized(input1)) {
            revert("One or more inputs are not initialized.");
        }`;
        let unwrap = `${UnderlyingTypes[input1]} unwrappedInput1 = ${unwrapType(input1, "input1")};`;
        let getResult = (inputName: string) => `FheOps(Precompiles.Fheos).${name}(${inputName});`;

        if (valueIsEncrypted(returnType)) {
            // input and return type are encrypted - not/neg other unary functions
            funcBody += `
        ${unwrap}
        bytes memory inputAsBytes = bytes.concat(bytes32(unwrappedInput1));
        bytes memory b = ${getResult("inputAsBytes")}
        uint256 result = Impl.getValue(b);
        return ${wrapType(returnType, "result")};
    }`
        } else if (returnType === "none") {
            // this is essentially req
            funcBody += `
        ${unwrap}
        bytes memory inputAsBytes = bytes.concat(bytes32(unwrappedInput1));
        ${getResult("inputAsBytes")}
    }`;
        } else if (valueIsPlaintext(returnType)){
            let returnTypeCamelCase = returnType.charAt(0).toUpperCase() + returnType.slice(1);
            let outputConvertor = `Common.bigIntTo${returnTypeCamelCase}(result);`
            funcBody += `
        ${unwrap}
        bytes memory inputAsBytes = bytes.concat(bytes32(unwrappedInput1));
        uint256 result = ${getResult("inputAsBytes")}
        return ${outputConvertor}
    }`
        }
    } else {
        throw new Error("unsupported function of 1 input that is not encrypted");
    }
    return funcBody;
}

export function SolTemplate3Arg(name: string, input1: AllTypes, input2: AllTypes,input3: AllTypes,returnType: AllTypes) {
    if (valueIsEncrypted(returnType)) {
        if (valueIsEncrypted(input1) && valueIsEncrypted(input2) && valueIsEncrypted(input3) && input1 === 'ebool') {
            if (input2 !== input3) {
                return "";
            }

            return `\n
    function ${name}(${input1} input1, ${input2} input2, ${input3} input3) internal pure returns (${returnType}) {
        if (!isInitialized(input1) || !isInitialized(input2) || !isInitialized(input3)) {
            revert("One or more inputs are not initialized.");
        }

        ${UnderlyingTypes[input1]} unwrappedInput1 = ${unwrapType(input1, "input1")};
        ${UnderlyingTypes[input2]} unwrappedInput2 = ${unwrapType(input2, "input2")};
        ${UnderlyingTypes[input3]} unwrappedInput3 = ${unwrapType(input3, `input3`)};

        ${UnderlyingTypes[returnType]} result = Impl.${name}(unwrappedInput1, unwrappedInput2, unwrappedInput3);
        return ${wrapType(returnType, "result")};
    }`;
        } else {
            return "";
        }
    } else {
        throw new Error(`Unsupported return type ${returnType} for 3 args`)
    }
}

function operatorFunctionName(funcName: string, forType: "ebool" | "euint8"  | "euint16" | "euint32") {
    return `operator${capitalize(funcName)}${capitalize(forType)}`;
}

export const OperatorOverloadDecl = (funcName: string, op: string, forType: EUintType, unary: boolean) => {
    let opOverloadName = operatorFunctionName(funcName, forType);
    let unaryParameters = unary ? 'lhs' : 'lhs, rhs';
    let funcParams = unaryParameters.split(',').map((key) => {return `${forType} ${key}`}).join(', ')

    return `\nusing {${opOverloadName} as ${op}, Bindings${capitalize(forType)}.${funcName}} for ${forType} global;
function ${opOverloadName}(${funcParams}) pure returns (${forType}) {
    return TFHE.${funcName}(${unaryParameters});
}\n`;
}

export const BindingsWithoutOperator = (funcName: string, forType: string) => {
    return `\nusing {Bindings${capitalize(forType)}.${funcName}} for ${forType} global;\n`;
}

export const BindingLibraryType = (type: string) => {
    return `\nlibrary Bindings${capitalize(type)} {`;
}

export const OperatorBinding = (funcName: string, forType: string, unary: boolean) => {
    let unaryParameters = unary ? 'lhs' : 'lhs, rhs';
    let funcParams = unaryParameters.split(',').map((key) => {return `${forType} ${key}`}).join(', ')

    if (funcName === "eq") {
        forType = "ebool"
    }

    return `
    function ${funcName}(${funcParams}) pure internal returns (${forType}) {
        return TFHE.${funcName}(${unaryParameters});
    }`;
}
