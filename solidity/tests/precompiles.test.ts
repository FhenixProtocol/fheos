import { ethers } from 'hardhat';
import { BaseContract } from "ethers";
import {createFheInstance, fromHexString} from './utils';
import { AddTestType,
    LteTestType,
    SubTestType,
    MulTestType,
    LtTestType,
    SelectTestType,
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
    NotTestType,
    AsEboolTestType,
    AsEuint8TestType,
    AsEuint16TestType,
    AsEuint32TestType,
    SealoutputTestType
} from './abis';

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
            function: ["add(euint8,euint8)", "euint8.add(euint8)", "euint8 + euint8"],
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
            function: ["add(euint16,euint16)", "euint16.add(euint16)", "euint16 + euint16"],
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
            function: ["add(euint32,euint32)", "euint32.add(euint32)", "euint32 + euint32"],
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
        for (const testFunc of test.function) {
            for (const testCase of test.cases) {
                it(`Test ${testFunc}${testCase.name}`, async () => {
                    const decryptedResult = await contract.add(testFunc, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }

    it(`Custom error test`, async () => {
        try {
            await contract.add("no such test", 1, 2)
            fail();
        } catch (error) {
            const revertData = error.data
            const decodedError = contract.interface.parseError(revertData);
            expect(decodedError.name).toBe("TestNotFound");
            expect(decodedError.args[0]).toBe("no such test");
        }
    });
});

describe('Test SealOutput', () =>  {
    let contract;
    let fheContract;
    let contractAddress;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        const baseContract = await deployContract('SealoutputTest');
        contract = baseContract as SealoutputTestType;
        contractAddress = await baseContract.getAddress();
        fheContract = await getFheContract(contractAddress);

        expect(contract).toBeTruthy();
        expect(fheContract).toBeTruthy();

    });

    const testCases = ["sealoutput(euint8)", "sealoutput(euint16)", "sealoutput(euint32)", "seal(euint8)"];

    for (const test of testCases) {
        it(`Test ${test}`, async () => {
            let plaintextInput = Math.floor(Math.random() * 1000) % 256;
            let encryptedOutput = await contract.sealoutput(test, plaintextInput, fromHexString(fheContract.permit.sealingKey.publicKey));
            let decryptedOutput = fheContract.instance.unseal(contractAddress, encryptedOutput);
            expect(decryptedOutput).toBe(BigInt(plaintextInput));
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
            function: ["lte(euint8,euint8)", "euint8.lte(euint8)"],
            cases,
        },
        {
            function: ["lte(euint16,euint16)", "euint16.lte(euint16)"],
            cases,
        },
        {
            function: ["lte(euint32,euint32)", "euint32.lte(euint32)"],
            cases,
        }
    ];
    for (const test of testCases) {
        for (const testFunc of test.function) {
            for (const testCase of test.cases) {
                it(`Test ${testFunc}${testCase.name}`, async () => {
                    const decryptedResult = await contract.lte(testFunc, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
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
            function: ["sub(euint8,euint8)", "euint8.sub(euint8)", "euint8 - euint8"],
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
            function: ["sub(euint16,euint16)", "euint16.sub(euint16)", "euint16 - euint16"],
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
            function: ["sub(euint32,euint32)", "euint32.sub(euint32)", "euint32 - euint32"],
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
        for (const testFunc of test.function) {
            for (const testCase of test.cases) {
                it(`Test ${testFunc}${testCase.name}`, async () => {
                    const decryptedResult = await contract.sub(testFunc, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
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
            function: ["mul(euint8,euint8)", "euint8.mul(euint8)", "euint8 * euint8"],
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
            function: ["mul(euint16,euint16)", "euint16.mul(euint16)", "euint16 * euint16"],
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
            function: ["mul(euint32,euint32)", "euint32.mul(euint32)", "euint32 * euint32"],
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
        for (const testFunc of test.function) {
            for (const testCase of test.cases) {
                it(`Test ${testFunc}${testCase.name}`, async () => {
                    const decryptedResult = await contract.mul(testFunc, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
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
            function: ["lt(euint8,euint8)", "euint8.lt(euint8)"],
            cases,
        },
        {
            function: ["lt(euint16,euint16)", "euint16.lt(euint16)"],
            cases,
        },
        {
            function: ["lt(euint32,euint32)", "euint32.lt(euint32)"],
            cases,
        }
    ];
    for (const test of testCases) {
        for (const testFunc of test.function) {
            for (const testCase of test.cases) {
                it(`Test ${testFunc}${testCase.name}`, async () => {
                    const decryptedResult = await contract.lt(testFunc, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Select', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('SelectTest') as SelectTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { control: true, a: 2, b: 3, expectedResult: 2, name: " true" },
        { control: false, a: 2, b: 3, expectedResult: 3, name: " false" },
    ];

    const testCases = [
        {
            function: "select: euint8",
            cases,
        },
        {
            function: "select: euint16",
            cases,
        },
        {
            function: "select: euint32",
            cases,
        },
        {
            function: "select: ebool",
            cases: [
                { control: true, a: 0, b: 1, expectedResult: 0, name: "true" },
                { control: false, a: 0, b: 1, expectedResult: 1, name: "false" },
            ],
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            it(`Test ${test.function}${testCase.name}`, async () => {
                const decryptedResult = await contract.select(test.function, testCase.control, BigInt(testCase.a), BigInt(testCase.b));
                expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
            });
        }
    }
});

describe('Test Req', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('ReqTest') as ReqTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 0, shouldCrash: true, name: " with crash" },
        { a: 1, shouldCrash: false, name: " no crash" },
    ];

    const testCases = [
        {
            function: "req(euint8)",
            cases,
        },
        {
            function: "req(euint16)",
            cases,
        },
        {
            function: "req(euint32)",
            cases,
        },
        {
            function: "req(ebool)",
            cases,
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            it(`Test ${test.function}${testCase.name}`, async () => {
                let hadEvaluationFailure = false;
                let err = "";
                try {
                    await contract.req(test.function, BigInt(testCase.a));
                } catch (e) {
                    console.log(e);
                    hadEvaluationFailure = true;
                    err = `${e}`;
                }
                expect(hadEvaluationFailure).toBe(testCase.shouldCrash);
                if (hadEvaluationFailure) {
                    expect(err.includes("execution reverted")).toBe(true);
                }
            });
        }
    }
});

describe('Test Div', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('DivTest') as DivTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 4, b: 2, expectedResult: 2, name: "" },
        { a: 4, b: 0, expectedResult: 2 ** 256, name: " Div by 0" },
    ];

    const testCases = [
        {
            function: ["div(euint8,euint8)", "euint8.div(euint8)", "euint8 / euint8"],
            cases: [
                { a: 4, b: 2, expectedResult: 2, name: "" },
                { a: 4, b: 3, expectedResult: 1, name: " with reminder" },
                { a: 4, b: 0, expectedResult: 2 ** 8 - 1, name: " div by 0" },
            ],
        },
        {
            function: ["div(euint16,euint16)", "euint16.div(euint16)", "euint16 / euint16"],
            cases: [
                { a: 4, b: 2, expectedResult: 2, name: "" },
                { a: 4, b: 3, expectedResult: 1, name: " with reminder" },
                { a: 4, b: 0, expectedResult: 2 ** 16 - 1, name: " div by 0" },
            ],
        },
        {
            function: ["div(euint32,euint32)", "euint32.div(euint32)", "euint32 / euint32"],
            cases: [
                { a: 4, b: 2, expectedResult: 2, name: "" },
                { a: 4, b: 3, expectedResult: 1, name: " with reminder" },
                { a: 4, b: 0, expectedResult: 2 ** 32 - 1, name: " div by 0" },
            ],
        }
    ];
    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const functionSignature of test.function) {
                it(`Test ${functionSignature}${testCase.name}`, async () => {
                    const decryptedResult = await contract.div(functionSignature, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Gt', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('GtTest') as GtTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 1, b: 2, expectedResult: 0, name: " a < b" },
        { a: 2, b: 1, expectedResult: 1, name: " a > b" },
        { a: 3, b: 3, expectedResult: 0, name: " a == b" },
    ];

    const testCases = [
        {
            function: ["gt(euint8,euint8)", "euint8.gt(euint8)"],
            cases,
        },
        {
            function: ["gt(euint16,euint16)", "euint16.gt(euint16)"],
            cases,
        },
        {
            function: ["gt(euint32,euint32)", "euint32.gt(euint32)"],
            cases,
        }
    ];
    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const functionSignature of test.function) {
                it(`Test ${functionSignature}${testCase.name}`, async () => {
                    const decryptedResult = await contract.gt(functionSignature, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Gte', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('GteTest') as GteTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 1, b: 2, expectedResult: 0, name: " a < b" },
        { a: 2, b: 1, expectedResult: 1, name: " a > b" },
        { a: 3, b: 3, expectedResult: 1, name: " a == b" },
    ];

    const testCases = [
        {
            function: ["gte(euint8,euint8)", "euint8.gte(euint8)"],
            cases,
        },
        {
            function: ["gte(euint16,euint16)", "euint16.gte(euint16)"],
            cases,
        },
        {
            function: ["gte(euint32,euint32)", "euint32.gte(euint32)"],
            cases,
        }
    ];
    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const functionSignature of test.function) {
                it(`Test ${functionSignature}${testCase.name}`, async () => {
                    const decryptedResult = await contract.gte(functionSignature, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Rem', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('RemTest') as RemTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 4, b: 3, expectedResult: 1, name: "" },
        { a: 4, b: 2, expectedResult: 0, name: " no reminder" },
        { a: 4, b: 0, expectedResult: 4, name: " div by 0" },
    ];

    const testCases = [
        {
            function: ["rem(euint8,euint8)", "euint8.rem(euint8)", "euint8 % euint8"],
            cases,
        },
        {
            function: ["rem(euint16,euint16)", "euint16.rem(euint16)", "euint16 % euint16"],
            cases,
        },
        {
            function: ["rem(euint32,euint32)", "euint32.rem(euint32)", "euint32 % euint32"],
            cases,
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const functionSignature of test.function) {
                it(`Test ${functionSignature}${testCase.name}`, async () => {
                    const decryptedResult = await contract.rem(functionSignature, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test And', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('AndTest') as AndTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 9, b: 15, expectedResult: 9 & 15, name: "" },
        { a: 7, b: 0, expectedResult: 0, name: " a & 0" },
        { a: 0, b: 5, expectedResult: 0, name: " 0 & b" },
    ];

    const testCases = [
        {
            function: ["and(euint8,euint8)", "euint8.and(euint8)", "euint8 & euint8"],
            cases,
        },
        {
            function: ["and(euint16,euint16)", "euint16.and(euint16)", "euint16 & euint16"],
            cases,
        },
        {
            function: ["and(euint32,euint32)", "euint32.and(euint32)", "euint32 & euint32"],
            cases,
        },
        {
            function: ["and(ebool,ebool)", "ebool.and(ebool)", "ebool & ebool"],
            cases: [
                { a: 9, b: 15, expectedResult: 1, name: "" },
                { a: 7, b: 0, expectedResult: 0, name: " a & 0" },
                { a: 0, b: 5, expectedResult: 0, name: " 0 & b" },
            ],
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const functionSignature of test.function) {
                it(`Test ${functionSignature}${testCase.name}`, async () => {
                    const decryptedResult = await contract.and(functionSignature, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Or', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('OrTest') as OrTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 9, b: 15, expectedResult: 9 | 15, name: "" },
        { a: 7, b: 0, expectedResult: 7, name: " a | 0" },
        { a: 0, b: 5, expectedResult: 5, name: " 0 | b" },
    ];

    const testCases = [
        {
            function: ["or(euint8,euint8)", "euint8.or(euint8)", "euint8 | euint8"],
            cases,
        },
        {
            function: ["or(euint16,euint16)", "euint16.or(euint16)", "euint16 | euint16"],
            cases,
        },
        {
            function: ["or(euint32,euint32)", "euint32.or(euint32)", "euint32 | euint32"],
            cases,
        },
        {
            function: ["or(ebool,ebool)", "ebool.or(ebool)", "ebool | ebool"],
            cases: [
                { a: 9, b: 15, expectedResult: 1, name: "" },
                { a: 7, b: 0, expectedResult: 1, name: " a | 0" },
                { a: 0, b: 5, expectedResult: 1, name: " 0 | b" },
                { a: 0, b: 0, expectedResult: 0, name: " 0 | 0" },
            ],
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const functionSignature of test.function) {
                it(`Test ${functionSignature}${testCase.name}`, async () => {
                    const decryptedResult = await contract.or(functionSignature, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Xor', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('XorTest') as XorTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 0b11110000, b: 0b10100101, expectedResult: 0b11110000 ^ 0b10100101, name: "" },
        { a: 7, b: 0, expectedResult: 7 ^ 0, name: " a ^ 0s" },
        { a: 7, b: 0b1111, expectedResult: 7 ^ 0b1111, name: " a ^ 1s" },
        { a: 0, b: 5, expectedResult: 0 ^ 5, name: " 0s ^ b" },
        { a: 0b1111, b: 5, expectedResult: 0b1111 ^ 5, name: " 1s ^ b" },
    ];

    const testCases = [
        {
            function: ["xor(euint8,euint8)", "euint8.xor(euint8)", "euint8 ^ euint8"],
            cases,
        },
        {
            function: ["xor(euint16,euint16)", "euint16.xor(euint16)", "euint16 ^ euint16"],
            cases,
        },
        {
            function: ["xor(euint32,euint32)", "euint32.xor(euint32)", "euint32 ^ euint32"],
            cases,
        },
        {
            function: ["xor(ebool,ebool)", "ebool.xor(ebool)", "ebool ^ ebool"],
            cases: [
                { a: 9, b: 15, expectedResult: 0, name: "" },
                { a: 7, b: 0, expectedResult: 1, name: " a ^ 0" },
                { a: 0, b: 5, expectedResult: 1, name: " 0 ^ b" },
                { a: 0, b: 0, expectedResult: 0, name: " 0 ^ 0" },
            ],
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const functionSignature of test.function) {
                it(`Test ${functionSignature}${testCase.name}`, async () => {
                    const decryptedResult = await contract.xor(functionSignature, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Eq', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('EqTest') as EqTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 1, b: 2, expectedResult: 0, name: " a < b" },
        { a: 2, b: 1, expectedResult: 0, name: " a > b" },
        { a: 3, b: 3, expectedResult: 1, name: " a == b" },
    ];

    const testCases = [
        {
            function: ["eq(euint8,euint8)", "euint8.eq(euint8)"],
            cases,
        },
        {
            function: ["eq(euint16,euint16)", "euint16.eq(euint16)"],
            cases,
        },
        {
            function: ["eq(euint32,euint32)", "euint32.eq(euint32)"],
            cases,
        },
        {
            function: ["eq(ebool,ebool)", "ebool.eq(ebool)"],
            cases: [
                { a: 1, b: 1, expectedResult: 1, name: " 1 == 1" },
                { a: 0, b: 0, expectedResult: 1, name: " 0 == 0" },
                { a: 0, b: 1, expectedResult: 0, name: " 0 != 1" },
                { a: 1, b: 0, expectedResult: 0, name: " 1 != 0" },
            ],
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const funcName of test.function) {
                it(`Test ${funcName}${testCase.name}`, async () => {
                    const decryptedResult = await contract.eq(funcName, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Ne', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('NeTest') as NeTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 1, b: 2, expectedResult: 1, name: " a < b" },
        { a: 2, b: 1, expectedResult: 1, name: " a > b" },
        { a: 3, b: 3, expectedResult: 0, name: " a == b" },
    ];

    const testCases = [
        {
            function: ["ne(euint8,euint8)", "euint8.ne(euint8)"],
            cases,
        },
        {
            function: ["ne(euint16,euint16)", "euint16.ne(euint16)"],
            cases,
        },
        {
            function: ["ne(euint32,euint32)", "euint32.ne(euint32)"],
            cases,
        },
        {
            function: ["ne(ebool,ebool)", "ebool.ne(ebool)"],
            cases: [
                { a: 1, b: 1, expectedResult: 0, name: " 1 == 1" },
                { a: 0, b: 0, expectedResult: 0, name: " 0 == 0" },
                { a: 0, b: 1, expectedResult: 1, name: " 0 != 1" },
                { a: 1, b: 0, expectedResult: 1, name: " 1 != 0" },
            ],
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const funcName of test.function) {
                it(`Test ${funcName}${testCase.name}`, async () => {
                    const decryptedResult = await contract.ne(funcName, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Min', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('MinTest') as MinTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 1, b: 2, expectedResult: 1, name: " a < b" },
        { a: 2, b: 1, expectedResult: 1, name: " a > b" },
        { a: 3, b: 3, expectedResult: 3, name: " a == b" },
    ];

    const testCases = [
        {
            function: ["min(euint8,euint8)", "euint8.min(euint8)"],
            cases,
        },
        {
            function: ["min(euint16,euint16)", "euint16.min(euint16)"],
            cases,
        },
        {
            function: ["min(euint32,euint32)", "euint32.min(euint32)"],
            cases,
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const funcName of test.function) {
                it(`Test ${funcName}${testCase.name}`, async () => {
                    const decryptedResult = await contract.min(funcName, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Max', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('MaxTest') as MaxTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 1, b: 2, expectedResult: 2, name: " a < b" },
        { a: 2, b: 1, expectedResult: 2, name: " a > b" },
        { a: 3, b: 3, expectedResult: 3, name: " a == b" },
    ];

    const testCases = [
        {
            function: ["max(euint8,euint8)", "euint8.max(euint8)"],
            cases,
        },
        {
            function: ["max(euint16,euint16)", "euint16.max(euint16)"],
            cases,
        },
        {
            function: ["max(euint32,euint32)", "euint32.max(euint32)"],
            cases,
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const funcName of test.function) {
                it(`Test ${funcName}${testCase.name}`, async () => {
                    const decryptedResult = await contract.max(funcName, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Shl', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('ShlTest') as ShlTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 0b10101010, b: 1, expectedResult: 0b101010100, name: " <<1" },
        { a: 0b10101010, b: 2, expectedResult: 0b1010101000, name: " <<2" },
        { a: 0b10101010, b: 3, expectedResult: 0b10101010000, name: " <<3" },
        { a: 0b10101010, b: 4, expectedResult: 0b101010100000, name: " <<4" },
        { a: 0b10101010, b: 5, expectedResult: 0b1010101000000, name: " <<5" },
    ];

    const testCases = [
        {
            function: ["shl(euint8,euint8)", "euint8.shl(euint8)"],
            cases : [
                { a: 0b10101010, b: 1, expectedResult: 0b01010100, name: " <<1" },
                { a: 0b10101010, b: 2, expectedResult: 0b10101000, name: " <<2" },
                { a: 0b10101010, b: 3, expectedResult: 0b01010000, name: " <<3" },
                { a: 0b10101010, b: 4, expectedResult: 0b10100000, name: " <<4" },
                { a: 0b10101010, b: 5, expectedResult: 0b01000000, name: " <<5" },
            ],
        },
        {
            function: ["shl(euint16,euint16)", "euint16.shl(euint16)"],
            cases,
        },
        {
            function: ["shl(euint32,euint32)", "euint32.shl(euint32)"],
            cases,
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const funcName of test.function) {
                it(`Test ${funcName}${testCase.name}`, async () => {
                    const decryptedResult = await contract.shl(funcName, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Shr', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('ShrTest') as ShrTestType;
        expect(contract).toBeTruthy();
    });


    const cases = [
        { a: 0b10101010, b: 1, expectedResult: 0b01010101, name: ">>1" },
        { a: 0b10101010, b: 2, expectedResult: 0b00101010, name: ">>2" },
        { a: 0b10101010, b: 3, expectedResult: 0b00010101, name: ">>3" },
        { a: 0b10101010, b: 4, expectedResult: 0b00001010, name: ">>4" },
        { a: 0b10101010, b: 5, expectedResult: 0b00000101, name: ">>5" },
    ];

    const testCases = [
        {
            function: ["shr(euint8,euint8)", "euint8.shr(euint8)"],
            cases,
        },
        {
            function: ["shr(euint16,euint16)", "euint16.shr(euint16)"],
            cases,
        },
        {
            function: ["shr(euint32,euint32)", "euint32.shr(euint32)"],
            cases,
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            for (const funcName of test.function) {
                it(`Test ${funcName}${testCase.name}`, async () => {
                    const decryptedResult = await contract.shr(funcName, BigInt(testCase.a), BigInt(testCase.b));
                    expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
                });
            }
        }
    }
});

describe('Test Not', () =>  {
    let contract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        contract = await deployContract('NotTest') as NotTestType;
        expect(contract).toBeTruthy();
    });

    const cases = [
        { a: 9, b: 15, expectedResult: 9 | 15, name: "" },
        { a: 7, b: 0, expectedResult: 7, name: " a | 0" },
        { a: 0, b: 5, expectedResult: 5, name: " 0 | b" },
    ];

    const testCases = [
        {
            function: "not(euint8)",
            cases: [{ value: 0b11110000, expectedResult: 0b00001111, name: "" }],
        },
        {
            function: "not(euint16)",
            cases: [{ value: 0b1111111100000000, expectedResult: 0b0000000011111111, name: "" }],
        },
        {
            function: "not(euint32)",
            cases: [{ value: 0b11111111111111110000000000000000, expectedResult: 0b00000000000000001111111111111111, name: "" }],
        },
        {
            function: "not(ebool)",
            cases: [{ value: 1, expectedResult: 0, name: " !true" }, { value: 0, expectedResult: 1, name: " !false" }],
        },
    ];

    for (const test of testCases) {
        for (const testCase of test.cases) {
            it(`Test ${test.function}${testCase.name}`, async () => {
                const decryptedResult = await contract.not(test.function, BigInt(testCase.value));
                expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
            });
        }
    }
});

describe('Test AsEbool', () =>  {
    let contract;
    let fheContract;

    const funcTypes = ["regular", "bound"];
    const cases = [{input: BigInt(0), output: false}, {input: BigInt(5), output: true}]
    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        const baseContract = await deployContract('AsEboolTest');
        contract = baseContract  as AsEboolTestType;

        const contractAddress = await baseContract.getAddress();
        fheContract = await getFheContract(contractAddress);

        expect(contract).toBeTruthy();
        expect(fheContract).toBeTruthy();
    });

    for (const funcType of funcTypes) {
        it(`From euint8 - ${funcType}`, async () => {
            for (const testCase of cases) {
                let decryptedResult = await contract.castFromEuint8ToEbool(testCase.input, funcType);
                expect(decryptedResult).toBe(testCase.output);
            }
        });
    }

    for (const funcType of funcTypes) {
        it(`From euint16 - ${funcType}`, async () => {
            for (const testCase of cases) {
                let decryptedResult = await contract.castFromEuint16ToEbool(testCase.input, funcType);
                expect(decryptedResult).toBe(testCase.output);
            }
        });
    }

    for (const funcType of funcTypes) {
        it(`From euint32 - ${funcType}`, async () => {
            for (const testCase of cases) {
                let decryptedResult = await contract.castFromEuint32ToEbool(testCase.input, funcType);
                expect(decryptedResult).toBe(testCase.output);
            }
        });
    }

    it(`From plaintext`, async () => {
        for (const testCase of cases) {
            let decryptedResult = await contract.castFromPlaintextToEbool(testCase.input);
            expect(decryptedResult).toBe(testCase.output);
        }
    });

    it(`From pre encrypted`, async () => {
        for (const testCase of cases) {
            // skip for 0 as currently encrypting 0 is not supported
            if (testCase.input === BigInt(0)) {
                continue;
            }

            const encInput = await fheContract.instance.encrypt_uint8(Number(testCase.input));
            let decryptedResult = await contract.castFromPreEncryptedToEbool(encInput.data);
            expect(decryptedResult).toBe(testCase.output);
        }
    });
});

describe('Test AsEuint8', () =>  {
    let contract;
    let fheContract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        const baseContract = await deployContract('AsEuint8Test');
        contract = baseContract  as AsEuint8TestType;

        const contractAddress = await baseContract.getAddress();
        fheContract = await getFheContract(contractAddress);

        expect(contract).toBeTruthy();
        expect(fheContract).toBeTruthy();
    });

    const funcTypes = ["regular", "bound"];
    const value = BigInt(1);
    for (const funcType of funcTypes) {
        it(`From ebool - ${funcType}`, async () => {
            let decryptedResult = await contract.castFromEboolToEuint8(value, funcType);
            expect(decryptedResult).toBe(value);
        });
    }

    for (const funcType of funcTypes) {
        it(`From euint16 - ${funcType}`, async () => {
            let decryptedResult = await contract.castFromEuint16ToEuint8(value, funcType);
            expect(decryptedResult).toBe(value);
        });
    }

    for (const funcType of funcTypes) {
        it(`From euint32 - ${funcType}`, async () => {
            let decryptedResult = await contract.castFromEuint32ToEuint8(value, funcType);
            expect(decryptedResult).toBe(value);
        });
    }

    it(`From plaintext`, async () => {
        let decryptedResult = await contract.castFromPlaintextToEuint8(value);
        expect(decryptedResult).toBe(value);
    });

    it(`From pre encrypted`, async () => {
        const encInput = await fheContract.instance.encrypt_uint8(Number(value));
        let decryptedResult = await contract.castFromPreEncryptedToEuint8(encInput.data);
        expect(decryptedResult).toBe(value);
    });
});

describe('Test AsEuint16', () =>  {
    let contract;
    let fheContract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        const baseContract = await deployContract('AsEuint16Test');
        contract = baseContract  as AsEuint16TestType;

        const contractAddress = await baseContract.getAddress();
        fheContract = await getFheContract(contractAddress);

        expect(contract).toBeTruthy();
        expect(fheContract).toBeTruthy();
    });

    const value = BigInt(1);
    const funcTypes = ["regular", "bound"];

    for (const funcType of funcTypes) {
        it(`From ebool - ${funcType}`, async () => {
            let decryptedResult = await contract.castFromEboolToEuint16(value, funcType);
            expect(decryptedResult).toBe(value);
        });
    }

    for (const funcType of funcTypes) {
        it(`From euint8 - ${funcType}`, async () => {
            let decryptedResult = await contract.castFromEuint8ToEuint16(value, funcType);
            expect(decryptedResult).toBe(value);
        });
    }

    for (const funcType of funcTypes) {
        it(`From euint32 - ${funcType}`, async () => {
            let decryptedResult = await contract.castFromEuint32ToEuint16(value, funcType);
            expect(decryptedResult).toBe(value);
        });
    }

    it(`From plaintext`, async () => {
        let decryptedResult = await contract.castFromPlaintextToEuint16(value);
        expect(decryptedResult).toBe(value);
    });

    it(`From pre encrypted`, async () => {
        const encInput = await fheContract.instance.encrypt_uint16(Number(value));
        let decryptedResult = await contract.castFromPreEncryptedToEuint16(encInput.data);
        expect(decryptedResult).toBe(value);
    });
});

describe('Test AsEuint32', () =>  {
    let contract;
    let fheContract;

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        const baseContract = await deployContract('AsEuint32Test');
        contract = baseContract  as AsEuint32TestType;

        const contractAddress = await baseContract.getAddress();
        fheContract = await getFheContract(contractAddress);

        expect(contract).toBeTruthy();
        expect(fheContract).toBeTruthy();
    });

    const value = BigInt(1);
    const funcTypes = ["regular", "bound"];

    for (const funcType of funcTypes) {
        it(`From ebool - ${funcType}`, async () => {
            let decryptedResult = await contract.castFromEboolToEuint32(value, funcType);
            expect(decryptedResult).toBe(value);
        });
    }

    for (const funcType of funcTypes) {
        it(`From euint8 - ${funcType}`, async () => {
            let decryptedResult = await contract.castFromEuint8ToEuint32(value, funcType);
            expect(decryptedResult).toBe(value);
        });
    }

    for (const funcType of funcTypes) {
        it(`From euint16 - ${funcType}`, async () => {
            let decryptedResult = await contract.castFromEuint16ToEuint32(value, funcType);
            expect(decryptedResult).toBe(value);
        });
    }

    it(`From plaintext`, async () => {
        let decryptedResult = await contract.castFromPlaintextToEuint32(value);
        expect(decryptedResult).toBe(value);
    });

    it(`From pre encrypted`, async () => {
        const encInput = await fheContract.instance.encrypt_uint32(Number(value));
        let decryptedResult = await contract.castFromPreEncryptedToEuint32(encInput.data);
        expect(decryptedResult).toBe(value);
    });
});