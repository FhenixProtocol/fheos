import { ethers } from 'hardhat';
import { createFheInstance } from './utils';
import { AddTestType,
    ReencryptTestType,
    LteTestType,
    SubTestType,
    MulTestType,
    LtTestType,
    CmuxTestType,
    ReqTestType,
    DivTestType,
    GtTestType,
    GteTestType,
    RemTestType,
    AndTestType,
    OrTestType,
    XorTestType,
    EqTestType,
    NeTestType,
    MinTestType,
    MaxTestType,
    ShlTestType,
    ShrTestType,
    NotTestType } from './abis';

import {BaseContract} from "ethers";

const bUnderflow = 4;


const deployContractFromSigner = async (con: any, signer: any, nonce?: number) => {
    return await con.deploy({
        from: signer,
        args: [],
        log: true,
        skipIfAlreadyDeployed: false,
        nonce,
    });
}

const syncNonce = async (con: any, signer: any, stateNonce: number) => {
    console.log(`Syncing nonce to ${stateNonce}`);
    try {
        await deployContractFromSigner(con, signer, stateNonce);
    } catch(e) {
        console.log("Fixed nonce issue");
    }

    return await deployContractFromSigner(con, signer);
}
const deployContract = async (contractName: string) => {
    const [signer] = await ethers.getSigners();
    const con = await ethers.getContractFactory(contractName);
    let deployedContract : BaseContract;
    try {
         deployedContract = await deployContractFromSigner(con, signer);

    } catch (e) {
        if (`${e}`.includes("nonce too")) {
            // find last occurence of ": " in e and get the number that comes after
            const match = `${e}`.match(/state: (\d+)/);
            const stateNonce = match ? parseInt(match[1], 10) : null;
            if (stateNonce === null) {
                throw new Error("Could not find nonce in error");
            }

            deployedContract = await syncNonce(con, signer, stateNonce);
        }
    }

    const contract = deployedContract.connect(signer);
    return contract;
}

const getFheContract = async (contractAddress: string) => {
    const fheContract = await createFheInstance(contractAddress);
    return fheContract;
}
describe('Test Add', () =>  {
    const aOverflow8 = 2 ** 8 - 1;
    const aOverflow16 = 2 ** 16 - 1;
    const aOverflow32 = 2 ** 32 - 1;
    const aUnderflow = 1;

    const bOverflow8 = aOverflow8 - 1;
    const bOverflow16 = aOverflow16 - 1;
    const bOverflow32 = aOverflow32 - 1;
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('AddTest') as AddTestType;
        expect(contract).toBeTruthy();

    });

    const testCases = [
        {
            function: "add(euint8,euint8)",
            cases: [
                {a: 1, b: 2, expectedResult: 3, name: ""},
                {
                    a: aOverflow8,
                    b: bOverflow8,
                    expectedResult: Number(BigInt.asUintN(8, BigInt(aOverflow8 + bOverflow8))),
                    name: " with overflow",
                },
            ],
            resType: 8,
        },
        {
            function: "add(euint16,euint16)",
            cases: [
                {a: 1, b: 2, expectedResult: 3, name: ""},
                {
                    a: aOverflow16,
                    b: bOverflow16,
                    expectedResult: Number(BigInt.asUintN(16, BigInt(aOverflow16 + bOverflow16))),
                    name: " with overflow",
                },
            ],
            resType: 16,
        },
        {
            function: "add(euint32,euint32)",
            cases: [
                {a: 1, b: 2, expectedResult: 3, name: ""},
                {
                    a: aOverflow32,
                    b: bOverflow32,
                    expectedResult: Number(BigInt.asUintN(32, BigInt(aOverflow32 + bOverflow32))),
                    name: " with overflow",
                },
            ],
            resType: 32,
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            it(`Test ${test.function}${testCase.name}`, async () => {
                const decryptedResult = await contract.add(test.function, BigInt(testCase.a), BigInt(testCase.b));
                expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
            });
        }
    }
});
describe('Test Reencrypt', () =>  {
    let contract;
    let fheContract;
    let contractAddress;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        const baseContract = await deployContract('ReencryptTest');
        contract = baseContract as ReencryptTestType;
        contractAddress = await baseContract.getAddress();
        fheContract = await getFheContract(contractAddress);

        expect(contract).toBeTruthy();
        expect(fheContract).toBeTruthy();

    });

    const testCases = ["reencrypt(euint8)", "reencrypt(euint16)", "reencrypt(euint32)"];

    for (const test of testCases) {
        it(`Test ${test}`, async () => {
            let plaintextInput = Math.floor(Math.random() * 1000) % 256;
            let encryptedOutput = await contract.reencrypt(test, plaintextInput, fheContract.publicKey);
            let decryptedOutput = fheContract.instance.decrypt(contractAddress, encryptedOutput);

            expect(decryptedOutput).toBe(plaintextInput);
        });
}
});
describe('Test Lte', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('LteTest') as LteTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 1, b: 2, expectedResult: 1, name: " a < b" },
        { a: 2, b: 1, expectedResult: 0, name: " a > b" },
        { a: 3, b: 3, expectedResult: 1, name: " a == b" },
    ];

    const testCases = [
        {
            function: "lte(euint8,euint8)",
            cases,
        },
        {
            function: "lte(euint16,euint16)",
            cases,
        },
        {
            function: "lte(euint32,euint32)",
            cases,
        }
    ];
    for (const test of testCases) {
        for (const testCase of test.cases) {
            it(`Test ${test.function}${testCase.name}`, async () => {
                const decryptedResult = await contract.lte(test.function, BigInt(testCase.a), BigInt(testCase.b));
                expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
            });
        }
    }
});
describe('Test Sub', () =>  {
    const aUnderflow = 1;
    const bUnderflow = 4;
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('SubTest') as SubTestType;
        expect(contract).toBeTruthy();
    });

    const testCases = [
        {
            function: "sub(euint8,euint8)",
            cases: [
                { a: 9, b: 4, expectedResult: 5, name: "" },
                {
                    a: aUnderflow,
                    b: bUnderflow,
                    expectedResult: Number(BigInt.asUintN(8, BigInt(aUnderflow - bUnderflow))),
                    name: " with underflow",
                },
            ],
        },
        {
            function: "sub(euint16,euint16)",
            cases: [
                { a: 9, b: 4, expectedResult: 5, name: "" },
                {
                    a: aUnderflow,
                    b: bUnderflow,
                    expectedResult: Number(BigInt.asUintN(16, BigInt(aUnderflow - bUnderflow))),
                    name: " with underflow",
                },
            ],
        },
        {
            function: "sub(euint32,euint32)",
            cases: [
                { a: 9, b: 4, expectedResult: 5, name: "" },
                {
                    a: aUnderflow,
                    b: bUnderflow,
                    expectedResult: Number(BigInt.asUintN(32, BigInt(aUnderflow - bUnderflow))),
                    name: " with underflow",
                },
            ],
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            it(`Test ${test.function}${testCase.name}`, async () => {
                const decryptedResult = await contract.sub(test.function, BigInt(testCase.a), BigInt(testCase.b));
                expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
            });
        }
    }
});
describe('Test Mul', () =>  {
    const overflow8 = 2 ** 8 / 2 + 1;
    const overflow16 = 2 ** 16 / 2 + 1;
    const overflow32 = 2 ** 32 / 2 + 1;
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('MulTest') as MulTestType;
        expect(contract).toBeTruthy();
    });

    const testCases = [
        {
            function: "mul(euint8,euint8)",
            cases: [
                { a: 2, b: 3, expectedResult: 6, name: "" },
                {
                    a: overflow8,
                    b: 2,
                    expectedResult: Number(BigInt.asUintN(8, BigInt(overflow8 * 2))),
                    name: " as overflow",
                },
            ],
        },
        {
            function: "mul(euint16,euint16)",
            cases: [
                { a: 2, b: 3, expectedResult: 6, name: "" },
                {
                    a: overflow16,
                    b: 2,
                    expectedResult: Number(BigInt.asUintN(16, BigInt(overflow16 * 2))),
                    name: " as overflow",
                },
            ],
        },
        {
            function: "mul(euint32,euint32)",
            cases: [
                { a: 2, b: 3, expectedResult: 6, name: "" },
                {
                    a: overflow32,
                    b: 2,
                    expectedResult: Number(BigInt.asUintN(32, BigInt(overflow32 * 2))),
                    name: " as overflow",
                },
            ],
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            it(`Test ${test.function}${testCase.name}`, async () => {
                const decryptedResult = await contract.mul(test.function, BigInt(testCase.a), BigInt(testCase.b));
                expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
            });
        }
    }
});
describe('Test Lt', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('LtTest') as LtTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 1, b: 2, expectedResult: 1, name: " a < b" },
        { a: 2, b: 1, expectedResult: 0, name: " a > b" },
        { a: 3, b: 3, expectedResult: 0, name: " a == b" },
    ];

    const testCases = [
        {
            function: "lt(euint8,euint8)",
            cases,
        },
        {
            function: "lt(euint16,euint16)",
            cases,
        },
        {
            function: "lt(euint32,euint32)",
            cases,
        }
    ];
    for (const test of testCases) {
        for (const testCase of test.cases) {
            it(`Test ${test.function}${testCase.name}`, async () => {
                const decryptedResult = await contract.lt(test.function, BigInt(testCase.a), BigInt(testCase.b));
                expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
            });
        }
    }
});



