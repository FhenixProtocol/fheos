// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract ShrTest {
    using Utils for *;
    
    function shr(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "shr(euint8,euint8)")) {
            return FHE.decrypt(FHE.shr(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "shr(euint16,euint16)")) {
            return FHE.decrypt(FHE.shr(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "shr(euint32,euint32)")) {
            return FHE.decrypt(FHE.shr(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.shr(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).shr(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.shr(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).shr(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.shr(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).shr(FHE.asEuint32(b)));
        }
    
        revert TestNotFound(test);
    }
}
