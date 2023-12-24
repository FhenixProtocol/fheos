// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {TFHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract EqTest {
    using Utils for *;
    
    function eq(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "eq(euint8,euint8)")) {
            if (TFHE.decrypt(TFHE.eq(TFHE.asEuint8(a), TFHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "eq(euint16,euint16)")) {
            if (TFHE.decrypt(TFHE.eq(TFHE.asEuint16(a), TFHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "eq(euint32,euint32)")) {
            if (TFHE.decrypt(TFHE.eq(TFHE.asEuint32(a), TFHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint8.eq(euint8)")) {
            if (TFHE.decrypt(TFHE.asEuint8(a).eq(TFHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.eq(euint16)")) {
            if (TFHE.decrypt(TFHE.asEuint16(a).eq(TFHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.eq(euint32)")) {
            if (TFHE.decrypt(TFHE.asEuint32(a).eq(TFHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "eq(ebool,ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (TFHE.decrypt(TFHE.eq(TFHE.asEbool(aBool), TFHE.asEbool(bBool)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ebool.eq(ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (TFHE.decrypt(TFHE.asEbool(aBool).eq(TFHE.asEbool(bBool)))) {
                return 1;
            }
        }
        revert TestNotFound(test);
    }
}
