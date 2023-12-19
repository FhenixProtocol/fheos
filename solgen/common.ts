export const EInputType = ['ebool', 'euint8', 'euint16', 'euint32'];
export const EComparisonType = ['ebool'];
export const EPlaintextType = ['bool', 'uint8', 'uint16', 'uint32', 'uint64', 'uint128', 'uint256'];
export type EUintType = 'ebool' | 'euint8' | 'euint16' | 'euint32';
export type PlaintextType = 'bool' | 'uint8' | 'uint16' | 'uint32' | 'uint64' | 'uint128' | 'uint256';
export type AllTypes = PlaintextType | EUintType | 'bytes memory' | "bytes32" | "uint8" | "encrypted" | "plaintext" | "none";

export const UnderlyingTypes: Record<EUintType, string> = {
    euint8: 'uint256',
    euint16: 'uint256',
    euint32: 'uint256',
    ebool: 'uint256',
};

interface OperatorMap {
    operator: string | null,
    func: string,
    unary: boolean,
    returnsBool: boolean,
}

export const ShorthandOperations: OperatorMap[] =
[
    {
        func: 'add',
        operator: '+',
        unary: false,
        returnsBool: false,
    },
    {
        func: 'sub',
        operator: '-',
        unary: false,
        returnsBool: false,
    },
    {
        func: 'mul',
        operator: '*',
        unary: false,
        returnsBool: false,
    },
    {
        func: 'div',
        operator: '/',
        unary: false,
        returnsBool: false,
    },
    // {
    //     func: 'not',
    //     operator: '~',
    //     unary: true
    // },
    // {
    //     func: 'neg',
    //     operator: '-',
    //     unary: true
    // },
    {
        func: 'or',
        operator: '|',
        unary: false,
        returnsBool: false,
    },
    {
        func: 'and',
        operator: '&',
        unary: false,
        returnsBool: false,
    },
    {
        func: 'xor',
        operator: '^',
        unary: false,
        returnsBool: false,
    },
    {
        func: 'gt',
        operator: '>',
        unary: false,
        returnsBool: true,
    },
    {
        func: 'gte',
        operator: '>=',
        unary: false,
        returnsBool: true,
    },
    {
        func: 'lt',
        operator: '<',
        unary: false,
        returnsBool: true,
    },
    {
        func: 'lte',
        operator: '<=',
        unary: false,
        returnsBool: true,
    },
    {
        func: 'rem',
        operator: '%',
        unary: false,
        returnsBool: false,
    },
    {
        func: 'max',
        operator: null,
        unary: false,
        returnsBool: false,
    },
    {
        func: 'min',
        operator: null,
        unary: false,
        returnsBool: false,
    },
    {
        func: 'eq',
        operator: '==',
        unary: false,
        returnsBool: true,
    },
    {
        func: 'ne',
        operator: '!=',
        unary: false,
        returnsBool: true,
    }
]

export const BindMathOperators = ['add', 'mul', 'div', 'sub', 'eq', 'ne', 'and', 'or', 'xor', 'gt', 'gte', 'lt', 'lte', 'rem', 'max', 'min'];

export const valueIsEncrypted = (value: string): value is EUintType => {
    return EInputType.includes(value);
}

export const valueIsPlaintext = (value: string): value is PlaintextType => {
    return EPlaintextType.includes(value);
}