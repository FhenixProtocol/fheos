// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AddTest {
    using Utils for *;
    
    function add(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "add(euint8,euint8)")) {
            return FHE.decrypt(FHE.add(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "add(euint16,euint16)")) {
            return FHE.decrypt(FHE.add(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "add(euint32,euint32)")) {
            return FHE.decrypt(FHE.add(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.add(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).add(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.add(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).add(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.add(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).add(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8 + euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) + FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 + euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) + FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 + euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) + FHE.asEuint32(b));
        }
    
        revert TestNotFound(test);
    }
}
