// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract GteTest {
    using Utils for *;
    
    function gte(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "gte(euint8,euint8)")) {
            if (FHE.decrypt(FHE.gte(FHE.asEuint8(a), FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "gte(euint16,euint16)")) {
            if (FHE.decrypt(FHE.gte(FHE.asEuint16(a), FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "gte(euint32,euint32)")) {
            if (FHE.decrypt(FHE.gte(FHE.asEuint32(a), FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint8.gte(euint8)")) {
            if (FHE.decrypt(FHE.asEuint8(a).gte(FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.gte(euint16)")) {
            if (FHE.decrypt(FHE.asEuint16(a).gte(FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.gte(euint32)")) {
            if (FHE.decrypt(FHE.asEuint32(a).gte(FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        }
        revert TestNotFound(test);
    }
}
