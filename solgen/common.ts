export const EInputType = ['ebool', 'euint8', 'euint16', 'euint32'];
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
    operator: string,
    func: string,
    unary: boolean
}

export const ShorthandOperations: OperatorMap[] =
[
    {
        func: 'add',
        operator: '+',
        unary: false
    },
    {
        func: 'sub',
        operator: '-',
        unary: false
    },
    {
        func: 'mul',
        operator: '*',
        unary: false
    },
    {
        func: 'div',
        operator: '/',
        unary: false
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
    // {
    //     func: 'xor',
    //     operator: '^',
    //     unary: true
    // },
    {
        func: 'or',
        operator: '|',
        unary: false,
    },
    {
        func: 'and',
        operator: '&',
        unary: false
    }
]

export const BindMathOperators = ['add', 'mul', 'div', 'sub', 'eq', 'and', 'or']

export const valueIsEncrypted = (value: string): value is EUintType => {
    return EInputType.includes(value);
}

export const valueIsPlaintext = (value: string): value is PlaintextType => {
    return EPlaintextType.includes(value);
}