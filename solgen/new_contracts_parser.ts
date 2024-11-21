import * as fs from "fs";
import * as path from "path";
import { AllTypes } from "./common";

export interface FunctionAnalysis {
  name: string;
  inputTypes: AllTypes[];
  originalText: string;
}

const parseGoFunction = (functionText: string): FunctionAnalysis | null => {
  // Extract function name
  const nameMatch = functionText.match(/func\s+(\w+)/);
  if (!nameMatch) return null;

  const name = nameMatch[1];

  // Default to encrypted for input types
  const inputTypes: AllTypes[] = ['encrypted', 'encrypted'];

  return {
    name,
    inputTypes,
    originalText: functionText
  };
};

export const getFunctionsFromGo = async (filePath: string): Promise<FunctionAnalysis[]> => {
  const fullPath = path.join(__dirname, filePath);
  const content = await fs.promises.readFile(fullPath, 'utf8');

  // Find all function declarations with their comments
  const functionRegex = /\/\/[^\n]*\n(\/\/[^\n]*\n)*func\s+\w+[^{]+{[^}]*}/g;
  const matches = content.matchAll(functionRegex);

  const functions: FunctionAnalysis[] = [];

  for (const match of matches) {
    const functionText = match[0];
    const parsed = parseGoFunction(functionText);
    if (parsed) {
      functions.push(parsed);
    }
  }

  return functions;
}; 