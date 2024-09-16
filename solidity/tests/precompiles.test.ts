import { createFheInstance, deployContract, fromHexString } from "./utils";
import {
  AddTestType,
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
  RandomTestType,
  AsEboolTestType,
  AsEuint8TestType,
  AsEuint16TestType,
  AsEuint32TestType,
  SealoutputTestType,
  DecryptTestType,
  AsEuint64TestType,
  AsEuint128TestType,
  AsEuint256TestType,
  AsEaddressTestType,
  RolTestType,
  RorTestType
} from "./abis";

const getFheContract = async (contractAddress: string) => {
  const fheContract = await createFheInstance(contractAddress);
  return fheContract;
};

describe("Test Add", () => {
  const getMaxValue = (bits: number) => {
    const val = BigInt.asUintN(bits, BigInt(1));
    return (val << BigInt(bits)) - BigInt(2);
  };

  const cases = [
    { overflow: false, name: "" },
    {
      overflow: true,
      name: " with overflow",
    },
  ];

  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("AddTest")) as AddTestType;
    expect(contract).toBeTruthy();
  });

  const testCases = [
    {
      function: ["add(euint8,euint8)", "euint8.add(euint8)", "euint8 + euint8"],
      cases,
      resType: 8,
    },
    {
      function: [
        "add(euint16,euint16)",
        "euint16.add(euint16)",
        "euint16 + euint16",
      ],
      cases,
      resType: 16,
    },
    {
      function: [
        "add(euint32,euint32)",
        "euint32.add(euint32)",
        "euint32 + euint32",
      ],
      cases,
      resType: 32,
    },
    {
      function: [
        "add(euint64,euint64)",
        "euint64.add(euint64)",
        "euint64 + euint64",
      ],
      cases,
      resType: 64,
    },
    {
      function: [
        "add(euint128,euint128)",
        "euint128.add(euint128)",
        "euint128 + euint128",
      ],
      cases,
      resType: 128,
    },
  ];

  for (const test of testCases) {
    for (const testFunc of test.function) {
      for (const testCase of test.cases) {
        it(`Test ${testFunc}${testCase.name}`, async () => {
          let a = BigInt.asUintN(test.resType, BigInt(2));
          if (testCase.overflow) {
            a = getMaxValue(test.resType);
          }

          const b = a - BigInt(1);
          const decryptedResult = await contract.add(testFunc, a, b);
          expect(decryptedResult).toBe(BigInt.asUintN(test.resType, a + b));
        });
      }
    }
  }

  it(`Custom error test`, async () => {
    try {
      await contract.add("no such test", 1, 2);
      fail();
    } catch (error) {
      const revertData = error.data;
      const decodedError = contract.interface.parseError(revertData);
      expect(decodedError.name).toBe("TestNotFound");
      expect(decodedError.args[0]).toBe("no such test");
    }
  });
});

describe("Test SealOutput", () => {
  let contract;
  let fheContract;
  let contractAddress;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    const baseContract = await deployContract("SealoutputTest");
    contract = baseContract as SealoutputTestType;
    contractAddress = await baseContract.getAddress();
    fheContract = await getFheContract(contractAddress);

    expect(contract).toBeTruthy();
    expect(fheContract).toBeTruthy();
  });

  const testCases = [
    "sealoutput(euint8)",
    "sealoutput(euint16)",
    "sealoutput(euint32)",
    "sealoutput(euint64)",
    "sealoutput(euint128)",
    "sealoutput(euint256)",
    "sealoutput(ebool)",
    "seal(euint8)",
  ];

  for (const test of testCases) {
    it(`Test ${test}`, async () => {
      let plaintextInput = Math.floor(Math.random() * 1000) % 256;
      let encryptedOutput = await contract.sealoutput(
        test,
        plaintextInput,
        fromHexString(fheContract.permit.sealingKey.publicKey)
      );
      let decryptedOutput = fheContract.instance.unseal(
        contractAddress,
        encryptedOutput
      );
      if (test.includes("ebool")) {
        expect(decryptedOutput).toBe(BigInt(Math.min(1, plaintextInput)));
      } else {
        expect(decryptedOutput).toBe(BigInt(plaintextInput));
      }
    });
  }
});

describe("Test Lte", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("LteTest")) as LteTestType;
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
    },
    {
      function: ["lte(euint64,euint64)", "euint64.lte(euint64)"],
      cases,
    },
    {
      function: ["lte(euint128,euint128)", "euint128.lte(euint128)"],
      cases,
    },
  ];
  for (const test of testCases) {
    for (const testFunc of test.function) {
      for (const testCase of test.cases) {
        it(`Test ${testFunc}${testCase.name}`, async () => {
          const decryptedResult = await contract.lte(
            testFunc,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Sub", () => {
  const aUnderflow = 1;
  const bUnderflow = 4;
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("SubTest")) as SubTestType;
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
          expectedResult: Number(
            BigInt.asUintN(8, BigInt(aUnderflow - bUnderflow))
          ),
          name: " with underflow",
        },
      ],
    },
    {
      function: [
        "sub(euint16,euint16)",
        "euint16.sub(euint16)",
        "euint16 - euint16",
      ],
      cases: [
        { a: 9, b: 4, expectedResult: 5, name: "" },
        {
          a: aUnderflow,
          b: bUnderflow,
          expectedResult: Number(
            BigInt.asUintN(16, BigInt(aUnderflow - bUnderflow))
          ),
          name: " with underflow",
        },
      ],
    },
    {
      function: [
        "sub(euint32,euint32)",
        "euint32.sub(euint32)",
        "euint32 - euint32",
      ],
      cases: [
        { a: 9, b: 4, expectedResult: 5, name: "" },
        {
          a: aUnderflow,
          b: bUnderflow,
          expectedResult: Number(
            BigInt.asUintN(32, BigInt(aUnderflow - bUnderflow))
          ),
          name: " with underflow",
        },
      ],
    },
    {
      function: [
        "sub(euint64,euint64)",
        "euint64.sub(euint64)",
        "euint64 - euint64",
      ],
      cases: [
        { a: 9, b: 4, expectedResult: 5, name: "" },
        {
          a: aUnderflow,
          b: bUnderflow,
          expectedResult: BigInt.asUintN(64, BigInt(aUnderflow - bUnderflow)),
          name: " with underflow",
        },
      ],
    },
    {
      function: [
        "sub(euint128,euint128)",
        "euint128.sub(euint128)",
        "euint128 - euint128",
      ],
      cases: [
        { a: 9, b: 4, expectedResult: 5, name: "" },
        {
          a: aUnderflow,
          b: bUnderflow,
          expectedResult: BigInt.asUintN(128, BigInt(aUnderflow - bUnderflow)),
          name: " with underflow",
        },
      ],
    },
  ];

  for (const test of testCases) {
    for (const testFunc of test.function) {
      for (const testCase of test.cases) {
        it(`Test ${testFunc}${testCase.name}`, async () => {
          const decryptedResult = await contract.sub(
            testFunc,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(BigInt(decryptedResult)).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Mul", () => {
  const overflow8 = 2 ** 8 / 2 + 1;
  const overflow16 = 2 ** 16 / 2 + 1;
  const overflow32 = 2 ** 32 / 2 + 1;
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("MulTest")) as MulTestType;
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
      function: [
        "mul(euint16,euint16)",
        "euint16.mul(euint16)",
        "euint16 * euint16",
      ],
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
      function: [
        "mul(euint32,euint32)",
        "euint32.mul(euint32)",
        "euint32 * euint32",
      ],
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
    {
      function: [
        "mul(euint64,euint64)",
        "euint64.mul(euint64)",
        "euint64 * euint64",
      ],
      cases: [
        { a: 2, b: 3, expectedResult: 6, name: "" },
        {
          a: 3000300300,
          b: 2000200200,
          expectedResult: BigInt(6001201260120060000n),
          name: " with large number",
        },
      ],
    },
  ];

  for (const test of testCases) {
    for (const testFunc of test.function) {
      for (const testCase of test.cases) {
        it(`Test ${testFunc}${testCase.name}`, async () => {
          const decryptedResult = await contract.mul(
            testFunc,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Lt", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("LtTest")) as LtTestType;
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
    },
    {
      function: ["lt(euint64,euint64)", "euint64.lt(euint64)"],
      cases,
    },
    {
      function: ["lt(euint128,euint128)", "euint128.lt(euint128)"],
      cases,
    },
  ];
  for (const test of testCases) {
    for (const testFunc of test.function) {
      for (const testCase of test.cases) {
        it(`Test ${testFunc}${testCase.name}`, async () => {
          const decryptedResult = await contract.lt(
            testFunc,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Select", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("SelectTest")) as SelectTestType;
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
      function: "select: euint64",
      cases,
    },
    {
      function: "select: euint128",
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
        const decryptedResult = await contract.select(
          test.function,
          testCase.control,
          BigInt(testCase.a),
          BigInt(testCase.b)
        );
        expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
      });
    }
  }
});

describe("Test Req", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("ReqTest")) as ReqTestType;
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
      function: "req(euint64)",
      cases,
    },
    {
      function: "req(euint128)",
      cases,
    },
    {
      function: "req(euint256)",
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
          let tx = await contract.req(test.function, BigInt(testCase.a));
          let result = await tx.wait();
        } catch (e) {
          hadEvaluationFailure = true;
          err = `${e}`;
          console.log(`err: ${err}`);
        }
        expect(hadEvaluationFailure).toBe(testCase.shouldCrash);
        if (hadEvaluationFailure) {
          expect(err.includes("execution reverted")).toBe(true);
          if (!testCase.shouldCrash) {
            console.log(`crashed in req even though it shouldn't have: ${err}`);
          }
        }
      });
    }
  }
});

describe("Test Decrypt", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("DecryptTest")) as DecryptTestType;
    expect(contract).toBeTruthy();
  });

  const testCases = [
    {
      function: "decrypt(euint8)",
      shouldPass: true,
    },
    {
      function: "decrypt(euint16)",
      shouldPass: true,
    },
    {
      function: "decrypt(euint32)",
      shouldPass: true,
    },
    {
      function: "decrypt(euint64)",
      shouldPass: true,
    },
    {
      function: "decrypt(euint128)",
      shouldPass: true,
    },
    {
      function: "decrypt(euint256)",
      shouldPass: true,
    },
    {
      function: "decrypt(ebool)",
      shouldPass: true,
    },
    {
      function: "decrypt(euint8) fail",
      shouldPass: false,
    },
    {
      function: "decrypt(euint16) fail",
      shouldPass: false,
    },
    {
      function: "decrypt(euint32) fail",
      shouldPass: false,
    },
    {
      function: "decrypt(euint64) fail",
      shouldPass: false,
    },
    {
      function: "decrypt(euint128) fail",
      shouldPass: false,
    },
    {
      function: "decrypt(euint256) fail",
      shouldPass: false,
    },
    {
      function: "decrypt(ebool) fail",
      shouldPass: false,
    },
  ];

  for (const test of testCases) {
    it(`Test ${test.function}`, async () => {
      let hadEvaluationFailure = false;
      let err = "";
      try {
        let tx = await contract.decrypt(test.function);
        let result = await tx.wait();
      } catch (e) {
        hadEvaluationFailure = true;
        err = `${e}`;
        console.log(`err: ${err}`);
      }
      expect(hadEvaluationFailure).toBe(!test.shouldPass);
      if (hadEvaluationFailure) {
        expect(err.includes("execution reverted")).toBe(true);
        if (test.shouldPass) {
          console.log(`crashed in req even though it shouldn't have: ${err}`);
        }
      }
    });
  }
});

describe("Test Div", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("DivTest")) as DivTestType;
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
      function: [
        "div(euint16,euint16)",
        "euint16.div(euint16)",
        "euint16 / euint16",
      ],
      cases: [
        { a: 4, b: 2, expectedResult: 2, name: "" },
        { a: 4, b: 3, expectedResult: 1, name: " with reminder" },
        { a: 4, b: 0, expectedResult: 2 ** 16 - 1, name: " div by 0" },
      ],
    },
    {
      function: [
        "div(euint32,euint32)",
        "euint32.div(euint32)",
        "euint32 / euint32",
      ],
      cases: [
        { a: 4, b: 2, expectedResult: 2, name: "" },
        { a: 4, b: 3, expectedResult: 1, name: " with reminder" },
        { a: 4, b: 0, expectedResult: 2 ** 32 - 1, name: " div by 0" },
      ],
    },
    // uint64 and uint128 and uint256 aren't permitted at this time
  ];
  for (const test of testCases) {
    for (const testCase of test.cases) {
      for (const functionSignature of test.function) {
        it(`Test ${functionSignature}${testCase.name}`, async () => {
          const decryptedResult = await contract.div(
            functionSignature,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Gt", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("GtTest")) as GtTestType;
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
    },
    {
      function: ["gt(euint64,euint64)", "euint64.gt(euint64)"],
      cases,
    },
    {
      function: ["gt(euint128,euint128)", "euint128.gt(euint128)"],
      cases,
    },
  ];
  for (const test of testCases) {
    for (const testCase of test.cases) {
      for (const functionSignature of test.function) {
        it(`Test ${functionSignature}${testCase.name}`, async () => {
          const decryptedResult = await contract.gt(
            functionSignature,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Gte", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("GteTest")) as GteTestType;
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
    },
    {
      function: ["gte(euint64,euint64)", "euint64.gte(euint64)"],
      cases,
    },
    {
      function: ["gte(euint128,euint128)", "euint128.gte(euint128)"],
      cases,
    },
  ];
  for (const test of testCases) {
    for (const testCase of test.cases) {
      for (const functionSignature of test.function) {
        it(`Test ${functionSignature}${testCase.name}`, async () => {
          const decryptedResult = await contract.gte(
            functionSignature,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Rem", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("RemTest")) as RemTestType;
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
      function: [
        "rem(euint16,euint16)",
        "euint16.rem(euint16)",
        "euint16 % euint16",
      ],
      cases,
    },
    {
      function: [
        "rem(euint32,euint32)",
        "euint32.rem(euint32)",
        "euint32 % euint32",
      ],
      cases,
    },
  ];

  for (const test of testCases) {
    for (const testCase of test.cases) {
      for (const functionSignature of test.function) {
        it(`Test ${functionSignature}${testCase.name}`, async () => {
          const decryptedResult = await contract.rem(
            functionSignature,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test And", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("AndTest")) as AndTestType;
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
      function: [
        "and(euint16,euint16)",
        "euint16.and(euint16)",
        "euint16 & euint16",
      ],
      cases,
    },
    {
      function: [
        "and(euint32,euint32)",
        "euint32.and(euint32)",
        "euint32 & euint32",
      ],
      cases,
    },
    {
      function: [
        "and(euint64,euint64)",
        "euint64.and(euint64)",
        "euint64 & euint64",
      ],
      cases,
    },
    {
      function: [
        "and(euint128,euint128)",
        "euint128.and(euint128)",
        "euint128 & euint128",
      ],
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
          const decryptedResult = await contract.and(
            functionSignature,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Or", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("OrTest")) as OrTestType;
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
      function: [
        "or(euint16,euint16)",
        "euint16.or(euint16)",
        "euint16 | euint16",
      ],
      cases,
    },
    {
      function: [
        "or(euint32,euint32)",
        "euint32.or(euint32)",
        "euint32 | euint32",
      ],
      cases,
    },
    {
      function: [
        "or(euint64,euint64)",
        "euint64.or(euint64)",
        "euint64 | euint64",
      ],
      cases,
    },
    {
      function: [
        "or(euint128,euint128)",
        "euint128.or(euint128)",
        "euint128 | euint128",
      ],
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
          const decryptedResult = await contract.or(
            functionSignature,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Xor", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("XorTest")) as XorTestType;
    expect(contract).toBeTruthy();
  });

  const cases = [
    {
      a: 0b11110000,
      b: 0b10100101,
      expectedResult: 0b11110000 ^ 0b10100101,
      name: "",
    },
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
      function: [
        "xor(euint16,euint16)",
        "euint16.xor(euint16)",
        "euint16 ^ euint16",
      ],
      cases,
    },
    {
      function: [
        "xor(euint32,euint32)",
        "euint32.xor(euint32)",
        "euint32 ^ euint32",
      ],
      cases,
    },
    {
      function: [
        "xor(euint64,euint64)",
        "euint64.xor(euint64)",
        "euint64 ^ euint64",
      ],
      cases,
    },
    {
      function: [
        "xor(euint128,euint128)",
        "euint128.xor(euint128)",
        "euint128 ^ euint128",
      ],
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
          const decryptedResult = await contract.xor(
            functionSignature,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Eq", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("EqTest")) as EqTestType;
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
      function: ["eq(euint64,euint64)", "euint64.eq(euint64)"],
      cases,
    },
    {
      function: ["eq(euint128,euint128)", "euint128.eq(euint128)"],
      cases,
    },
    {
      function: ["eq(euint256,euint256)", "euint256.eq(euint256)"],
      cases,
    },
    {
      function: ["eq(eaddress,eaddress)", "eaddress.eq(eaddress)"],
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
          const decryptedResult = await contract.eq(
            funcName,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Ne", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("NeTest")) as NeTestType;
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
      function: ["ne(euint64,euint64)", "euint64.ne(euint64)"],
      cases,
    },
    {
      function: ["ne(euint128,euint128)", "euint128.ne(euint128)"],
      cases,
    },
    {
      function: ["ne(euint256,euint256)", "euint256.ne(euint256)"],
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
          const decryptedResult = await contract.ne(
            funcName,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Min", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("MinTest")) as MinTestType;
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
    {
      function: ["min(euint64,euint64)", "euint64.min(euint64)"],
      cases,
    },
    {
      function: ["min(euint128,euint128)", "euint128.min(euint128)"],
      cases,
    },
  ];

  for (const test of testCases) {
    for (const testCase of test.cases) {
      for (const funcName of test.function) {
        it(`Test ${funcName}${testCase.name}`, async () => {
          const decryptedResult = await contract.min(
            funcName,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Max", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("MaxTest")) as MaxTestType;
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
    {
      function: ["max(euint64,euint64)", "euint64.max(euint64)"],
      cases,
    },
    {
      function: ["max(euint128,euint128)", "euint128.max(euint128)"],
      cases,
    },
  ];

  for (const test of testCases) {
    for (const testCase of test.cases) {
      for (const funcName of test.function) {
        it(`Test ${funcName}${testCase.name}`, async () => {
          const decryptedResult = await contract.max(
            funcName,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Shl", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("ShlTest")) as ShlTestType;
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
      cases: [
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
    {
      function: ["shl(euint64,euint64)", "euint64.shl(euint64)"],
      cases,
    },
    {
      function: ["shl(euint128,euint128)", "euint128.shl(euint128)"],
      cases,
    },
  ];

  for (const test of testCases) {
    for (const testCase of test.cases) {
      for (const funcName of test.function) {
        it(`Test ${funcName}${testCase.name}`, async () => {
          const decryptedResult = await contract.shl(
            funcName,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Shr", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("ShrTest")) as ShrTestType;
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
    {
      function: ["shr(euint64,euint64)", "euint64.shr(euint64)"],
      cases,
    },
    {
      function: ["shr(euint128,euint128)", "euint128.shr(euint128)"],
      cases,
    },
  ];

  for (const test of testCases) {
    for (const testCase of test.cases) {
      for (const funcName of test.function) {
        it(`Test ${funcName}${testCase.name}`, async () => {
          const decryptedResult = await contract.shr(
            funcName,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  }
});

describe("Test Rol", () => {
  let contract;

  beforeAll(async () => {
    contract = (await deployContract("RolTest")) as RolTestType;
  });

  const generateTestCases = (bitSize: number) => {
    const basePattern = BigInt(`0b${'11100000'.repeat(bitSize / 8)}`);
    const mask = (1n << BigInt(bitSize)) - 1n;

    return [1, 2, 3, 4, 5].map(rotateAmount => {
      const rotated = ((basePattern << BigInt(rotateAmount)) | (basePattern >> BigInt(bitSize - rotateAmount))) & mask;
      return {
        a: basePattern,
        b: rotateAmount,
        expectedResult: rotated,
        name: ` rol ${rotateAmount}`
      };
    });
  };

  const bitSizes = [8, 16, 32, 64, 128];

  const funcsToTest = bitSizes.map(size => [`rol(euint${size},euint${size})`, `euint${size}.rol(euint${size})`]);

  funcsToTest.forEach((test, index) => {
    const bitSize = bitSizes[index];
    const testCases = generateTestCases(bitSize);

    for (const testCase of testCases) {
      for (const funcName of test) {
        it(`Test ${funcName}${testCase.name}`, async () => {
          const decryptedResult = await contract.rol(
            funcName,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );

          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      }
    }
  });
});

describe("Test Ror", () => {
  let contract;

  beforeAll(async () => {
    contract = (await deployContract("RorTest")) as RorTestType;
  });

  const generateTestCases = (bitSize: number) => {
    const basePattern = BigInt(`0b${'11100000'.repeat(bitSize / 8)}`);
    const mask = (1n << BigInt(bitSize)) - 1n;

    return [1, 2, 3, 4, 5].map(rotateAmount => {
      const rotated = ((basePattern >> BigInt(rotateAmount)) | (basePattern << BigInt(bitSize - rotateAmount))) & mask;
      return {
        a: basePattern,
        b: rotateAmount,
        expectedResult: rotated,
        name: ` ror ${rotateAmount}`
      };
    });
  };

  const bitSizes = [8, 16, 32, 64, 128];

  const funcsToTest = bitSizes.map(size => [`ror(euint${size},euint${size})`, `euint${size}.ror(euint${size})`]);

  funcsToTest.forEach((test, index) => {
    const bitSize = bitSizes[index];
    const testCases = generateTestCases(bitSize);

    for (const testCase of testCases) {
      for (const funcName of test) {
        it(`Test ${funcName}${testCase.name}`, async () => {
          const decryptedResult = await contract.ror(
            funcName,
            BigInt(testCase.a),
            BigInt(testCase.b)
          );
          
          expect(decryptedResult).toBe(BigInt(testCase.expectedResult));
        });
      } 
    }
  });
});


describe("Test Not", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("NotTest")) as NotTestType;
    expect(contract).toBeTruthy();
  });

  const testCases = [
    {
      function: "not(ebool)",
      bits: 1,
    },
    {
      function: "not(euint8)",
      bits: 8,
    },
    {
      function: "not(euint16)",
      bits: 16,
    },
    {
      function: "not(euint32)",
      bits: 32,
    },
    {
      function: "not(euint64)",
      bits: 64,
    },
    {
      function: "not(euint128)",
      bits: 128,
    },
  ];

  for (const test of testCases) {
    for (const securityZone of [1]) {
      for (const input of [true, false]) {
        it(`Test ${test.function} !${input} - security zone ${securityZone}`, async () => {
          let val = BigInt(+input);
          let expectedResult = BigInt(+!input);
          if (test.bits !== 1) {
            val = BigInt.asUintN(test.bits, BigInt(+input));
            expectedResult = (1n << BigInt(test.bits)) - BigInt(1 + +input);
          }

          const decryptedResult = await contract.not(
            test.function,
            val,
            securityZone
          );

          expect(decryptedResult).toBe(expectedResult);
        });
      }
    }
  }
});

describe("Test Random", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("RandomTest")) as RandomTestType;
    expect(contract).toBeTruthy();
  });

  const testCases = [
    {
      function: "randomEuint8()",
      bits: 8,
    },
    {
      function: "randomEuint16()",
      bits: 16,
    },
    {
      function: "randomEuint32()",
      bits: 32,
    },
    {
      function: "randomEuint64()",
      bits: 64,
    },
    {
      function: "randomEuint128()",
      bits: 128,
    },
    {
      function: "randomEuint256()",
      bits: 256,
    },
  ];

  for (const test of testCases) {
    it(`Test ${test.function}`, async () => {
      const decryptedResult = await contract.random(test.function);

      expect(decryptedResult).toBeLessThan(2 ** test.bits);
    });
  }
});

describe("Test Random with seed", () => {
  let contract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    contract = (await deployContract("RandomSeedTest")) as RandomTestType;
    expect(contract).toBeTruthy();
  });

  it(`Test Random With Seed`, async () => {
    const firstResult = await contract.randomSeed(1337);
    const secondResult = await contract.randomSeed(87654);
    const thirdResult = await contract.randomSeed(1337);

    expect(firstResult).toBe(thirdResult);
    expect(firstResult).not.toBe(secondResult);
  });
});

describe("Test AsEbool", () => {
  let contract;
  let fheContract;

  const funcTypes = ["regular", "bound"];
  const cases = [
    { input: BigInt(0), output: false },
    { input: BigInt(5), output: true },
  ];
  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    const baseContract = await deployContract("AsEboolTest");
    contract = baseContract as AsEboolTestType;

    const contractAddress = await baseContract.getAddress();
    fheContract = await getFheContract(contractAddress);

    expect(contract).toBeTruthy();
    expect(fheContract).toBeTruthy();
  });

  for (const funcType of funcTypes) {
    it(`From euint8 - ${funcType}`, async () => {
      for (const testCase of cases) {
        let decryptedResult = await contract.castFromEuint8ToEbool(
          testCase.input,
          funcType
        );
        expect(decryptedResult).toBe(testCase.output);
      }
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint16 - ${funcType}`, async () => {
      for (const testCase of cases) {
        let decryptedResult = await contract.castFromEuint16ToEbool(
          testCase.input,
          funcType
        );
        expect(decryptedResult).toBe(testCase.output);
      }
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint32 - ${funcType}`, async () => {
      for (const testCase of cases) {
        let decryptedResult = await contract.castFromEuint32ToEbool(
          testCase.input,
          funcType
        );
        expect(decryptedResult).toBe(testCase.output);
      }
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint64 - ${funcType}`, async () => {
      for (const testCase of cases) {
        let decryptedResult = await contract.castFromEuint64ToEbool(
          testCase.input,
          funcType
        );
        expect(decryptedResult).toBe(testCase.output);
      }
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint128 - ${funcType}`, async () => {
      for (const testCase of cases) {
        let decryptedResult = await contract.castFromEuint128ToEbool(
          testCase.input,
          funcType
        );
        expect(decryptedResult).toBe(testCase.output);
      }
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint256 - ${funcType}`, async () => {
      for (const testCase of cases) {
        let decryptedResult = await contract.castFromEuint256ToEbool(
          testCase.input,
          funcType
        );
        expect(decryptedResult).toBe(testCase.output);
      }
    });
  }

  it(`From plaintext`, async () => {
    for (const testCase of cases) {
      let decryptedResult = await contract.castFromPlaintextToEbool(
        testCase.input
      );
      expect(decryptedResult).toBe(testCase.output);
    }
  });

  it(`From pre encrypted`, async () => {
    for (const testCase of cases) {
      // skip for 0 as currently encrypting 0 is not supported
      if (testCase.input === BigInt(0)) {
        continue;
      }

      const encInput = await fheContract.instance.encrypt_bool(
        !!Number(testCase.input)
      );
      let decryptedResult = await contract.castFromPreEncryptedToEbool(
        encInput
      );
      expect(decryptedResult).toBe(testCase.output);
    }
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    for (const testCase of cases) {
      // skip for 0 as currently encrypting 0 is not supported
      if (testCase.input === BigInt(0)) {
        continue;
      }

      const encInput = await fheContract.instance.encrypt_bool(
        !!Number(testCase.input),
        1 // non-default security zone
      );
      let decryptedResult = await contract.castFromPreEncryptedToEbool(
        encInput
      );
      expect(decryptedResult).toBe(testCase.output);
    }
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    for (const testCase of cases) {
      // skip for 0 as currently encrypting 0 is not supported
      if (testCase.input === BigInt(0)) {
        continue;
      }

      const encInput = await fheContract.instance.encrypt_bool(
        !!Number(testCase.input),
        1 // non-default security zone
      );
      let decryptedResult = await contract.castFromPreEncryptedToEbool(
        encInput
      );
      expect(decryptedResult).toBe(testCase.output);
    }
  });
});

describe("Test AsEuint8", () => {
  let contract;
  let fheContract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    const baseContract = await deployContract("AsEuint8Test");
    contract = baseContract as AsEuint8TestType;

    const contractAddress = await baseContract.getAddress();
    fheContract = await getFheContract(contractAddress);

    expect(contract).toBeTruthy();
    expect(fheContract).toBeTruthy();
  });

  const funcTypes = ["regular", "bound"];
  const value = BigInt(1);
  for (const funcType of funcTypes) {
    it(`From ebool - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEboolToEuint8(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint16 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint16ToEuint8(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint32 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint32ToEuint8(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint64 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint64ToEuint8(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint128 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint128ToEuint8(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint256 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint256ToEuint8(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  it(`From plaintext`, async () => {
    let decryptedResult = await contract.castFromPlaintextToEuint8(value);
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted`, async () => {
    const encInput = await fheContract.instance.encrypt_uint8(Number(value));
    let decryptedResult = await contract.castFromPreEncryptedToEuint8(encInput);
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_uint8(Number(value), 1);
    let decryptedResult = await contract.castFromPreEncryptedToEuint8(encInput);
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_uint8(Number(value), 1);
    let decryptedResult = await contract.castFromPreEncryptedToEuint8(encInput);
    expect(decryptedResult).toBe(value);
  });
});

describe("Test AsEuint16", () => {
  let contract;
  let fheContract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    const baseContract = await deployContract("AsEuint16Test");
    contract = baseContract as AsEuint16TestType;

    const contractAddress = await baseContract.getAddress();
    fheContract = await getFheContract(contractAddress);

    expect(contract).toBeTruthy();
    expect(fheContract).toBeTruthy();
  });

  const value = BigInt(1);
  const funcTypes = ["regular", "bound"];

  for (const funcType of funcTypes) {
    it(`From ebool - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEboolToEuint16(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint8 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint8ToEuint16(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint32 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint32ToEuint16(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint64 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint64ToEuint16(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint128 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint128ToEuint16(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint256 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint256ToEuint16(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  it(`From plaintext`, async () => {
    let decryptedResult = await contract.castFromPlaintextToEuint16(value);
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted`, async () => {
    const encInput = await fheContract.instance.encrypt_uint16(Number(value));
    let decryptedResult = await contract.castFromPreEncryptedToEuint16(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_uint16(
      Number(value),
      1
    );
    let decryptedResult = await contract.castFromPreEncryptedToEuint16(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_uint16(
      Number(value),
      1
    );
    let decryptedResult = await contract.castFromPreEncryptedToEuint16(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });
});

describe("Test AsEuint32", () => {
  let contract;
  let fheContract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    const baseContract = await deployContract("AsEuint32Test");
    contract = baseContract as AsEuint32TestType;

    const contractAddress = await baseContract.getAddress();
    fheContract = await getFheContract(contractAddress);

    expect(contract).toBeTruthy();
    expect(fheContract).toBeTruthy();
  });

  const value = BigInt(1);
  const funcTypes = ["regular", "bound"];

  for (const funcType of funcTypes) {
    it(`From ebool - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEboolToEuint32(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint8 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint8ToEuint32(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint16 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint16ToEuint32(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint64 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint64ToEuint32(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint128 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint128ToEuint32(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint256 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint256ToEuint32(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  it(`From plaintext`, async () => {
    let decryptedResult = await contract.castFromPlaintextToEuint32(value);
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted`, async () => {
    const encInput = await fheContract.instance.encrypt_uint32(Number(value));
    let decryptedResult = await contract.castFromPreEncryptedToEuint32(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_uint32(
      Number(value),
      1
    );
    let decryptedResult = await contract.castFromPreEncryptedToEuint32(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_uint32(
      Number(value),
      1
    );
    let decryptedResult = await contract.castFromPreEncryptedToEuint32(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });
});

describe("Test AsEuint64", () => {
  let contract;
  let fheContract;

  // We don't really need it as test but it is a test since it is async
  it(`Test Contract Deployment`, async () => {
    const baseContract = await deployContract("AsEuint64Test");
    contract = baseContract as AsEuint64TestType;

    const contractAddress = await baseContract.getAddress();
    fheContract = await getFheContract(contractAddress);

    expect(contract).toBeTruthy();
    expect(fheContract).toBeTruthy();
  });

  const value = BigInt(1);
  const funcTypes = ["regular", "bound"];

  for (const funcType of funcTypes) {
    it(`From ebool - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEboolToEuint64(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint8 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint8ToEuint64(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint16 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint16ToEuint64(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint32 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint32ToEuint64(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint128 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint128ToEuint64(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint256 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint256ToEuint64(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  it(`From plaintext`, async () => {
    let decryptedResult = await contract.castFromPlaintextToEuint64(value);
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted`, async () => {
    const encInput = await fheContract.instance.encrypt_uint64(value);
    let decryptedResult = await contract.castFromPreEncryptedToEuint64(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_uint64(value, 1);
    let decryptedResult = await contract.castFromPreEncryptedToEuint64(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_uint64(value, 1);
    let decryptedResult = await contract.castFromPreEncryptedToEuint64(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });
});

describe("Test AsEuint128", () => {
  let contract;
  let fheContract;

  it(`Test Contract Deployment`, async () => {
    const baseContract = await deployContract("AsEuint128Test");
    contract = baseContract as AsEuint128TestType;

    const contractAddress = await baseContract.getAddress();
    fheContract = await getFheContract(contractAddress);

    expect(contract).toBeTruthy();
    expect(fheContract).toBeTruthy();
  });

  const value = BigInt(1);
  const funcTypes = ["regular", "bound"];

  for (const funcType of funcTypes) {
    it(`From ebool - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEboolToEuint128(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint8 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint8ToEuint128(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint16 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint16ToEuint128(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint32 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint32ToEuint128(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint64 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint64ToEuint128(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint256 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint256ToEuint128(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  it(`From plaintext`, async () => {
    let decryptedResult = await contract.castFromPlaintextToEuint128(value);
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted`, async () => {
    const encInput = await fheContract.instance.encrypt_uint128(value);
    let decryptedResult = await contract.castFromPreEncryptedToEuint128(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_uint128(value, 1);
    let decryptedResult = await contract.castFromPreEncryptedToEuint128(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_uint128(value, 1);
    let decryptedResult = await contract.castFromPreEncryptedToEuint128(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });
});

describe("Test AsEuint256", () => {
  let contract;
  let fheContract;

  it(`Test Contract Deployment`, async () => {
    const baseContract = await deployContract("AsEuint256Test");
    contract = baseContract as AsEuint256TestType;

    const contractAddress = await baseContract.getAddress();
    fheContract = await getFheContract(contractAddress);

    expect(contract).toBeTruthy();
    expect(fheContract).toBeTruthy();
  });

  const value = BigInt(1);
  const funcTypes = ["regular", "bound"];

  for (const funcType of funcTypes) {
    it(`From ebool - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEboolToEuint256(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint8 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint8ToEuint256(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint16 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint16ToEuint256(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint32 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint32ToEuint256(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint64 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint64ToEuint256(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  for (const funcType of funcTypes) {
    it(`From euint128 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint128ToEuint256(
        value,
        funcType
      );
      expect(decryptedResult).toBe(value);
    });
  }

  it(`From plaintext`, async () => {
    let decryptedResult = await contract.castFromPlaintextToEuint256(value);
    expect(decryptedResult).toBe(value);
  });

  it(`From plaintext`, async () => {
    // We want to make sure that numbers that are higher than max of uint128 but lower than max of uint256 are handled correctly
    // Before a fix was introduced, there were some code sections that were converting the results of decryption to 64bits representations.
    // This was causing the numbers to be truncated and the results to be incorrect.
    // This test is to make sure that the fix is working correctly.
    // The number is 2^128 + 1 which is not representable in uint128 but is representable in uint256 and also is not aligned by 64 bits, previously, with the bug present the result was just 1
    const bigNumber = BigInt(2) ** BigInt(128) + BigInt(1);
    let decryptedResult = await contract.castFromPlaintextToEuint256(bigNumber);
    expect(decryptedResult).toBe(bigNumber);
  });

  it(`From pre encrypted`, async () => {
    const encInput = await fheContract.instance.encrypt_uint256(value);
    let decryptedResult = await contract.castFromPreEncryptedToEuint256(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_uint256(value, 1);
    let decryptedResult = await contract.castFromPreEncryptedToEuint256(
      encInput
    );
    expect(decryptedResult).toBe(value);
  });
});

describe("Test AsEaddress", () => {
  let contract;
  let fheContract;

  it(`Test Contract Deployment`, async () => {
    const baseContract = await deployContract("AsEaddressTest");
    contract = baseContract as AsEaddressTestType;

    const contractAddress = await baseContract.getAddress();
    fheContract = await getFheContract(contractAddress);

    expect(contract).toBeTruthy();
    expect(fheContract).toBeTruthy();
  });

  const value = BigInt(455113547441765951074000967332144802967768096399n); //0x4fB7FF4e004FcADbff708d8d873592B2044d5E8f
  const funcTypes = ["regular", "bound"];

  for (const funcType of funcTypes) {
    it(`From euint256 - ${funcType}`, async () => {
      let decryptedResult = await contract.castFromEuint256ToEaddress(
        value,
        funcType
      );
      let decimal = BigInt(decryptedResult);
      expect(decimal).toBe(value);
    });
  }

  it(`From plaintext`, async () => {
    let decryptedResult = await contract.castFromPlaintextToEaddress(value);
    let decimal = BigInt(decryptedResult);
    expect(decimal).toBe(value);
  });

  it(`From plaintext`, async () => {
    // We want to make sure that numbers that are higher than max of uint128 but lower than max of uint256 are handled correctly
    // Before a fix was introduced, there were some code sections that were converting the results of decryption to 64bits representations.
    // This was causing the numbers to be truncated and the results to be incorrect.
    // This test is to make sure that the fix is working correctly.
    // The number is 2^128 + 1 which is not representable in uint128 but is representable in uint256 and also is not aligned by 64 bits, previously, with the bug present the result was just 1
    const bigNumber = BigInt(2) ** BigInt(128) + BigInt(1);
    let decryptedResult = await contract.castFromPlaintextToEaddress(bigNumber);
    let decimal = BigInt(decryptedResult);
    expect(decimal).toBe(bigNumber);
  });

  it(`From plaintext address`, async () => {
    let decryptedResult = await contract.castFromPlaintextAddressToEaddress(
      "0x4fB7FF4e004FcADbff708d8d873592B2044d5E8f"
    );
    let decimal = BigInt(decryptedResult);
    expect(decimal).toBe(value);
  });

  it(`From pre encrypted`, async () => {
    const encInput = await fheContract.instance.encrypt_address(value);
    let decryptedResult = await contract.castFromPreEncryptedToEaddress(
      encInput
    );
    let decimal = BigInt(decryptedResult);
    expect(decimal).toBe(value);
  });

  it(`From pre encrypted - Security Zone 1`, async () => {
    const encInput = await fheContract.instance.encrypt_address(value, 1);
    let decryptedResult = await contract.castFromPreEncryptedToEaddress(
      encInput
    );
    let decimal = BigInt(decryptedResult);
    expect(decimal).toBe(value);
  });
});
