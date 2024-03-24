// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract NeTest {
    using Utils for *;
    
    function ne(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "ne(euint8,euint8)")) {
            if (FHE.decrypt(FHE.ne(FHE.asEuint8(a), FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ne(euint16,euint16)")) {
            if (FHE.decrypt(FHE.ne(FHE.asEuint16(a), FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ne(euint32,euint32)")) {
            if (FHE.decrypt(FHE.ne(FHE.asEuint32(a), FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ne(euint64,euint64)")) {
            if (FHE.decrypt(FHE.ne(FHE.asEuint64(a), FHE.asEuint64(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ne(euint128,euint128)")) {
            if (FHE.decrypt(FHE.ne(FHE.asEuint128(a), FHE.asEuint128(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ne(euint256,euint256)")) {
            if (FHE.decrypt(FHE.ne(FHE.asEuint256(a), FHE.asEuint256(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint8.ne(euint8)")) {
            if (FHE.asEuint8(a).ne(FHE.asEuint8(b)).decrypt()) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.ne(euint16)")) {
            if (FHE.asEuint16(a).ne(FHE.asEuint16(b)).decrypt()) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.ne(euint32)")) {
            if (FHE.asEuint32(a).ne(FHE.asEuint32(b)).decrypt()) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint64.ne(euint64)")) {
            if (FHE.asEuint64(a).ne(FHE.asEuint64(b)).decrypt()) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "euint128.ne(euint128)")) {
            if (FHE.asEuint128(a).ne(FHE.asEuint128(b)).decrypt()) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "euint256.ne(euint256)")) {
            if (FHE.asEuint256(a).ne(FHE.asEuint256(b)).decrypt()) {
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
            if (FHE.decrypt(FHE.ne(FHE.asEbool(aBool), FHE.asEbool(bBool)))) {
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
            if (FHE.asEbool(aBool).ne(FHE.asEbool(bBool)).decrypt()) {
                return 1;
            }
            return 0;
        }
        revert TestNotFound(test);
    }
}
