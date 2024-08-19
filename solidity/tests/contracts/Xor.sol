// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE, euint8, euint16, euint32, euint64, euint128, euint256, ebool} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract XorTest {
    using Utils for *;
    
    function xor(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "xor(euint8,euint8)")) {
            return FHE.decrypt(FHE.xor(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "xor(euint16,euint16)")) {
            return FHE.decrypt(FHE.xor(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "xor(euint32,euint32)")) {
            return FHE.decrypt(FHE.xor(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "xor(euint64,euint64)")) {
            return FHE.decrypt(FHE.xor(FHE.asEuint64(a), FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "xor(euint128,euint128)")) {
            return FHE.decrypt(FHE.xor(FHE.asEuint128(a), FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "euint8.xor(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).xor(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.xor(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).xor(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.xor(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).xor(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint64.xor(euint64)")) {
            return FHE.decrypt(FHE.asEuint64(a).xor(FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "euint128.xor(euint128)")) {
            return FHE.decrypt(FHE.asEuint128(a).xor(FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "euint8 ^ euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) ^ FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 ^ euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) ^ FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 ^ euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) ^ FHE.asEuint32(b));
        } else if (Utils.cmp(test, "euint64 ^ euint64")) {
            return FHE.decrypt(FHE.asEuint64(a) ^ FHE.asEuint64(b));
        } else if (Utils.cmp(test, "euint128 ^ euint128")) {
            return FHE.decrypt(FHE.asEuint128(a) ^ FHE.asEuint128(b));
        } else if (Utils.cmp(test, "xor(ebool,ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.xor(FHE.asEbool(aBool), FHE.asEbool(bBool)))) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "ebool.xor(ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.asEbool(aBool).xor(FHE.asEbool(bBool)).decrypt()) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "ebool ^ ebool")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.asEbool(aBool) ^ FHE.asEbool(bBool))) {
                return 1;
            }
            return 0;
        }
    
        revert TestNotFound(test);
    }
}
