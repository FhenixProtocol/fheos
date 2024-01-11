// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract LtTest {
    using Utils for *;
    
    function lt(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "lt(euint8,euint8)")) {
            if (FHE.decrypt(FHE.lt(FHE.asEuint8(a), FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "lt(euint16,euint16)")) {
            if (FHE.decrypt(FHE.lt(FHE.asEuint16(a), FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "lt(euint32,euint32)")) {
            if (FHE.decrypt(FHE.lt(FHE.asEuint32(a), FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint8.lt(euint8)")) {
            if (FHE.decrypt(FHE.asEuint8(a).lt(FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.lt(euint16)")) {
            if (FHE.decrypt(FHE.asEuint16(a).lt(FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.lt(euint32)")) {
            if (FHE.decrypt(FHE.asEuint32(a).lt(FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        }
        revert TestNotFound(test);
    }
}
