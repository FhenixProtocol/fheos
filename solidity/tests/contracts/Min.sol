// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract MinTest {
    using Utils for *;
    
    function min(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "min(euint8,euint8)")) {
            return FHE.decrypt(FHE.min(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "min(euint16,euint16)")) {
            return FHE.decrypt(FHE.min(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "min(euint32,euint32)")) {
            return FHE.decrypt(FHE.min(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.min(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).min(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.min(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).min(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.min(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).min(FHE.asEuint32(b)));
        }
    
        revert TestNotFound(test);
    }
}
