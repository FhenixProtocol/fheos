// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract MaxTest {
    using Utils for *;
    
    function max(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "max(euint8,euint8)")) {
            return FHE.decrypt(FHE.max(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "max(euint16,euint16)")) {
            return FHE.decrypt(FHE.max(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "max(euint32,euint32)")) {
            return FHE.decrypt(FHE.max(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.max(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).max(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.max(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).max(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.max(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).max(FHE.asEuint32(b)));
        }
    
        revert TestNotFound(test);
    }
}
