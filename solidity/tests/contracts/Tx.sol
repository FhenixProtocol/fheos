// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AddCaller {
    euint32 private counter;
    
    function addTx(inEuint32 value, bytes calldata publicKey) public returns (bytes memory) {
        counter = counter.add(value.toU32());
        return counter.seal(publicKey);
    }

    function addViaContractCallAsPlain(
        uint32 value,
        address contractAddress,
        bytes calldata publicKey
    ) public returns (bytes memory) {
        AddCallee addContract = AddCallee(contractAddress);
        counter = addContract.add(FHE.decrypt(counter), value);

        return counter.seal(publicKey);
    }

    function addViaContractCall(
        inEuint32 value,
        address addContract,
        bytes calldata publicKey
    ) public returns (bytes memory) {
        AddCallee addContract = AddCallee(contractAddress);
        counter = addContract.add(counter, value);

        return counter.seal(publicKey);
    }

    function addViaContractCallU32(
        inEuint32 value,
        address addContract,
        bytes calldata publicKey
    ) public returns (bytes memory) {
        AddCallee addContract = AddCallee(contractAddress);
        counter = addContract.add(counter, value.toU32());

        return counter.seal(publicKey);
    }

    function addDelegate(
        inEuint32 value,
        address addContract,
        bytes calldata publicKey
    ) public returns (bytes memory) {
        AddCallee addContract = AddCallee(contractAddress);
        addContract.addDelegate(value);

        return counter.seal(publicKey);
    }

    function addDelegatePlain(
        uint32 value,
        address addContract,
        bytes calldata publicKey
    ) public returns (bytes memory) {
        AddCallee addContract = AddCallee(contractAddress);

        (bool success, bytes memory data) = addContract.delegatecall(
            abi.encodeWithSignature("addDelegate(uint32)", value.toU32())
        );

        return counter.seal(publicKey);
    }

    function addDelegate(
        inEuint32 value,
        address addContract,
        bytes calldata publicKey
    ) public returns (bytes memory) {
        AddCallee addContract = AddCallee(contractAddress);

        (bool success, bytes memory data) = addContract.delegatecall(
            abi.encodeWithSignature("addDelegate(inEuint32)", value)
        );

        return counter.seal(publicKey);
    }

    function addDelegateU32(
        inEuint32 value,
        address addContract,
        bytes calldata publicKey
    ) public returns (bytes memory) {
        AddCallee addContract = AddCallee(contractAddress);

        (bool success, bytes memory data) = addContract.delegatecall(
            abi.encodeWithSignature("addDelegate(euint32)", value.toU32())
        );

        return counter.seal(publicKey);
    }
}

contract AddCallee {
    euint32 private counter;

    function add(euint32 a, inEuint32 b) public pure returns (euint32 output) {
        return a + b.toU32();
    }

    function add(euint32 a, euint32 b) public pure returns (euint32 output) {
        return a + b;
    }

    function add(uint32 a, uint32 b) public pure returns (euint32 output) {
        return FHE.asEuint32(a + b);
    }

    function addDelegate(euint32) public {
        counter = counter + b;
    }

    function addDelegate(inEuint32) public {
        counter = counter + b.toU32();
    }
}
