import { Counter } from "../types/tests/contracts/Counter";
import { createFheInstance, deployContract } from "./utils";

describe("Test Unitialized Variables", () => {
  let contract: Counter;
  let contractAddr: string;

  it("Contract Deployment", async () => {
    contract = (await deployContract("Counter")) as Counter;
    expect(contract).toBeTruthy();
    contractAddr = await contract.getAddress();
    expect(contractAddr).toBeTruthy();
  });

  it("Test reading uninitialized state variable", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encCounter = await contract.getCounter(permit.publicKey);
    const counter = await instance.unseal(contractAddr, encCounter);
    expect(Number(counter)).toEqual(0);
  });

  it("Test using uninitialized state variable", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);

    const encAmount = instance.encrypt_uint8(33);
    await contract.add(encAmount);

    const encCounter = await contract.getCounter(permit.publicKey);
    const counter = await instance.unseal(contractAddr, encCounter);
    expect(Number(counter)).toEqual(33); // Uninitialized state variable should be considered as 0
  });

  it("Test reading uninitialized mapping state variable", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);

    let encCounter = await contract.getCounterMapping(
      BigInt(0),
      permit.publicKey
    );
    let counter = await instance.unseal(contractAddr, encCounter);
    expect(Number(counter)).toEqual(0);

    // Try a different index
    encCounter = await contract.getCounterMapping(BigInt(12), permit.publicKey);
    counter = await instance.unseal(contractAddr, encCounter);
    expect(Number(counter)).toEqual(0);
  });

  it("Test using uninitialized mapping state variable", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);

    const encAmount = instance.encrypt_uint8(34);
    await contract.addMapping(BigInt(0), encAmount);

    const encCounter = await contract.getCounterMapping(
      BigInt(0),
      permit.publicKey
    );
    let counter = await instance.unseal(contractAddr, encCounter);
    expect(Number(counter)).toEqual(34); // Uninitialized state variable should be considered as 0

    // Try a different index
    const diffCounter = await contract.getCounterMapping(
      BigInt(12),
      permit.publicKey
    );
    counter = await instance.unseal(contractAddr, diffCounter);
    expect(Number(counter)).toEqual(0);
  });
});
