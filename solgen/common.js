"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.valueIsPlaintext = exports.valueIsEncrypted = exports.BindMathOperators = exports.ShorthandOperations = exports.UnderlyingTypes = exports.EPlaintextType = exports.EInputType = void 0;
exports.EInputType = ['ebool', 'euint8', 'euint16', 'euint32'];
exports.EPlaintextType = ['bool', 'uint8', 'uint16', 'uint32', 'uint64', 'uint128', 'uint256'];
exports.UnderlyingTypes = {
    euint8: 'uint256',
    euint16: 'uint256',
    euint32: 'uint256',
    ebool: 'uint256',
};
exports.ShorthandOperations = [
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
];
exports.BindMathOperators = ['add', 'mul', 'div', 'sub', 'eq', 'and', 'or'];
var valueIsEncrypted = function (value) {
    return exports.EInputType.includes(value);
};
exports.valueIsEncrypted = valueIsEncrypted;
var valueIsPlaintext = function (value) {
    return exports.EPlaintextType.includes(value);
};
exports.valueIsPlaintext = valueIsPlaintext;
