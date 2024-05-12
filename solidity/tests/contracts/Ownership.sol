// SPDX-License-Identifier: MIT
// solhint-disable one-contract-per-file
// solhint-disable avoid-low-level-calls
pragma solidity ^0.8.19;

import { FHE, euint32, inEuint32 } from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

contract Ownership {
    euint32 private counter = FHE.asEuint32(1);
    euint32 private contractResponse;

    function inc(euint32 c) public returns (euint32) {
        return FHE.add(c, FHE.asEuint32(1));
    }

    function setPlain(uint32 val) public {
        counter = FHE.asEuint32(val);
    }

    function resetContractResponse() public {
        contractResponse = FHE.asEuint32(1337);
    }

    function get() public view returns (uint256) {
        return FHE.decrypt(counter);
    }

    function getEnc() public view returns (uint256) {
        return euint32.unwrap(counter);
    }

    function set(uint256 value) public returns (euint32) {
        counter = euint32.wrap(value);
        return FHE.add(counter, FHE.asEuint32(1));
    }

    function setStatic(uint256 value) public returns (euint32) {
        return FHE.add(euint32.wrap(value), FHE.asEuint32(1));
    }

    function setQ(euint32 value) public view returns (euint32) {
        return FHE.add(value, FHE.asEuint32(1));
    }

    function getResp() public view returns (uint256) {
        return FHE.decrypt(contractResponse);
    }

    function callTest(address next) public {
        contractResponse =  Ownership(next).set(euint32.unwrap(inc(counter)));
    }

    function staticTest(address next) public  {
        (bool success, bytes memory ret) = next.staticcall(abi.encodeWithSignature("setStatic(uint256)", euint32.unwrap(inc(counter))));
        require(success, "staticcall failed");

        contractResponse = abi.decode(ret, (euint32));
    }

    function delegateTest(address next) public {
        (bool success, bytes memory ret) = next.delegatecall(abi.encodeWithSignature("set(uint256)", euint32.unwrap(inc(counter))));
        require(success, "delegatecall failed");

        contractResponse = abi.decode(ret, (euint32));
    }

    function queryTest(address next) public returns (uint256) {
        contractResponse =  Ownership(next).setQ(inc(counter));
    }
}
