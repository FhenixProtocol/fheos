// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import { TFHE } from "../../FHE.sol";
import { Utils } from "./utils/Utils.sol";

error TestNotFound(string test);

contract AndTest {
    using Utils for *;

    function and(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "and(euint8,euint8)")) {
            return TFHE.decrypt(TFHE.and(TFHE.asEuint8(a), TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "and(euint16,euint16)")) {
            return TFHE.decrypt(TFHE.and(TFHE.asEuint16(a), TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "and(euint32,euint32)")) {
            return TFHE.decrypt(TFHE.and(TFHE.asEuint32(a), TFHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.and(euint8)")) {
            return TFHE.decrypt(TFHE.asEuint8(a).and(TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.and(euint16)")) {
            return TFHE.decrypt(TFHE.asEuint16(a).and(TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.and(euint32)")) {
            return TFHE.decrypt(TFHE.asEuint32(a).and(TFHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8 & euint8")) {
            return TFHE.decrypt(TFHE.asEuint8(a) & TFHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 & euint16")) {
            return TFHE.decrypt(TFHE.asEuint16(a) & TFHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 & euint32")) {
            return TFHE.decrypt(TFHE.asEuint32(a) & TFHE.asEuint32(b));
        } else if (Utils.cmp(test, "and(ebool,ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (TFHE.decrypt(TFHE.and(TFHE.asEbool(aBool), TFHE.asEbool(bBool)))) {
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
            if (TFHE.decrypt(TFHE.asEbool(aBool).and(TFHE.asEbool(bBool)))) {
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
            if (TFHE.decrypt(TFHE.asEbool(aBool) & TFHE.asEbool(bBool))) {
                return 1;
            }
            return 0;
        } else {
            revert TestNotFound(test);
        }
    }

}