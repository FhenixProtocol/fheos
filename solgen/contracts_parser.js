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
exports.getFunctionsFromGo = void 0;
var fs = require("fs");
var path = require("path");
// helps us know how many input parameters there are
var specificFunctions = [
    { name: 'get3VerifiedOperands', amount: 3, paramTypes: ["encrypted", "encrypted", "encrypted"] },
    { name: 'get2VerifiedOperands', amount: 2, paramTypes: ["encrypted", "uint8", "plaintext"] },
    { name: 'getCiphertext', amount: 1, paramTypes: ["encrypted"] },
    { name: 'tfhe.UintType(', amount: 1, paramTypes: ["encrypted"] }
];
function analyzeGoFile(filePath) {
    return __awaiter(this, void 0, void 0, function () {
        var fileContent, lines, highLevelFunctionRegex, solgenCommentRegex, solgenReturnsComment, solgenInputPlaintextComment, solgenOutputPlaintextComment, solgenInput2Comment, solgenBooleanMathOp, specificFunctionAnalysis, isInsideHighLevelFunction, isBooleanMathOp, braceDepth, funcName, returnType, inputs, _i, lines_1, line, trimmedLine, _a, specificFunctions_1, keyfn, needsSameType, amount;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0: return [4 /*yield*/, fs.promises.readFile(filePath, 'utf-8')];
                case 1:
                    fileContent = _b.sent();
                    lines = fileContent.split('\n');
                    highLevelFunctionRegex = /func\s+/;
                    solgenCommentRegex = /solgen:/;
                    solgenReturnsComment = / return /;
                    solgenInputPlaintextComment = /input plaintext/;
                    solgenOutputPlaintextComment = /output plaintext/;
                    solgenInput2Comment = /input2 /;
                    solgenBooleanMathOp = /bool math/;
                    specificFunctionAnalysis = [];
                    isInsideHighLevelFunction = false;
                    isBooleanMathOp = false;
                    braceDepth = 0;
                    funcName = "";
                    returnType = undefined;
                    inputs = [];
                    for (_i = 0, lines_1 = lines; _i < lines_1.length; _i++) {
                        line = lines_1[_i];
                        trimmedLine = line.trim();
                        //console.log(`testing: ${trimmedLine}`)
                        if (isInsideHighLevelFunction) {
                            if (solgenCommentRegex.test(trimmedLine)) {
                                if (solgenReturnsComment.test(trimmedLine)) {
                                    returnType = trimmedLine.split('return')[1].trim();
                                }
                                if (solgenInputPlaintextComment.test(trimmedLine)) {
                                    inputs = inputs.map(function () { return "plaintext"; });
                                }
                                if (solgenInput2Comment.test(trimmedLine)) {
                                    // @ts-ignore
                                    inputs[1] = trimmedLine.split('input2 ')[1].trim();
                                }
                                if (solgenOutputPlaintextComment.test(trimmedLine)) {
                                    returnType = 'plaintext';
                                }
                                if (solgenBooleanMathOp.test(trimmedLine)) {
                                    isBooleanMathOp = true;
                                }
                            }
                            braceDepth += (trimmedLine.match(/\{/g) || []).length;
                            braceDepth -= (trimmedLine.match(/}/g) || []).length;
                            //console.log(`brace depth: ${braceDepth}`)
                            // Check if we've exited the high-level function
                            if (braceDepth === 0) {
                                isInsideHighLevelFunction = false;
                                continue;
                            }
                            // Look for specific functions within high-level function
                            for (_a = 0, specificFunctions_1 = specificFunctions; _a < specificFunctions_1.length; _a++) {
                                keyfn = specificFunctions_1[_a];
                                if (trimmedLine.includes(keyfn.name)) {
                                    needsSameType = /lhs.UintType\s+!=\s+rhs.UintType/.test(trimmedLine);
                                    amount = keyfn.amount;
                                    if (funcName === "reencrypt") {
                                        // console.log(`func name: ${funcName}`)
                                        amount = 2;
                                        returnType = "bytes memory";
                                        needsSameType = false;
                                        inputs = ["encrypted", "bytes32"];
                                    }
                                    // we generate these manually for now
                                    if (['trivialencrypt', 'cast', 'verify'].includes(funcName.toLowerCase())) {
                                        continue;
                                    }
                                    specificFunctionAnalysis.push({
                                        name: funcName,
                                        paramsCount: amount,
                                        needsSameType: needsSameType,
                                        returnType: returnType,
                                        inputTypes: inputs.slice(0, amount),
                                        isBooleanMathOp: isBooleanMathOp
                                    });
                                }
                            }
                        }
                        else if (highLevelFunctionRegex.test(trimmedLine)) {
                            // console.log(trimmedLine)
                            returnType = "encrypted";
                            inputs = ["encrypted", "encrypted", "encrypted"];
                            funcName = trimmedLine.split(' ')[1].split('(')[0].toLowerCase();
                            // If we match the high-level function, set the flag and initialize brace counting
                            isInsideHighLevelFunction = true;
                            isBooleanMathOp = false;
                            braceDepth = 1; // starts with the opening brace of the function
                        }
                    }
                    return [2 /*return*/, specificFunctionAnalysis.length > 0 ? specificFunctionAnalysis : null];
            }
        });
    });
}
var getFunctionsFromGo = function (file) { return __awaiter(void 0, void 0, void 0, function () {
    var goFilePath, result;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0:
                goFilePath = path.join(__dirname, file);
                return [4 /*yield*/, analyzeGoFile(goFilePath)];
            case 1:
                result = _a.sent();
                if (result) {
                    // console.log(result)
                    return [2 /*return*/, result];
                }
                else {
                    throw new Error('No specific function found.');
                }
                return [2 /*return*/];
        }
    });
}); };
exports.getFunctionsFromGo = getFunctionsFromGo;
