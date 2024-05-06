// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AndTest {
    using Utils for *;
    
    function and(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "and(euint8,euint8)")) {
            return FHE.decrypt(FHE.and(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "and(euint16,euint16)")) {
            return FHE.decrypt(FHE.and(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "and(euint32,euint32)")) {
            return FHE.decrypt(FHE.and(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "and(euint64,euint64)")) {
            return FHE.decrypt(FHE.and(FHE.asEuint64(a), FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "and(euint128,euint128)")) {
            return FHE.decrypt(FHE.and(FHE.asEuint128(a), FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "euint8.and(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).and(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.and(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).and(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.and(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).and(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint64.and(euint64)")) {
            return FHE.decrypt(FHE.asEuint64(a).and(FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "euint128.and(euint128)")) {
            return FHE.decrypt(FHE.asEuint128(a).and(FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "euint8 & euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) & FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 & euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) & FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 & euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) & FHE.asEuint32(b));
        } else if (Utils.cmp(test, "euint64 & euint64")) {
            return FHE.decrypt(FHE.asEuint64(a) & FHE.asEuint64(b));
        } else if (Utils.cmp(test, "euint128 & euint128")) {
            return FHE.decrypt(FHE.asEuint128(a) & FHE.asEuint128(b));
        } else if (Utils.cmp(test, "and(ebool,ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.and(FHE.asEbool(aBool), FHE.asEbool(bBool)))) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "ebool.and(ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.asEbool(aBool).and(FHE.asEbool(bBool)).decrypt()) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "ebool & ebool")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.asEbool(aBool) & FHE.asEbool(bBool))) {
                return 1;
            }
            return 0;
        }
    
        revert TestNotFound(test);
    }
}
