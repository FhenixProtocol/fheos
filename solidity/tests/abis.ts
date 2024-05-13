import { BaseContract } from 'ethers';
export interface AddTestType extends BaseContract {
    add: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface SealoutputTestType extends BaseContract {
    sealoutput: (test: string, a: bigint, pubkey: Uint8Array) => Promise<string>;
}
export interface LteTestType extends BaseContract {
    lte: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface SubTestType extends BaseContract {
    sub: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface MulTestType extends BaseContract {
    mul: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface LtTestType extends BaseContract {
    lt: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface SelectTestType extends BaseContract {
    select: (test: string, c: boolean, a: bigint, b: bigint) => Promise<bigint>;
}
export interface ReqTestType extends BaseContract {
    req: (test: string, a: bigint) => Promise<{}>;
}
export interface DivTestType extends BaseContract {
    div: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface GtTestType extends BaseContract {
    gt: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface GteTestType extends BaseContract {
    gte: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface RemTestType extends BaseContract {
    rem: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface AndTestType extends BaseContract {
    and: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface OrTestType extends BaseContract {
    or: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface XorTestType extends BaseContract {
    xor: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface EqTestType extends BaseContract {
    eq: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface NeTestType extends BaseContract {
    ne: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface MinTestType extends BaseContract {
    min: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface MaxTestType extends BaseContract {
    max: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface ShlTestType extends BaseContract {
    shl: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface ShrTestType extends BaseContract {
    shr: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface NotTestType extends BaseContract {
    not: (test: string, a: bigint) => Promise<bigint>;
}
export interface AsEboolTestType extends BaseContract {
    castFromEuint8ToEbool: (val: bigint, test: string) => Promise<boolean>;
    castFromEuint16ToEbool: (val: bigint, test: string) => Promise<boolean>;
    castFromEuint32ToEbool: (val: bigint, test: string) => Promise<boolean>;
    castFromEuint64ToEbool: (val: bigint, test: string) => Promise<boolean>;
    castFromEuint128ToEbool: (val: bigint, test: string) => Promise<boolean>;
    castFromEuint256ToEbool: (val: bigint, test: string) => Promise<boolean>;
    castFromPlaintextToEbool: (val: bigint) => Promise<boolean>;
    castFromPreEncryptedToEbool: (val: Uint8Array) => Promise<boolean>;
}
export interface AsEuint8TestType extends BaseContract {
    castFromEboolToEuint8: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint16ToEuint8: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint32ToEuint8: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint64ToEuint8: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint128ToEuint8: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint256ToEuint8: (val: bigint, test: string) => Promise<bigint>;
    castFromPlaintextToEuint8: (val: bigint) => Promise<bigint>;
    castFromPreEncryptedToEuint8: (val: Uint8Array) => Promise<bigint>;
}
export interface AsEuint16TestType extends BaseContract {
    castFromEboolToEuint16: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint8ToEuint16: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint32ToEuint16: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint64ToEuint16: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint128ToEuint16: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint256ToEuint16: (val: bigint, test: string) => Promise<bigint>;
    castFromPlaintextToEuint16: (val: bigint) => Promise<bigint>;
    castFromPreEncryptedToEuint16: (val: Uint8Array) => Promise<bigint>;
}
export interface AsEuint32TestType extends BaseContract {
    castFromEboolToEuint32: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint8ToEuint32: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint16ToEuint32: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint64ToEuint32: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint128ToEuint32: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint256ToEuint32: (val: bigint, test: string) => Promise<bigint>;
    castFromPlaintextToEuint32: (val: bigint) => Promise<bigint>;
    castFromPreEncryptedToEuint32: (val: Uint8Array) => Promise<bigint>;
}
export interface AsEuint64TestType extends BaseContract {
    castFromEboolToEuint64: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint8ToEuint64: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint16ToEuint64: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint32ToEuint64: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint128ToEuint64: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint256ToEuint64: (val: bigint, test: string) => Promise<bigint>;
    castFromPlaintextToEuint64: (val: bigint) => Promise<bigint>;
    castFromPreEncryptedToEuint64: (val: Uint8Array) => Promise<bigint>;
}
export interface AsEuint128TestType extends BaseContract {
    castFromEboolToEuint128: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint8ToEuint128: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint16ToEuint128: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint32ToEuint128: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint64ToEuint128: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint256ToEuint128: (val: bigint, test: string) => Promise<bigint>;
    castFromPlaintextToEuint128: (val: bigint) => Promise<bigint>;
    castFromPreEncryptedToEuint128: (val: Uint8Array) => Promise<bigint>;
}
export interface AsEuint256TestType extends BaseContract {
    castFromEboolToEuint256: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint8ToEuint256: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint16ToEuint256: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint32ToEuint256: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint64ToEuint256: (val: bigint, test: string) => Promise<bigint>;
    castFromEuint128ToEuint256: (val: bigint, test: string) => Promise<bigint>;
    castFromEaddressToEuint256: (val: bigint, test: string) => Promise<bigint>;
    castFromPlaintextToEuint256: (val: bigint) => Promise<bigint>;
    castFromPreEncryptedToEuint256: (val: Uint8Array) => Promise<bigint>;
}
export interface AsEaddressTestType extends BaseContract {
    castFromEuint256ToEaddress: (val: bigint, test: string) => Promise<bigint>;
    castFromPlaintextToEaddress: (val: bigint) => Promise<bigint>;
    castFromPreEncryptedToEaddress: (val: Uint8Array) => Promise<bigint>;
}


