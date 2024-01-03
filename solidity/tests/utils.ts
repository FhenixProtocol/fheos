import {FhenixClient, getPermit, Permit} from "fhenix.js";
import { ethers } from 'hardhat';

export interface FheContract {
    instance: FhenixClient;
    permit: Permit;
}

export async function createFheInstance(contractAddress: string): Promise<FheContract> {
    const provider = ethers.provider;

    // Get the chainId
    //const fhenix = await provider.getNetwork();
    //const chainId = fhenix.chainId;
    let instance = await FhenixClient.Create({provider});
    const permit = await getPermit(contractAddress, provider);
    instance.storePermit(permit);
    // // workaround for call not working the first time on a fresh chain
    // let fhePublicKey = await ethers.provider.send("eth_getNetworkPublicKey");
    // const instance = createInstance({ chainId: Number(chainId), publicKey: fhePublicKey });
    // const genTokenResponse = instance.then((ins) => {
    //     return ins.generateToken({ verifyingContract: contractAddress });
    // });

    return {instance, permit}
}

export const fromHexString = (hexString: string): Uint8Array => {
    const arr = hexString.replace(/^(0x)/, '').match(/.{1,2}/g);
    if (!arr) return new Uint8Array();
    return Uint8Array.from(arr.map((byte) => parseInt(byte, 16)));
};
