import {
  EInputType,
  SEALING_FUNCTION_NAME,
  AllowedOperations,
  capitalize,
} from "../common";

const toVarSuffix = (inputType: string) => capitalize(inputType.slice(1).replace("uint", ""));
const toInType = (inputType: string) => "inE" + inputType.slice(1);
const toInTypeFunction = (inputType: string) => toInType(inputType) + " calldata";
const toAsType = (inputType: string) => "asE" + inputType.slice(1);

const IsOperationAllowed = (
  functionName: string,
  dataType: string,
): boolean => {
  const inputIdx= EInputType.indexOf(dataType)
  const regexes = AllowedOperations[inputIdx];
  for (let regex of regexes) {
    if (!new RegExp(regex).test(functionName.toLowerCase())) {
      return false;
    }
  }

  return true;
};

export function generateBenchContract(
  name: string,
  body: string,
  importStatement: string = `import {ebool, euint8} from "../../../FHE.sol";`,
) {
  return `// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
${importStatement}

contract ${capitalize(name)}Bench {
${body}
}
`;
}

export function benchContract1Arg(name: string) {
  let importTypes = ``;
  let privateVarsA = ``;
  let funcLoad = "";
  let funcBench = "";

  for (let inputType of EInputType) {
    if (IsOperationAllowed(name, inputType)) {
      importTypes += "\n\t" + inputType + ", " + toInType(inputType) + ",";
      privateVarsA += `\t${inputType} internal a${toVarSuffix(inputType)};\n`;
      funcLoad += `
    function load${toVarSuffix(inputType)}(${toInTypeFunction(inputType)} _a) public {
        a${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_a);
    }`;

      funcBench += `
    function bench${capitalize(name)}${toVarSuffix(inputType)}() public view {
        FHE.${name}(a${toVarSuffix(inputType)});
    }`;
    }
  }

  const func = privateVarsA + funcLoad + "\n" + funcBench;
  importTypes = importTypes.slice(0, -1); // remove last comma
  const importStatement = `import {${importTypes}
} from "../../../FHE.sol";`;

  return generateBenchContract(name, func, importStatement);
}

export function benchContract2Arg(name: string) {
  let importTypes = ``;
  let privateVarsA = ``;
  let privateVarsB = ``;
  let funcLoad = "";
  let funcBench = "";

  for (let inputType of EInputType) {
    if (IsOperationAllowed(name, inputType)) {
      importTypes += "\n\t" + inputType + ", " + toInType(inputType) + ",";
      privateVarsA += `\t${inputType} internal a${toVarSuffix(inputType)};\n`;
      privateVarsB += `\t${inputType} internal b${toVarSuffix(inputType)};\n`;
      funcLoad += `
    function load${toVarSuffix(inputType)}(${toInTypeFunction(inputType)} _a, ${toInTypeFunction(inputType)} _b) public {
        a${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_a);
        b${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_b);
    }`;

      funcBench += `
    function bench${capitalize(name)}${toVarSuffix(inputType)}() public view {
        FHE.${name}(a${toVarSuffix(inputType)}, b${toVarSuffix(inputType)});
    }`;
    }
  }

  const func = privateVarsA + "\n" + privateVarsB + funcLoad + "\n" + funcBench;
  importTypes = importTypes.slice(0, -1); // remove last comma
  const importStatement = `import {${importTypes}
} from "../../../FHE.sol";`;

  return generateBenchContract(name, func, importStatement);
}

export function benchContract3Arg(name: string) {
  let importTypes = ``;
  let privateVarsA = `\tebool internal control;\n\n`;
  let privateVarsB = ``;
  let funcLoad = "";
  let funcBench = "";

  for (let inputType of EInputType) {
    if (IsOperationAllowed(name, inputType)) {
      importTypes += "\n\t" + inputType + ", " + toInType(inputType) + ",";
      privateVarsA += `\t${inputType} internal a${toVarSuffix(inputType)};\n`;
      privateVarsB += `\t${inputType} internal b${toVarSuffix(inputType)};\n`;
      funcLoad += `
    function load${toVarSuffix(inputType)}(inEbool calldata _control, ${toInTypeFunction(inputType)} _a, ${toInTypeFunction(inputType)} _b) public {
        control = FHE.asEbool(_control);
        a${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_a);
        b${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_b);
    }`;

      funcBench += `
    function bench${capitalize(name)}${toVarSuffix(inputType)}() public view {
        FHE.${name}(control, a${toVarSuffix(inputType)}, b${toVarSuffix(inputType)});
    }`;
    }
  }

  const func = privateVarsA + "\n" + privateVarsB + funcLoad + "\n" + funcBench;
  importTypes = importTypes.slice(0, -1); // remove last comma
  const importStatement = `import {${importTypes}
} from "../../../FHE.sol";`;

  return generateBenchContract(name, func, importStatement);
}

export function benchContractReencrypt() {
  let importTypes = ``;
  let privateVarsA = `\tbytes32 internal pubkey;\n`;
  let funcLoad = "";
  let funcBench = "";

  for (let inputType of EInputType) {
    if (IsOperationAllowed(SEALING_FUNCTION_NAME, inputType)) {
      importTypes += "\n\t" + inputType + ", " + toInType(inputType) + ",";
      privateVarsA += `\t${inputType} internal a${toVarSuffix(inputType)};\n`;
      funcLoad += `
    function load${toVarSuffix(inputType)}(${toInTypeFunction(inputType)} _a, bytes32 _pubkey) public {
        a${toVarSuffix(inputType)} = FHE.${toAsType(inputType)}(_a);
        pubkey = _pubkey;
    }`;

      funcBench += `
    function bench${capitalize(SEALING_FUNCTION_NAME)}${toVarSuffix(inputType)}() public view {
        FHE.${SEALING_FUNCTION_NAME}(a${toVarSuffix(inputType)}, pubkey);
    }`;
    }
  }

  const func = privateVarsA + funcLoad + "\n" + funcBench;
  importTypes = importTypes.slice(0, -1); // remove last comma
  const importStatement = `import {${importTypes}
} from "../../../FHE.sol";`;

  return generateBenchContract(SEALING_FUNCTION_NAME, func, importStatement);
}

export function AsTypeBenchmarkContract(type: string) {
  let privateVarsA = "";
  let importTypes = "";
  let loads = "";
  let casts = "";

  // Although casts from eaddress to types with < 256 bits are possible, we don't want to bench them.
  let eaddressAllowedTypes = ["euint256"];
  let fromTypeCollection = type === "eaddress" ? eaddressAllowedTypes : EInputType;

  for (let fromType of fromTypeCollection) {
    importTypes += "\n\t" + fromType + ", " + toInType(fromType) + ",";

    privateVarsA += `\t${fromType} internal a${toVarSuffix(fromType)};\n`;

    loads += `\n\tfunction load${toVarSuffix(fromType)}(${toInTypeFunction(fromType)} _a) public {
        a${toVarSuffix(fromType)} = FHE.${toAsType(fromType)}(_a);
    }`;
    casts += `\n\tfunction benchCast${capitalize(fromType)}To${capitalize(type)}() public view {
        FHE.${toAsType(type)}(a${toVarSuffix(fromType)});
    }`;
  }

  // deal with casting from built-in types
  let builtInTypes = {"uint256": "Uint256", "bytes memory": "Bytes"};
  for (const [builtInType, varSuffix] of Object.entries(builtInTypes)) {
    privateVarsA += `\t${builtInType} internal a${varSuffix};\n`;

    loads += `\n\tfunction load${varSuffix}(${builtInType} _a) public {
        a${varSuffix} = _a;
    }`;
    casts += `\n\tfunction benchCast${varSuffix}To${capitalize(type)}() public view {
        FHE.${toAsType(type)}(a${varSuffix});
    }`;
  }

  const body = privateVarsA + loads + "\n" + casts;
  importTypes = importTypes.slice(0, -1); // remove last comma
  const importStatement = `import {${importTypes}
} from "../../../FHE.sol";`;

  return generateBenchContract(`As${capitalize(type)}`, body, importStatement);
}