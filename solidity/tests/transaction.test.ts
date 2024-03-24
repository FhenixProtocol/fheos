import { AddCaller, AddCallee } from "../types/tests/contracts/Tx.sol";
import { createFheInstance, deployContract } from "./utils";

describe.only("Test Transactions Scenarios", () => {
  let contractCaller: AddCaller;
  let contractCallee: AddCallee;
  let contractAddr: string;

  it.only("Contract Deployment", async () => {
    contractCallee = (await deployContract("AddCallee")) as AddCallee;
    expect(contractCallee).toBeTruthy();

    contractCaller = (await deployContract("AddCaller", [contractCallee])) as AddCaller;
    expect(contractCaller).toBeTruthy();
    contractAddr = await contractCaller.getAddress();
    expect(contractAddr).toBeTruthy();
  });

  it.only("Test simple Tx - Call Static", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encInput = await instance.encrypt_uint32(7);
    const encCounterResponse = await contractCaller.addTx.staticCall(encInput, permit.publicKey);
    const counter = instance.unseal(contractAddr, encCounterResponse);
    console.log("counter", counter);
    expect(Number(counter)).toEqual(7);
  });

  // it("Test using uninitialized state variable", async () => {
  //   const { instance, permit } = await createFheInstance(contractAddr);
  //
  //   const encAmount = await instance.encrypt_uint8(33);
  //   await contract.add(encAmount.data);
  //
  //   const encCounter = await contract.getCounter(permit.publicKey);
  //   const counter = await instance.unseal(contractAddr, encCounter);
  //   expect(Number(counter)).toEqual(33); // Uninitialized state variable should be considered as 0
  // });
  //
  // it("Test reading uninitialized mapping state variable", async () => {
  //   const { instance, permit } = await createFheInstance(contractAddr);
  //
  //   let encCounter = await contract.getCounterMapping(
  //     BigInt(0),
  //     permit.publicKey
  //   );
  //   let counter = await instance.unseal(contractAddr, encCounter);
  //   expect(Number(counter)).toEqual(0);
  //
  //   // Try a different index
  //   encCounter = await contract.getCounterMapping(BigInt(12), permit.publicKey);
  //   counter = await instance.unseal(contractAddr, encCounter);
  //   expect(Number(counter)).toEqual(0);
  // });
  //
  // it("Test using uninitialized mapping state variable", async () => {
  //   const { instance, permit } = await createFheInstance(contractAddr);
  //
  //   const encAmount = await instance.encrypt_uint8(34);
  //   await contract.addMapping(BigInt(0), encAmount.data);
  //
  //   const encCounter = await contract.getCounterMapping(
  //     BigInt(0),
  //     permit.publicKey
  //   );
  //   let counter = await instance.unseal(contractAddr, encCounter);
  //   expect(Number(counter)).toEqual(34); // Uninitialized state variable should be considered as 0
  //
  //   // Try a different index
  //   const diffCounter = await contract.getCounterMapping(
  //     BigInt(12),
  //     permit.publicKey
  //   );
  //   counter = await instance.unseal(contractAddr, diffCounter);
  //   expect(Number(counter)).toEqual(0);
  // });
});
