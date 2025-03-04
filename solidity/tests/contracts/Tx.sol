// SPDX-License-Identifier: MIT
// solhint-disable one-contract-per-file
// solhint-disable avoid-low-level-calls
pragma solidity ^0.8.19;

import { FHE, euint32, inEuint32 } from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AddCaller {
    // Don't change the order of these variables!
    // Two things must be preserved:
    // 1. The counter variable must be the first one, to allow delegate calls
    // 2. The counterPublic has to come right before the addContract variable, which can't be initialized.
    //    Only this catches the sstore bug
    euint32 private counter = FHE.asEuint32(0);
    uint64 private counterPublic = 0;
    AddCallee public addContract;
    error DelegateCallFailed();
    event RandomSanityEvent(uint32 a, uint32 b);

    constructor(address addCallee) {
        addContract = AddCallee(addCallee);
    }

    function resetCounter() public {
        counter = FHE.asEuint32(0);
    }

    function addTx(inEuint32 calldata value, bytes32 publicKey) public returns (string memory) {
        counter = counter.add(FHE.asEuint32(value));
        return counter.seal(publicKey);
    }

    function addPlainSecurityZone(uint256 lhs, uint256 rhs, int32 securityZone, bytes32 publicKey) public pure returns (string memory) {
        euint32 lhsEncrypted = FHE.asEuint32(lhs, securityZone);
        euint32 rhsEncrypted = FHE.asEuint32(rhs, securityZone);
        euint32 result = lhsEncrypted.add(rhsEncrypted);
        return result.seal(publicKey);
    }

    function subViaContractCallAsPlain(uint32 value, bytes32 publicKey) public returns (string memory) {
        counter = addContract.sub(FHE.decrypt(counter), value);
        return counter.seal(publicKey);
    }

    function addViaContractCall(inEuint32 calldata value, bytes32 publicKey) public returns (string memory) {
        counter = addContract.add(counter, value);
        return counter.seal(publicKey);
    }

    function addViaContractCallU32(inEuint32 calldata value, bytes32 publicKey) public returns (string memory) {
        counter = addContract.add(counter, FHE.asEuint32(value));
        return counter.seal(publicKey);
    }

    function addViaViewContractCallU32(inEuint32 calldata value, bytes32 publicKey) public returns (string memory) {
        counter = addContract.addView(counter, FHE.asEuint32(value));
        return counter.seal(publicKey);
    }

    function addDelegatePlain(uint32 value, bytes32 publicKey) public returns (string memory) {
        // value is added twice in this case, to verify that the delegatecall operated on the correct value
        counter = counter.add(FHE.asEuint32(value));

        (bool success, /* bytes memory data */) = address(addContract).delegatecall(
            abi.encodeWithSelector(
                AddCallee.addDelegatePlain.selector,
                value
            )
        );

        if (!success) {
            revert DelegateCallFailed();
        }
        return counter.seal(publicKey);
    }

    function addDelegate(inEuint32 calldata value, bytes32 publicKey) public returns (string memory) {
        // value is added twice, to verify that the delegatecall operated on the correct value
        counter = counter.add(FHE.asEuint32(value));

        (bool success, /* bytes memory data */) = address(addContract).delegatecall(
            abi.encodeWithSelector(
                AddCallee.addDelegateInEuint.selector,
                value
            )
        );

        if (!success) {
            revert DelegateCallFailed();
        }
        return counter.seal(publicKey);
    }

    function addDelegateU32(inEuint32 calldata value, bytes32 publicKey) public returns (string memory) {
        // value is added twice, to verify that the delegatecall operated on the correct value
        counter = counter.add(FHE.asEuint32(value));

        (bool success, /* bytes memory data */) = address(addContract).delegatecall(
            abi.encodeWithSelector(
                AddCallee.addDelegateEuint.selector,
                FHE.asEuint32(value)
            )
        );

        if (!success) {
            revert DelegateCallFailed();
        }
        return counter.seal(publicKey);
    }

    // utility funcs
    function returnPlainCounter() public view returns (uint256) {
        return FHE.decrypt(counter);
    }

    function getCounter(bytes32 publicKey) public view returns (string memory) {
        return counter.seal(publicKey);
    }

    function sStoreSanity() public {
        // catches a bug we had a in the past on the sstore hook
        counterPublic = counterPublic + 1;
    }
}

contract AddCallee {
    euint32 private counter;
    euint32 private constForView = FHE.asEuint32(5);

    function add(euint32 a, inEuint32 calldata b) public pure returns (euint32 output) {
        return a + FHE.asEuint32(b);
    }

    function add(euint32 a, euint32 b) public pure returns (euint32 output) {
        return a + b;
    }

    function addView(euint32 a, euint32 b) public view returns (euint32 output) {
        return a + b + constForView;
    }

    // Note: sub is needed as a workaround when working with decrypted integers on tests.
    // There's an issue where gas estimation reverts with overflow.
    // Because, as a shortcut, "a" is decrypted to max_uint.
    function sub(uint32 a, uint32 b) public pure returns (euint32 output) {
        return FHE.asEuint32(a - b);
    }

    function addDelegatePlain(uint32 value) public {
        counter = counter + FHE.asEuint32(value);
    }

    function addDelegateEuint(euint32 value) public {
        counter = counter + value;
    }

    function addDelegateInEuint(inEuint32 calldata value) public {
        counter = counter + FHE.asEuint32(value);
    }
}
