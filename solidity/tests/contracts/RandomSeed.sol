// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import { FHE, euint64, Common } from "../../FHE.sol";

contract RandomSeedTest {
    function randomEuint64(uint64 seed) internal pure returns (euint64) {
        uint256 result = FHE.random(Common.EUINT64_TFHE, seed, 0);
        return euint64.wrap(result);
    }

    function randomSeed(uint64 seed) public pure returns (uint64) {
        return randomEuint64(seed).decrypt();
    }
}
