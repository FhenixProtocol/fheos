import * as fs from 'fs';
import * as path from 'path';

type ParamTypes = "encrypted" | "uint8" | "plaintext" | "bytes32";

interface FunctionAnalysis {
    name: string;
    paramsCount: number;
    needsSameType: boolean;
    inputTypes: ParamTypes[],
    returnType?: string;
    isBooleanMathOp: boolean;
}

// helps us know how many input parameters there are
const specificFunctions = [
    {name: 'get3VerifiedOperands', amount: 3, paramTypes: ["encrypted", "encrypted", "encrypted"]},
    {name: 'get2VerifiedOperands', amount: 2, paramTypes: ["encrypted", "uint8", "plaintext"]},
    {name: 'getCiphertext', amount: 1, paramTypes: ["encrypted"]},
    {name: 'tfhe.UintType(', amount: 1, paramTypes: ["encrypted"]}
];

async function analyzeGoFile(filePath: string): Promise<FunctionAnalysis[] | null> {
    const fileContent = await fs.promises.readFile(filePath, 'utf-8');
    const lines = fileContent.split('\n');

    const highLevelFunctionRegex = /func\s+/;
    const solgenCommentRegex = /solgen:/;
    const solgenReturnsComment = / return /;
    const solgenInputPlaintextComment = /input plaintext/;
    const solgenOutputPlaintextComment = /output plaintext/;
    const solgenInput2Comment = /input2 /;
    const solgenBooleanMathOp = /bool math/;
    const specificFunctionAnalysis: FunctionAnalysis[] = [];

    let isInsideHighLevelFunction = false;
    let isBooleanMathOp = false;
    let braceDepth = 0;
    let funcName = "";
    let returnType = undefined;
    let inputs: ParamTypes[] = [];

    for (const line of lines) {
        const trimmedLine = line.trim();
        //console.log(`testing: ${trimmedLine}`)
        if (isInsideHighLevelFunction) {
            if (solgenCommentRegex.test(trimmedLine)) {
                if (solgenReturnsComment.test(trimmedLine)) {
                    returnType = trimmedLine.split('return')[1].trim();
                }
                if (solgenInputPlaintextComment.test(trimmedLine)) {
                    inputs = inputs.map(() => {return "plaintext"});
                }
                if (solgenInput2Comment.test(trimmedLine)) {
                    // @ts-ignore
                    inputs[1] = trimmedLine.split('input2 ')[1].trim();
                }
                if(solgenOutputPlaintextComment.test(trimmedLine)) {
                    returnType = 'plaintext';
                }
                if(solgenBooleanMathOp.test(trimmedLine)) {
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
            for (const keyfn of specificFunctions) {
                if (trimmedLine.includes(keyfn.name)) {
                    let needsSameType = /lhs.UintType\s+!=\s+rhs.UintType/.test(trimmedLine);
                    let amount = keyfn.amount;
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
                        isBooleanMathOp : isBooleanMathOp
                    });
                }
            }
        } else if (highLevelFunctionRegex.test(trimmedLine)) {
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

    return specificFunctionAnalysis.length > 0 ? specificFunctionAnalysis : null;
}

export const getFunctionsFromGo = async (file: string) => {
    let goFilePath = path.join(__dirname, file);

    let result = await analyzeGoFile(goFilePath);
    if (result) {
        // console.log(result)
        return result
    } else {
        throw new Error('No specific function found.');
    }
}