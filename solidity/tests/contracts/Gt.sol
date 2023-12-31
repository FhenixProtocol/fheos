// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract GtTest {
    using Utils for *;
    
    function gt(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "gt(euint8,euint8)")) {
            if (FHE.decrypt(FHE.gt(FHE.asEuint8(a), FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "gt(euint16,euint16)")) {
            if (FHE.decrypt(FHE.gt(FHE.asEuint16(a), FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "gt(euint32,euint32)")) {
            if (FHE.decrypt(FHE.gt(FHE.asEuint32(a), FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint8.gt(euint8)")) {
            if (FHE.decrypt(FHE.asEuint8(a).gt(FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.gt(euint16)")) {
            if (FHE.decrypt(FHE.asEuint16(a).gt(FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.gt(euint32)")) {
            if (FHE.decrypt(FHE.asEuint32(a).gt(FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        }
        revert TestNotFound(test);
    }
}
