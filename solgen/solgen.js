"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
var contracts_parser_1 = require("./contracts_parser");
var fs = require("fs");
var templates_1 = require("./templates");
var common_1 = require("./common");
var generateMetadataPayload = function () { return __awaiter(void 0, void 0, void 0, function () {
    var result;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0: return [4 /*yield*/, (0, contracts_parser_1.getFunctionsFromGo)('../precompiles/contracts.go')];
            case 1:
                result = _a.sent();
                return [2 /*return*/, result.map(function (value) {
                        return {
                            functionName: value.name,
                            hasDifferentInputTypes: !value.needsSameType,
                            inputCount: value.paramsCount,
                            returnValueType: value.returnType,
                            inputs: value.inputTypes,
                            isBooleanMathOp: value.isBooleanMathOp
                        };
                    })];
        }
    });
}); };
// Function to generate all combinations of parameters.
function generateCombinations(arr, current, index) {
    if (current === void 0) { current = []; }
    if (index === void 0) { index = 0; }
    if (index === arr.length) {
        // Only add to the result if there are elements in the current combination
        return current.length === 0 ? [] : [current];
    }
    var result = [];
    // Add each element of the current sub-array to the combination
    for (var _i = 0, _a = arr[index]; _i < _a.length; _i++) {
        var item = _a[_i];
        result = result.concat(generateCombinations(arr, current.concat(item), index + 1));
    }
    return result;
}
var getReturnType = function (inputs, isBooleanMathOp, returnType) {
    if (returnType === 'plaintext') {
        if (inputs.length != 1) {
            throw new Error("expecting exactly one input for functions returning plaintext");
        }
        var inputType = inputs[0].split(' ')[1];
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
    // if (isBooleanMathOp) {
    //     return "ebool";
    // }
    var maxRank = 0;
    for (var _i = 0, inputs_1 = inputs; _i < inputs_1.length; _i++) {
        var input = inputs_1[_i];
        var inputType = input.split(' ')[1];
        maxRank = Math.max(common_1.EInputType.indexOf(inputType), common_1.EPlaintextType.indexOf(inputType), maxRank);
    }
    return common_1.EInputType[maxRank];
};
function getAllFunctionDeclarations(functionName, functions, isBooleanMathOp, returnValueType) {
    var functionDecl = "function ".concat(functionName);
    // Generate all combinations of input parameters.
    var allCombinations = generateCombinations(functions);
    // Create function declarations for each combination.
    return allCombinations.map(function (combination) {
        var returnType = getReturnType(combination, isBooleanMathOp, returnValueType);
        var returnStr = "internal pure returns (".concat(returnType, ");");
        return "".concat(functionDecl, "(").concat(combination.join(', '), ") ").concat(returnStr);
    });
}
/** Generates a Solidity test contract based on the provided metadata */
var generateSolidityTestContract = function (metadata) {
    var functionName = metadata.functionName, inputCount = metadata.inputCount, hasDifferentInputTypes = metadata.hasDifferentInputTypes, returnValueType = metadata.returnValueType, inputs = metadata.inputs, isBooleanMathOp = metadata.isBooleanMathOp;
    if (functionName === "req") {
        return (0, templates_1.testContractReq)();
    }
    if (functionName === "reencrypt") {
        return (0, templates_1.testContractReencrypt)();
    }
    if (inputCount === 2 && inputs[0] === "encrypted" && inputs[1] === "encrypted") {
        if (returnValueType === "ebool") {
            return (0, templates_1.testContract2ArgBoolRes)(functionName, isBooleanMathOp);
        }
        return (0, templates_1.testContract2Arg)(functionName, isBooleanMathOp);
    }
    if (inputCount === 1 && inputs[0] === "encrypted" && returnValueType === "encrypted") {
        return (0, templates_1.testContract1Arg)(functionName);
    }
    if (inputCount === 3) {
        return (0, templates_1.testContract3Arg)(functionName);
    }
    console.log("Function ".concat(functionName, " with ").concat(inputCount, " inputs that are ").concat(inputs, " is not implemented"));
    return ["", ""];
};
/**
 * Generates a Solidity function based on the provided metadata
 * This generates all the different types of function headers that can exist
 */
var genSolidityFunctionHeaders = function (metadata) {
    var functionName = metadata.functionName, inputCount = metadata.inputCount, hasDifferentInputTypes = metadata.hasDifferentInputTypes, returnValueType = metadata.returnValueType, inputs = metadata.inputs, isBooleanMathOp = metadata.isBooleanMathOp;
    var functions = [];
    inputs.forEach(function (input, idx) {
        var inputVariants = [];
        switch (input) {
            case "encrypted":
                for (var _i = 0, EInputType_1 = common_1.EInputType; _i < EInputType_1.length; _i++) {
                    var inputType = EInputType_1[_i];
                    if (inputs.length === 2 && !isBooleanMathOp && common_1.EComparisonType.includes(inputType)) {
                        continue;
                    }
                    inputVariants.push("input".concat(idx, " ").concat(inputType));
                }
                break;
            case "plaintext":
                for (var _a = 0, EPlaintextType_1 = common_1.EPlaintextType; _a < EPlaintextType_1.length; _a++) {
                    var inputType = EPlaintextType_1[_a];
                    inputVariants.push("input".concat(idx, " ").concat(inputType));
                }
                break;
            default:
                inputVariants.push("input".concat(idx, " ").concat(input));
        }
        functions.push(inputVariants);
    });
    return getAllFunctionDeclarations(functionName, functions, isBooleanMathOp, returnValueType);
};
// Regular expression to match the Solidity function signature pattern
var functionPattern = /function (\w+)\((.*?)\) internal pure returns \((.*?)\);/;
/**
 * Parses a Solidity function definition into its components.
 *
 * @param funcDef Solidity function definition as string.
 * @returns An object containing the functionName, inputTypes, and outputType.
 */
var parseFunctionDefinition = function (funcDef) {
    var match = funcDef.match(functionPattern);
    if (!match) {
        throw new Error("Invalid function definition format for ".concat(funcDef));
    }
    var functionName = match[1], inputs = match[2], outputType = match[3];
    var inputTypes = inputs.split(',').map(function (input) {
        return input.trim().split(/\s+/).pop();
    });
    var inputPlaintext = inputTypes[0].startsWith('e') ? "none" : "all";
    return {
        funcName: functionName,
        inputs: inputTypes,
        returnType: outputType,
        inputPlaintext: inputPlaintext
    };
};
// Helper function to capitalize type name for asEuintX function call.
// This will generate the Solidity function body based on the function definition provided.
var generateSolidityFunction = function (parsedFunction) {
    var funcName = parsedFunction.funcName, inputs = parsedFunction.inputs, returnType = parsedFunction.returnType;
    switch (inputs.length) {
        case 1:
            return (0, templates_1.SolTemplate1Arg)(funcName, inputs[0], returnType);
        case 2:
            return (0, templates_1.SolTemplate2Arg)(funcName, inputs[0], inputs[1], returnType);
        case 3:
            return (0, templates_1.SolTemplate3Arg)(funcName, inputs[0], inputs[1], inputs[2], returnType);
        default:
            throw new Error("Unknown number of inputs");
    }
};
var main = function () { return __awaiter(void 0, void 0, void 0, function () {
    var metadata, solidityHeaders, testContracts, testContractsAbis, importLineHelper, _i, metadata_1, func, testContract, outputFile, _a, solidityHeaders_1, fn, funcDefinition, _b, _c, fromType, _d, EInputType_2, toType, _e, EInputType_3, type, functionName, testContract, _f, _g, testContract;
    return __generator(this, function (_h) {
        switch (_h.label) {
            case 0: return [4 /*yield*/, generateMetadataPayload()];
            case 1:
                metadata = _h.sent();
                solidityHeaders = [];
                testContracts = {};
                testContractsAbis = "";
                importLineHelper = "import { ";
                for (_i = 0, metadata_1 = metadata; _i < metadata_1.length; _i++) {
                    func = metadata_1[_i];
                    // Decrypt is already tested in every test contract
                    if (func.functionName !== "decrypt") {
                        testContract = generateSolidityTestContract(func);
                        if (testContract[0] !== "") {
                            testContracts[(0, templates_1.capitalize)(func.functionName)] = testContract[0];
                            testContractsAbis += testContract[1];
                            importLineHelper += "".concat((0, templates_1.capitalize)(func.functionName), "TestType,\n");
                        }
                    }
                    // this generates solidity header functions for all the different possible types
                    solidityHeaders = solidityHeaders.concat(genSolidityFunctionHeaders(func));
                }
                outputFile = (0, templates_1.preamble)();
                for (_a = 0, solidityHeaders_1 = solidityHeaders; _a < solidityHeaders_1.length; _a++) {
                    fn = solidityHeaders_1[_a];
                    funcDefinition = generateSolidityFunction(parseFunctionDefinition(fn));
                    outputFile += funcDefinition;
                }
                outputFile += "\n// ********** TYPE CASTING ************* //\n";
                // generate casting functions
                for (_b = 0, _c = common_1.EInputType.concat('uint256', 'bytes memory'); _b < _c.length; _b++) {
                    fromType = _c[_b];
                    for (_d = 0, EInputType_2 = common_1.EInputType; _d < EInputType_2.length; _d++) {
                        toType = EInputType_2[_d];
                        if (fromType === toType) {
                            continue;
                        }
                        outputFile += (0, templates_1.AsTypeFunction)(fromType, toType);
                    }
                }
                for (_e = 0, EInputType_3 = common_1.EInputType; _e < EInputType_3.length; _e++) {
                    type = EInputType_3[_e];
                    functionName = "as".concat((0, templates_1.capitalize)(type));
                    testContract = (0, templates_1.AsTypeTestingContract)(type);
                    testContracts[functionName] = testContract[0];
                    testContractsAbis += testContract[1];
                    importLineHelper += "".concat((0, templates_1.capitalize)(functionName), "TestType,\n");
                }
                importLineHelper = importLineHelper.slice(0, -2) + " } from './abis';\n";
                outputFile += (0, templates_1.AsTypeFunction)("bool", "ebool");
                outputFile += (0, templates_1.PostFix)();
                outputFile += "\n// ********** OPERATOR OVERLOADING ************* //\n";
                // generate operator overloading
                common_1.ShorthandOperations.forEach(function (value) {
                    for (var _i = 0, EInputType_4 = common_1.EInputType; _i < EInputType_4.length; _i++) {
                        var encType = EInputType_4[_i];
                        if (!(0, common_1.valueIsEncrypted)(encType)) {
                            throw new Error("InputType mismatch");
                        }
                        if (!common_1.EComparisonType.includes(encType)) {
                            outputFile += (0, templates_1.OperatorOverloadDecl)(value.func, value.operator, encType, value.unary);
                        }
                    }
                });
                outputFile += "\n// ********** BINDING DEFS ************* //\n";
                common_1.EInputType.forEach(function (encryptedType) {
                    if (!common_1.EComparisonType.includes(encryptedType)) {
                        common_1.BindMathOperators.forEach(function (bindMathOp) {
                            if (common_1.ShorthandOperations.filter(function (value) { return value.func === bindMathOp; }).length === 0) {
                                // console.log(`${bindMathOp}`)
                                outputFile += (0, templates_1.BindingsWithoutOperator)(bindMathOp, encryptedType);
                            }
                        });
                        outputFile += (0, templates_1.BindingLibraryType)(encryptedType);
                        common_1.BindMathOperators.forEach(function (fnToBind) {
                            var foundFnDef = solidityHeaders.find(function (funcHeader) {
                                var fnDef = parseFunctionDefinition(funcHeader);
                                var input = fnDef.inputs[0];
                                if (!common_1.EInputType.includes(input)) {
                                    return false;
                                }
                                return (fnDef.funcName === fnToBind && fnDef.inputs.every(function (item) { return item === input; }));
                            });
                            if (foundFnDef) {
                                var fnDef = parseFunctionDefinition(foundFnDef);
                                outputFile += (0, templates_1.OperatorBinding)(fnDef.funcName, encryptedType, fnDef.inputs.length === 1);
                            }
                        });
                        outputFile += (0, templates_1.PostFix)();
                    }
                });
                return [4 /*yield*/, fs.promises.writeFile('FHE.sol', outputFile)];
            case 2:
                _h.sent();
                for (_f = 0, _g = Object.entries(testContracts); _f < _g.length; _f++) {
                    testContract = _g[_f];
                    fs.writeFileSync("../solidity/tests/contracts/".concat(testContract[0], ".sol"), testContract[1]);
                }
                fs.writeFileSync("../solidity/tests/abis.ts", (0, templates_1.genAbiFile)(testContractsAbis));
                console.log(importLineHelper);
                return [2 /*return*/];
        }
    });
}); };
main();
