// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract ShlTest {
    using Utils for *;
    
    function shl(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "shl(euint8,euint8)")) {
            return FHE.decrypt(FHE.shl(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "shl(euint16,euint16)")) {
            return FHE.decrypt(FHE.shl(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "shl(euint32,euint32)")) {
            return FHE.decrypt(FHE.shl(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.shl(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).shl(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.shl(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).shl(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.shl(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).shl(FHE.asEuint32(b)));
        }
    
        revert TestNotFound(test);
    }
}
