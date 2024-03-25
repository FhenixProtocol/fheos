// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract DivTest {
    using Utils for *;
    
    function div(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "div(euint8,euint8)")) {
            return FHE.decrypt(FHE.div(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "div(euint16,euint16)")) {
            return FHE.decrypt(FHE.div(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "div(euint32,euint32)")) {
            return FHE.decrypt(FHE.div(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "div(euint64,euint64)")) {
            return FHE.decrypt(FHE.div(FHE.asEuint64(a), FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "div(euint128,euint128)")) {
            return FHE.decrypt(FHE.div(FHE.asEuint128(a), FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "div(euint256,euint256)")) {
            return FHE.decrypt(FHE.div(FHE.asEuint256(a), FHE.asEuint256(b))); 
        } else if (Utils.cmp(test, "euint8.div(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).div(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.div(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).div(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.div(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).div(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint64.div(euint64)")) {
            return FHE.decrypt(FHE.asEuint64(a).div(FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "euint128.div(euint128)")) {
            return FHE.decrypt(FHE.asEuint128(a).div(FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "euint256.div(euint256)")) {
            return FHE.decrypt(FHE.asEuint256(a).div(FHE.asEuint256(b)));
        } else if (Utils.cmp(test, "euint8 / euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) / FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 / euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) / FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 / euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) / FHE.asEuint32(b));
        } else if (Utils.cmp(test, "euint64 / euint64")) {
            return FHE.decrypt(FHE.asEuint64(a) / FHE.asEuint64(b));
        } else if (Utils.cmp(test, "euint128 / euint128")) {
            return FHE.decrypt(FHE.asEuint128(a) / FHE.asEuint128(b));
        } else if (Utils.cmp(test, "euint256 / euint256")) {
            return FHE.decrypt(FHE.asEuint256(a) / FHE.asEuint256(b));
        }
    
        revert TestNotFound(test);
    }
}
