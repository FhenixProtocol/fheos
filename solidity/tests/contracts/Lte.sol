// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract LteTest {
    using Utils for *;
    
    function lte(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "lte(euint8,euint8)")) {
            if (FHE.decrypt(FHE.lte(FHE.asEuint8(a), FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "lte(euint16,euint16)")) {
            if (FHE.decrypt(FHE.lte(FHE.asEuint16(a), FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "lte(euint32,euint32)")) {
            if (FHE.decrypt(FHE.lte(FHE.asEuint32(a), FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint8.lte(euint8)")) {
            if (FHE.decrypt(FHE.asEuint8(a).lte(FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.lte(euint16)")) {
            if (FHE.decrypt(FHE.asEuint16(a).lte(FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.lte(euint32)")) {
            if (FHE.decrypt(FHE.asEuint32(a).lte(FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        }
        revert TestNotFound(test);
    }
}
