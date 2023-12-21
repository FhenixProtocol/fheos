// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "../../FHE.sol";
import "./utils/Utils.sol";

error TestNotFound(string test);

contract NeTest {
    using Utils for *;

    function ne(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "ne(euint8,euint8)")) {
            if (TFHE.decrypt(TFHE.ne(TFHE.asEuint8(a), TFHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ne(euint16,euint16)")) {
            if (TFHE.decrypt(TFHE.ne(TFHE.asEuint16(a), TFHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ne(euint32,euint32)")) {
            if (TFHE.decrypt(TFHE.ne(TFHE.asEuint32(a), TFHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint8.ne(euint8)")) {
            if (TFHE.decrypt(TFHE.asEuint8(a).ne(TFHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.ne(euint16)")) {
            if (TFHE.decrypt(TFHE.asEuint16(a).ne(TFHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.ne(euint32)")) {
            if (TFHE.decrypt(TFHE.asEuint32(a).ne(TFHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ne(ebool,ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (TFHE.decrypt(TFHE.ne(TFHE.asEbool(aBool), TFHE.asEbool(bBool)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ebool.ne(ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (TFHE.decrypt(TFHE.asEbool(aBool).ne(TFHE.asEbool(bBool)))) {
                return 1;
            }
        } else {
            revert TestNotFound(test);
        }
    }

}