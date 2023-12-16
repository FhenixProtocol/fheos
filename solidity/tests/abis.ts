import { Contract } from 'ethers';
export interface AddTestType extends Contract {
    add: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface ReencryptTestType extends Contract {
    reencrypt: (test: string, a: bigint, pubkey: Uint8Array) => Promise<Uint8Array>; // Adjust the method signature
}
export interface LteTestType extends Contract {
    lte: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface SubTestType extends Contract {
    sub: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface MulTestType extends Contract {
    mul: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface LtTestType extends Contract {
    lt: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface CmuxTestType extends Contract {
    cmux: (test: string,c: boolean, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface ReqTestType extends Contract {
    req: (test: string, a: bigint) => Promise<()>; // Adjust the method signature
}
export interface DivTestType extends Contract {
    div: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface GtTestType extends Contract {
    gt: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface GteTestType extends Contract {
    gte: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface RemTestType extends Contract {
    rem: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface AndTestType extends Contract {
    and: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface OrTestType extends Contract {
    or: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface XorTestType extends Contract {
    xor: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface EqTestType extends Contract {
    eq: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface NeTestType extends Contract {
    ne: (test: string, a: bigint, b: bigint) => Promise<bigint>;
}
export interface MinTestType extends Contract {
    min: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface MaxTestType extends Contract {
    max: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface ShlTestType extends Contract {
    shl: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface ShrTestType extends Contract {
    shr: (test: string, a: bigint, b: bigint) => Promise<bigint>; // Adjust the method signature
}
export interface NotTestType extends Contract {
    not: (test: string, a: bigint) => Promise<bigint>; // Adjust the method signature
}


