import {WrappingERC20} from "../types/tests/contracts/wERC20.sol";
import { createFheInstance, deployContract } from "./utils";
import { ethers } from 'hardhat';

describe('Test WERC20', () =>  {
    let contract: WrappingERC20;
    let contractAddr: string;
    const amountToSend = BigInt(1);

    // We don't really need it as test but it is a test since it is async
    it(`Test Contract Deployment`, async () => {
        const baseContract = await deployContract('WrappingERC20', ["TEST", "TST"]);
        contract = baseContract as WrappingERC20;

        contractAddr = await baseContract.getAddress();

        expect(contract).toBeTruthy();
    });


    it(`Execute Transaction`, async () => {
        const [signer] = await ethers.getSigners();
        const { instance, permit } = await createFheInstance(contractAddr);

        let signerAddress = await signer.getAddress();

        const encrypted = await instance.encrypt_uint32(Number(amountToSend));

        let decryptedResult = await contract.transferEncrypted(signerAddress, encrypted);

        // for (let i = 0; i < 2; i++) {
        // }
    });
});

