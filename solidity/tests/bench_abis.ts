import { BaseContract } from 'ethers';
export interface AddBenchType extends BaseContract {
    add: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface SealoutputBenchType extends BaseContract {
    sealoutput: (test: string, a: bigint, pubkey: Uint8Array) => Promise<string>;
}
export interface LteBenchType extends BaseContract {
    lte: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface SubBenchType extends BaseContract {
    sub: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface MulBenchType extends BaseContract {
    mul: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface LtBenchType extends BaseContract {
    lt: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface SelectBenchType extends BaseContract {
    select: (test: string, c: boolean, a: bigint, b: bigint) => Promise<bigint>;
}
export interface ReqBenchType extends BaseContract {
    req: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface DivBenchType extends BaseContract {
    div: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface GtBenchType extends BaseContract {
    gt: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface GteBenchType extends BaseContract {
    gte: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface RemBenchType extends BaseContract {
    rem: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface AndBenchType extends BaseContract {
    and: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface OrBenchType extends BaseContract {
    or: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface XorBenchType extends BaseContract {
    xor: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface EqBenchType extends BaseContract {
    eq: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface NeBenchType extends BaseContract {
    ne: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface MinBenchType extends BaseContract {
    min: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface MaxBenchType extends BaseContract {
    max: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface ShlBenchType extends BaseContract {
    shl: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface ShrBenchType extends BaseContract {
    shr: (_a: bytes, _b: bytes) => Promise<bigint>;
}
export interface NotBenchType extends BaseContract {
    not: (_a: bytes, _b: bytes) => Promise<bigint>;
}


