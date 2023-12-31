// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract OrTest {
    using Utils for *;
    
    function or(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "or(euint8,euint8)")) {
            return FHE.decrypt(FHE.or(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "or(euint16,euint16)")) {
            return FHE.decrypt(FHE.or(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "or(euint32,euint32)")) {
            return FHE.decrypt(FHE.or(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.or(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).or(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.or(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).or(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.or(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).or(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8 | euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) | FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 | euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) | FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 | euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) | FHE.asEuint32(b));
        } else if (Utils.cmp(test, "or(ebool,ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.or(FHE.asEbool(aBool), FHE.asEbool(bBool)))) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "ebool.or(ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.asEbool(aBool).or(FHE.asEbool(bBool)))) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "ebool | ebool")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.asEbool(aBool) | FHE.asEbool(bBool))) {
                return 1;
            }
            return 0;
        }
    
        revert TestNotFound(test);
    }
}
