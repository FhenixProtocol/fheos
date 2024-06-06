import { BaseContract } from 'ethers';
export interface AddBenchType extends BaseContract {
    add: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface SealoutputBenchType extends BaseContract {
    sealoutput: (test: string, a: bigint, pubkey: Uint8Array) => Promise<string>;
}
export interface LteBenchType extends BaseContract {
    lte: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface SubBenchType extends BaseContract {
    sub: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface MulBenchType extends BaseContract {
    mul: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface LtBenchType extends BaseContract {
    lt: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface SelectBenchType extends BaseContract {
    select: (test: string, c: boolean, a: bigint, b: bigint) => Promise<bigint>;
}
export interface ReqBenchType extends BaseContract {
    req: (a: bytes) => Promise<{}>;
}
export interface DivBenchType extends BaseContract {
    div: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface GtBenchType extends BaseContract {
    gt: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface GteBenchType extends BaseContract {
    gte: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface RemBenchType extends BaseContract {
    rem: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface AndBenchType extends BaseContract {
    and: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface OrBenchType extends BaseContract {
    or: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface XorBenchType extends BaseContract {
    xor: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface EqBenchType extends BaseContract {
    eq: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface NeBenchType extends BaseContract {
    ne: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface MinBenchType extends BaseContract {
    min: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface MaxBenchType extends BaseContract {
    max: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface ShlBenchType extends BaseContract {
    shl: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface ShrBenchType extends BaseContract {
    shr: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface NotBenchType extends BaseContract {
    not: (test: string, a: bigint) => Promise<bigint>;
}


