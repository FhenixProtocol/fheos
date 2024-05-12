// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract EqTest {
    using Utils for *;
    
    function eq(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "eq(euint8,euint8)")) {
            if (FHE.decrypt(FHE.eq(FHE.asEuint8(a), FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "eq(euint16,euint16)")) {
            if (FHE.decrypt(FHE.eq(FHE.asEuint16(a), FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "eq(euint32,euint32)")) {
            if (FHE.decrypt(FHE.eq(FHE.asEuint32(a), FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "eq(euint64,euint64)")) {
            if (FHE.decrypt(FHE.eq(FHE.asEuint64(a), FHE.asEuint64(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "eq(euint128,euint128)")) {
            if (FHE.decrypt(FHE.eq(FHE.asEuint128(a), FHE.asEuint128(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "eq(euint256,euint256)")) {
            if (FHE.decrypt(FHE.eq(FHE.asEuint256(a), FHE.asEuint256(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "eq(eaddress,eaddress)")) {
            if (FHE.decrypt(FHE.eq(FHE.asEaddress(a), FHE.asEaddress(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint8.eq(euint8)")) {
            if (FHE.asEuint8(a).eq(FHE.asEuint8(b)).decrypt()) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.eq(euint16)")) {
            if (FHE.asEuint16(a).eq(FHE.asEuint16(b)).decrypt()) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.eq(euint32)")) {
            if (FHE.asEuint32(a).eq(FHE.asEuint32(b)).decrypt()) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint64.eq(euint64)")) {
            if (FHE.asEuint64(a).eq(FHE.asEuint64(b)).decrypt()) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "euint128.eq(euint128)")) {
            if (FHE.asEuint128(a).eq(FHE.asEuint128(b)).decrypt()) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "euint256.eq(euint256)")) {
            if (FHE.asEuint256(a).eq(FHE.asEuint256(b)).decrypt()) {
                return 1;
            }
            return 0;
        } else if (Utils.cmp(test, "eaddress.eq(eaddress)")) {
            if (FHE.asEaddress(a).eq(FHE.asEaddress(b)).decrypt()) {
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
            if (FHE.decrypt(FHE.eq(FHE.asEbool(aBool), FHE.asEbool(bBool)))) {
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
            if (FHE.asEbool(aBool).eq(FHE.asEbool(bBool)).decrypt()) {
                return 1;
            }
            return 0;
        }
        revert TestNotFound(test);
    }
}
