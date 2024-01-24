import { BaseContract } from 'ethers';
import { FhenixClient, getPermit, Permit } from 'fhenixjs';
import { ethers } from 'hardhat';

export interface FheContract {
  instance: FhenixClient;
  permit: Permit;
}

export async function createFheInstance(
  contractAddress: string
): Promise<FheContract> {
  const provider = ethers.provider;

  // Get the chainId
  //const fhenix = await provider.getNetwork();
  //const chainId = fhenix.chainId;
  let instance = new FhenixClient({ provider });
  const permit = await getPermit(contractAddress, provider);
  instance.storePermit(permit);
  // // workaround for call not working the first time on a fresh chain
  // let fhePublicKey = await ethers.provider.send("eth_getNetworkPublicKey");
  // const instance = createInstance({ chainId: Number(chainId), publicKey: fhePublicKey });
  // const genTokenResponse = instance.then((ins) => {
  //     return ins.generateToken({ verifyingContract: contractAddress });
  // });

  return { instance, permit };
}

export const fromHexString = (hexString: string): Uint8Array => {
  const arr = hexString.replace(/^(0x)/, '').match(/.{1,2}/g);
  if (!arr) return new Uint8Array();
  return Uint8Array.from(arr.map((byte) => parseInt(byte, 16)));
};

export const deployContract = async (contractName: string, args?: any[]) => {
  const [signer] = await ethers.getSigners();
  const con = await ethers.getContractFactory(contractName);
  let deployedContract: BaseContract;
  try {
    deployedContract = await deployContractFromSigner(con, signer, undefined, args);
  } catch (e) {
    if (`${e}`.includes('nonce too')) {
      // find last occurence of ": " in e and get the number that comes after
      const match = `${e}`.match(/state: (\d+)/);
      const stateNonce = match ? parseInt(match[1], 10) : null;
      if (stateNonce === null) {
        throw new Error('Could not find nonce in error');
      }

      deployedContract = await syncNonce(con, signer, stateNonce);
    } else {
      throw e;
    }
  }

  const contract = deployedContract.connect(signer);
  return contract;
};

export const deployContractFromSigner = async (
  con: any,
  signer: any,
  nonce?: number,
  args?: any[]
) => {

  let argsToUse = args || [];

  return await con.deploy(...argsToUse, {
    from: signer,
    log: true,
    skipIfAlreadyDeployed: false,
    nonce,
  });
};

export const syncNonce = async (con: any, signer: any, stateNonce: number) => {
  console.log(`Syncing nonce to ${stateNonce}`);
  try {
    await deployContractFromSigner(con, signer, stateNonce);
  } catch (e) {
    console.log('Fixed nonce issue');
  }

  return await deployContractFromSigner(con, signer);
};
