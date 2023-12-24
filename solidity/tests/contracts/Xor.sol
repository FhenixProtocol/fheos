// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {TFHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract XorTest {
    using Utils for *;
    
    function xor(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "xor(euint8,euint8)")) {
            return TFHE.decrypt(TFHE.xor(TFHE.asEuint8(a), TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "xor(euint16,euint16)")) {
            return TFHE.decrypt(TFHE.xor(TFHE.asEuint16(a), TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "xor(euint32,euint32)")) {
            return TFHE.decrypt(TFHE.xor(TFHE.asEuint32(a), TFHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.xor(euint8)")) {
            return TFHE.decrypt(TFHE.asEuint8(a).xor(TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.xor(euint16)")) {
            return TFHE.decrypt(TFHE.asEuint16(a).xor(TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.xor(euint32)")) {
            return TFHE.decrypt(TFHE.asEuint32(a).xor(TFHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8 ^ euint8")) {
            return TFHE.decrypt(TFHE.asEuint8(a) ^ TFHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 ^ euint16")) {
            return TFHE.decrypt(TFHE.asEuint16(a) ^ TFHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 ^ euint32")) {
            return TFHE.decrypt(TFHE.asEuint32(a) ^ TFHE.asEuint32(b));
        } else if (Utils.cmp(test, "xor(ebool,ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (TFHE.decrypt(TFHE.xor(TFHE.asEbool(aBool), TFHE.asEbool(bBool)))) {
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
            if (TFHE.decrypt(TFHE.asEbool(aBool).xor(TFHE.asEbool(bBool)))) {
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
            if (TFHE.decrypt(TFHE.asEbool(aBool) ^ TFHE.asEbool(bBool))) {
                return 1;
            }
            return 0;
        } else {
            revert TestNotFound(test);
        }
    }
}
