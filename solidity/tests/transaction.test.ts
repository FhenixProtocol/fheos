import { AddCaller, AddCallee } from "../types/tests/contracts/Tx.sol";
import { Ownership } from "../types/tests/contracts/Ownership";
import { createFheInstance, deployContract } from "./utils";
import { fail } from "assert";

describe("Test Transactions Scenarios", () => {
  let contractCaller: AddCaller;
  let contractCallee: AddCallee;
  let contractAddr: string;

  beforeAll(async () => {
    contractCallee = (await deployContract("AddCallee")) as AddCallee;
    expect(contractCallee).toBeTruthy();

    contractCaller = (await deployContract("AddCaller", [
      contractCallee,
    ])) as AddCaller;
    expect(contractCaller).toBeTruthy();
    contractAddr = await contractCaller.getAddress();
    expect(contractAddr).toBeTruthy();
  });

  beforeEach(async () => {
    const encCounterResponse = await contractCaller.resetCounter();
    await encCounterResponse.wait();
  });

  it("Contract Deployment", async () => {
    contractCallee = (await deployContract("AddCallee")) as AddCallee;
    expect(contractCallee).toBeTruthy();

    contractCaller = (await deployContract("AddCaller", [
      contractCallee,
    ])) as AddCaller;
    expect(contractCaller).toBeTruthy();
    contractAddr = await contractCaller.getAddress();
    expect(contractAddr).toBeTruthy();
  });

  it("Basic Add Tx", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encInput = await instance.encrypt_uint32(8);

    // 1 - static call
    const encCounter = await contractCaller.addTx.staticCall(
      encInput,
      permit.publicKey
    );
    const counterStatic = instance.unseal(contractAddr, encCounter);
    expect(Number(counterStatic)).toEqual(8);

    // 2 - real call + query
    const encCounterReceipt = await contractCaller.addTx(
      encInput,
      permit.publicKey
    );
    await encCounterReceipt.wait();

    const getCounterResponse = await contractCaller.getCounter(
      permit.publicKey
    );
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(8);
  });

  it("Sub via contract call as Plaintext", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);

    // Set the counter to 10. We need it to be positive, so we can subtract from it.
    // We can't add to it when it's decrypted, because it will overflow on gas estimation (decrypt shortcuts with MAX_UINT).
    const setReceipt = await contractCaller.addTx(
      await instance.encrypt_uint32(1337 + 9),
      permit.publicKey
    );
    await setReceipt.wait();

    // 1 - static call
    const encCounter =
      await contractCaller.subViaContractCallAsPlain.staticCall(
        9,
        permit.publicKey
      );
    const counterStatic = instance.unseal(contractAddr, encCounter);
    expect(Number(counterStatic)).toEqual(1337);

    // 2 - real call + query
    const encCounterReceipt = await contractCaller.subViaContractCallAsPlain(
      9,
      permit.publicKey
    );
    await encCounterReceipt.wait();

    const getCounterResponse = await contractCaller.getCounter(
      permit.publicKey
    );
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(1337);
  });

  it("Add via contract call - pass pass InEuint32", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encInput = await instance.encrypt_uint32(10);

    // 1 - static call
    const encCounter = await contractCaller.addViaContractCall.staticCall(
      encInput,
      permit.publicKey
    );
    const counterStatic = instance.unseal(contractAddr, encCounter);
    expect(Number(counterStatic)).toEqual(10);

    // 2 - real call + query
    const encCounterReceipt = await contractCaller.addViaContractCall(
      encInput,
      permit.publicKey
    );
    await encCounterReceipt.wait();

    const getCounterResponse = await contractCaller.getCounter(
      permit.publicKey
    );
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(10);
  });

  it("Add via contract call - pass Euint32", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encInput = await instance.encrypt_uint32(11);

    // 1 - static call
    const encCounter = await contractCaller.addViaContractCallU32.staticCall(
      encInput,
      permit.publicKey
    );
    const counterStatic = instance.unseal(contractAddr, encCounter);
    expect(Number(counterStatic)).toEqual(11);

    // 2 - real call + query
    const encCounterReceipt = await contractCaller.addViaContractCallU32(
      encInput,
      permit.publicKey
    );
    await encCounterReceipt.wait();

    const getCounterResponse = await contractCaller.getCounter(
      permit.publicKey
    );
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(11);
  });

  it("Add via VIEW contract call - pass Euint32", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encInput = await instance.encrypt_uint32(11);

    // 1 - static call
    const encCounter =
      await contractCaller.addViaViewContractCallU32.staticCall(
        encInput,
        permit.publicKey
      );
    const counterStatic = instance.unseal(contractAddr, encCounter);
    expect(Number(counterStatic)).toEqual(16);

    // 2 - real call + query
    const encCounterReceipt = await contractCaller.addViaViewContractCallU32(
      encInput,
      permit.publicKey
    );
    await encCounterReceipt.wait();

    const getCounterResponse = await contractCaller.getCounter(
      permit.publicKey
    );
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(16);
  });

  it("Add via DELEGATE contract call as Plaintext", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);

    // 1 - static call
    const encCounter = await contractCaller.addDelegatePlain.staticCall(
      12,
      permit.publicKey
    );
    const counterStatic = instance.unseal(contractAddr, encCounter);
    expect(Number(counterStatic)).toEqual(24);

    // 2 - real call + query
    const encCounterReceipt = await contractCaller.addDelegatePlain(
      12,
      permit.publicKey
    );
    await encCounterReceipt.wait();

    const getCounterResponse = await contractCaller.getCounter(
      permit.publicKey
    );
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(24);
  });

  it("Add via DELEGATE contract call - pass Euint32", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encInput = await instance.encrypt_uint32(13);

    // 1 - static call
    const encCounter = await contractCaller.addDelegate.staticCall(
      encInput,
      permit.publicKey
    );
    const counterStatic = instance.unseal(contractAddr, encCounter);
    expect(Number(counterStatic)).toEqual(26);

    // 2 - real call + query
    const encCounterReceipt = await contractCaller.addDelegate(
      encInput,
      permit.publicKey
    );
    await encCounterReceipt.wait();

    const getCounterResponse = await contractCaller.getCounter(
      permit.publicKey
    );
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(26);
  });

  it("Specify valid Security Zone (eth call, not tx)", async () => {
    const { permit, instance } = await createFheInstance(contractAddr);

    let encRes;
    try {
      encRes = await contractCaller.addPlainSecurityZone(
        1299,
        38,
        1,
        permit.publicKey
      );
    } catch (e) {
      console.error(`failed operation on securityzone 1: ${e}`);
      fail("Should not have reverted");
    }

    const result = instance.unseal(contractAddr, encRes);
    expect(Number(result)).toEqual(1337);
  });

  it("Specify invalid Security Zone (eth call, not tx)", async () => {
    const { permit, instance } = await createFheInstance(contractAddr);

    try {
      await contractCaller.addPlainSecurityZone(
        1299,
        38,
        3,
        permit.publicKey
      );

      fail("Should have reverted");
    } catch (err) {
      expect(err.message).toContain("execution reverted");
    }
  });

  it("Add via DELEGATE contract call - pass Uint32", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encInput = await instance.encrypt_uint32(14);

    // 1 - static call
    const encCounter = await contractCaller.addDelegate.staticCall(
      encInput,
      permit.publicKey
    );
    const counterStatic = instance.unseal(contractAddr, encCounter);
    expect(Number(counterStatic)).toEqual(28);

    // 2 - real call + query
    const encCounterReceipt = await contractCaller.addDelegateU32(
      encInput,
      permit.publicKey
    );
    await encCounterReceipt.wait();

    const getCounterResponse = await contractCaller.getCounter(
      permit.publicKey
    );
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(28);
  });

  it("sStore sanity check", async () => {
    try {
      const tx = await contractCaller.sStoreSanity();
      await tx.wait();
    } catch (e) {
      console.error("failed sstore sanity check");
      fail("Should not have reverted");
    }
  });

  it("Random sanity check", async () => {
    let result;
    try {
      const tx = await contractCaller.randomSanity();
      await tx.wait();
      const filter = contractCaller.filters.RandomSanityEvent
      const events = await contractCaller.queryFilter(filter, -1)
      result = events[0].args;
    } catch (e) {
      console.error("failed sstore sanity check");
      fail("Should not have reverted");
    }

    expect(result[0] === result[1]).toBeFalsy();
  });

  // function addDelegate(inEuint32 calldata value, bytes32 publicKey) public returns (bytes memory) {
});

describe("Test CT Ownership", () => {
  let methods = [
    { name: "callTest", stateExpectation: 2, expectedResponse: 3 },
    { name: "staticTest", stateExpectation: 1, expectedResponse: 4 },
    { name: "queryTest", stateExpectation: 1, expectedResponse: 2 },
    { name: "delegateTest", stateExpectation: 1, expectedResponse: 3 },
  ];
  let expectedCounter = 1;

  let contracts = new Array<Ownership>(methods.length + 1);
  let contractAddrs = new Array<string>(methods.length + 1);
  beforeAll(async () => {
    // contracts.length = methods.length + 1
    for (let i = 0; i <= methods.length; i++) {
      const contract = (await deployContract("Ownership")) as Ownership;
      expect(contract).toBeTruthy();

      contracts[i] = contract;

      const contractAddr = await contract.getAddress();
      expect(contractAddr).toBeTruthy();

      contractAddrs[i] = contractAddr;
    }
  });

  it.skip("Full chain", async () => {
    for (let i = 0; i < methods.length; i++) {
      const contract = contracts[i];
      const method = methods[i];

      let response = await contract.resetContractResponse();
      await response.wait();

      response = await contract[method.name](contracts[i + 1]);
      await response.wait();

      const returnCounter = await contract.getResp();
      expect(Number(returnCounter)).toEqual(method.expectedResponse);

      const counter = await contracts[i + 1].get();
      expect(Number(counter)).toEqual(method.stateExpectation);
    }
  });

  it("Broken chain", async () => {
    const caller = contracts[2];
    const callee = contracts[3];
    const { instance, permit } = await createFheInstance(contractAddrs[2]);
    const encInput = await instance.encrypt_uint32(100);
    const tx = await caller.setEnc(encInput);
    await tx.wait();
    const response = await caller.getEnc();
    let failed = false;
    try {
      const resp = await callee.set(response);
      let result = await resp.wait();
    } catch (err) {
      failed = true;
    }

    expect(failed).toBeTruthy();
  });
});
