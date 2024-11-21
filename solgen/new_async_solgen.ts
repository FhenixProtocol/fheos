// import { getFunctionsFromGo } from "./new_contracts_parser";
// import * as fs from "fs";
// import {
//   AllTypes,
//   EInputType,
//   valueIsEncrypted,
//   capitalize,
// } from "./common";

// interface AsyncSolgenConfig {
//   sameInputs: boolean;
//   forcedReturnType?: string;
// }

// interface FunctionMetadata {
//   name: string;
//   config: AsyncSolgenConfig;
//   inputTypes: AllTypes[];
// }

// const parseAsyncSolgenComments = (functionText: string): AsyncSolgenConfig => {
//   const config: AsyncSolgenConfig = {
//     sameInputs: false
//   };

//   // Parse async-solgen-same-inputs
//   const sameInputsMatch = functionText.match(/async-solgen-same-inputs:(true|false)/);
//   if (sameInputsMatch) {
//     config.sameInputs = sameInputsMatch[1] === 'true';
//   }

//   // Parse async-solgen-return-type
//   const returnTypeMatch = functionText.match(/async-solgen-return-type:(\w+)/);
//   if (returnTypeMatch) {
//     config.forcedReturnType = returnTypeMatch[1];
//   }

//   return config;
// };

// const generateFunctionHeader = (metadata: FunctionMetadata): string => {
//   const { name, config, inputTypes } = metadata;
  
//   let docString = `
//     /// @notice Performs the ${name} operation on encrypted inputs
//     /// @dev This is an async operation that will be executed by the FHE network
//     `;

//   if (config.sameInputs) {
//     docString += `/// @dev Both inputs must be of the same type\n`;
//   }

//   docString += `/// @param lhs - Left hand side operand
//     /// @param rhs - Right hand side operand
//     /// @return Result of type ${config.forcedReturnType || inputTypes[0]}
//     `;

//   return docString;
// };

// const generateSolidityFunction = (metadata: FunctionMetadata): string => {
//   const { name, config, inputTypes } = metadata;
  
//   // Generate function parameters
//   let params = '';
//   if (inputTypes.length == 1) {
//     params = `${inputTypes[0]} value`;
//   } else if (inputTypes.length == 2) {
//     if (config.sameInputs) {
//       // Both inputs must be the same type
//       params = `${inputTypes[0]} lhs, ${inputTypes[0]} rhs`;
//     } else {
//       params = `${inputTypes[0]} lhs, ${inputTypes[1]} rhs`;
//     }
//   }

//   // Determine return type
//   const returnType = config.forcedReturnType || inputTypes[0];

//   // Generate function header with documentation
//   let functionBody = generateFunctionHeader(metadata);

//   // Generate function implementation
//   functionBody += `function ${name}(${params}) internal returns (${returnType}) {
//         if (!isInitialized(lhs)) {
//             lhs = as${capitalize(inputTypes[0])}(0);
//         }
//         if (!isInitialized(rhs)) {
//             rhs = as${capitalize(inputTypes[0])}(0);
//         }
        
//         uint256 unwrappedInput1 = ${inputTypes[0]}.unwrap(lhs);
//         uint256 unwrappedInput2 = ${inputTypes[0]}.unwrap(rhs);

//         uint256 result = calcBinaryPlaceholderValueHash(unwrappedInput1, unwrappedInput2, FunctionId.${name});
//         ITaskManager(TASK_MANAGER_ADDRESS).createTask(result, "${name}", unwrappedInput1, unwrappedInput2);
//         return ${returnType}.wrap(result);
//     }`;

//   return functionBody;
// }

// const generateSolidityCode = (functions: FunctionMetadata[]): string => {
//   let code = `// SPDX-License-Identifier: BSD-3-Clause-Clear
// pragma solidity >=0.8.19 <0.9.0;

// import {FheOps} from "./FheOS.sol";
// import {Common} from "./Common.sol";

// /// @title AsyncFHE Library
// /// @notice Provides async FHE operations that are executed by the FHE network
// /// @dev All operations in this library are async and return temporary values
// library AsyncFHE {
//     address constant TASK_MANAGER_ADDRESS = address(129);

//     enum FunctionId {
//         ${functions.map(f => f.name).join(',\n        ')}
//     }

// `;

//   // Generate all functions
//   for (const fn of functions) {
//     code += generateSolidityFunction(fn);
//     code += '\n';
//   }

//   code += '}\n';
//   return code;
// };

// const main = async () => {
//   console.log("Starting async solgen...");
  
//   // Read and parse the Go file
//   const goFunctions = await getFunctionsFromGo("../precompiles/contracts.go");

//   // Filter functions with async-solgen comments and parse their configs
//   const asyncFunctions: FunctionMetadata[] = goFunctions
//     .filter(fn => fn.originalText.includes('async-solgen'))
//     .map(fn => ({
//       name: fn.name,
//       config: parseAsyncSolgenComments(fn.originalText),
//       inputTypes: fn.inputTypes
//     }));

//   // Generate Solidity code
//   const outputFile = generateSolidityCode(asyncFunctions);

//   // Write output file
//   await fs.promises.writeFile("AsyncFHE.sol", outputFile);
//   console.log("Generated AsyncFHE.sol");
// };

// main().catch(console.error); 