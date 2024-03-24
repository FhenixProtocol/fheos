// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import { FHE, euint32, inEuint32 } from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AddCaller {
    euint32 private counter = FHE.asEuint32(0);
    AddCallee public addContract;

    constructor(address addCallee) {
        addContract = AddCallee(addCallee);
    }

    function resetCounter() public {
        counter = FHE.asEuint32(0);
    }

    function addTx(inEuint32 calldata value, bytes32 publicKey) public returns (bytes memory) {
        counter = counter.add(FHE.asEuint32(value));
        return counter.seal(publicKey);
    }

    function getCounter(bytes32 publicKey) public view returns (bytes memory) {
        return counter.seal(publicKey);
    }

    function addViaContractCallAsPlain(uint32 value, bytes32 publicKey) public returns (bytes memory) {
        counter = addContract.add(FHE.decrypt(counter), value);
        return counter.seal(publicKey);
    }

    function addViaContractCall(inEuint32 calldata value, bytes32 publicKey) public returns (bytes memory) {
        counter = addContract.add(counter, value);
        return counter.seal(publicKey);
    }

    function addViaContractCallU32(inEuint32 calldata value, bytes32 publicKey) public returns (bytes memory) {
        counter = addContract.add(counter, FHE.asEuint32(value));
        return counter.seal(publicKey);
    }

    function addDelegatePlain(uint32 value, bytes32 publicKey) public returns (bytes memory) {
        address(addContract).delegatecall(
            abi.encodeWithSignature("addDelegate(uint32)", value)
        );

        return counter.seal(publicKey);
    }

    function addDelegate(inEuint32 calldata value, bytes32 publicKey) public returns (bytes memory) {
        address(addContract).delegatecall(
            abi.encodeWithSignature("addDelegate(inEuint32)", value)
        );

        return counter.seal(publicKey);
    }

    function addDelegateU32(inEuint32 calldata value, bytes32 publicKey) public returns (bytes memory) {
        address(addContract).delegatecall(
            abi.encodeWithSignature("addDelegate(euint32)", FHE.asEuint32(value))
        );

        return counter.seal(publicKey);
    }

    //utility func
    function returnPlainCounter() public view returns (uint256) {
        return FHE.decrypt(counter);
    }
}

contract AddCallee {
    euint32 private counter;

    function add(euint32 a, inEuint32 calldata b) public pure returns (euint32 output) {
        return a + FHE.asEuint32(b);
    }

    function add(euint32 a, euint32 b) public pure returns (euint32 output) {
        return a + b;
    }

    function add(uint32 a, uint32 b) public pure returns (euint32 output) {
        // TODO: this fails for some reason. "a", a decrypted uint, cannot be used with "+"
        return FHE.asEuint32(a + b);
    }

    function addDelegate(euint32 value) public {
        counter = counter + value;
    }

    function addDelegate(inEuint32 calldata value) public {
        counter = counter + FHE.asEuint32(value);
    }
}
