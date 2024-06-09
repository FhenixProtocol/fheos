// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
import {
    euint256, inEuint256
} from "../../../FHE.sol";

contract AsEaddressBench {
    euint256 internal a256;
    uint256 internal aUint256;
    bytes internal aBytes;

    function load256(inEuint256 calldata _a) public {
        a256 = FHE.asEuint256(_a);
    }
    function loadUint256(uint256 _a) public {
        aUint256 = _a;
    }
    function loadBytes(bytes memory _a) public {
        aBytes = _a;
    }

    function benchCastEuint256ToEaddress() public view {
        FHE.asEaddress(a256);
    }
    function benchCastUint256ToEaddress() public view {
        FHE.asEaddress(aUint256);
    }
    function benchCastBytesToEaddress() public view {
        FHE.asEaddress(aBytes);
    }
}
