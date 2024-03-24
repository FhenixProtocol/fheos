import { AddCaller, AddCallee } from "../types/tests/contracts/Tx.sol";
import { createFheInstance, deployContract } from "./utils";
import {ethers} from "hardhat";

describe.only("Test Transactions Scenarios", () => {
  let contractCaller: AddCaller;
  let contractCallee: AddCallee;
  let contractAddr: string;

  beforeAll(async () => {
    contractCallee = (await deployContract("AddCallee")) as AddCallee;
    expect(contractCallee).toBeTruthy();

    contractCaller = (await deployContract("AddCaller", [contractCallee])) as AddCaller;
    expect(contractCaller).toBeTruthy();
    contractAddr = await contractCaller.getAddress();
    expect(contractAddr).toBeTruthy();
  });

  beforeEach(async () => {
    const encCounterResponse = await contractCaller.resetCounter();
    await encCounterResponse.wait();
  });

  it.only("Contract Deployment", async () => {
    contractCallee = (await deployContract("AddCallee")) as AddCallee;
    expect(contractCallee).toBeTruthy();

    contractCaller = (await deployContract("AddCaller", [contractCallee])) as AddCaller;
    expect(contractCaller).toBeTruthy();
    contractAddr = await contractCaller.getAddress();
    expect(contractAddr).toBeTruthy();
  });

  it.only("Test simple Tx - static call ", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encInput = await instance.encrypt_uint32(7);
    const encCounterResponse = await contractCaller.addTx.staticCall(encInput, permit.publicKey);
    const counter = instance.unseal(contractAddr, encCounterResponse);
    expect(Number(counter)).toEqual(7);
  });

  it.only("Test simple Tx - query for state change", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encInput = await instance.encrypt_uint32(8);
    const encCounterResponse = await contractCaller.addTx(encInput, permit.publicKey);
    await encCounterResponse.wait();

    const getCounterResponse = await contractCaller.getCounter(permit.publicKey);
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(8);
  });
});
