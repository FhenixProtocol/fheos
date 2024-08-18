// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract RandomTest {
    using Utils for *;
    
    function random(string calldata test) public pure returns (uint256) {
        if (Utils.cmp(test, "random_euint8()")) {
            return FHE.decrypt(FHE.randomEuint8());
        } else if (Utils.cmp(test, "random_euint16()")) {
            return FHE.decrypt(FHE.randomEuint16());
        } else if (Utils.cmp(test, "random_euint32()")) {
            return FHE.decrypt(FHE.randomEuint32());
        } else if (Utils.cmp(test, "random_euint64()")) {
            return FHE.decrypt(FHE.randomEuint64());
        } else if (Utils.cmp(test, "random_euint128()")) {
            return FHE.decrypt(FHE.randomEuint128());
        } else if (Utils.cmp(test, "random_euint256()")) {
            return FHE.decrypt(FHE.randomEuint256());
        }
        revert TestNotFound(test);
    }
}
