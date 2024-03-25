// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract MulTest {
    using Utils for *;
    
    function mul(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "mul(euint8,euint8)")) {
            return FHE.decrypt(FHE.mul(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "mul(euint16,euint16)")) {
            return FHE.decrypt(FHE.mul(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "mul(euint32,euint32)")) {
            return FHE.decrypt(FHE.mul(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "mul(euint64,euint64)")) {
            return FHE.decrypt(FHE.mul(FHE.asEuint64(a), FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "mul(euint128,euint128)")) {
            return FHE.decrypt(FHE.mul(FHE.asEuint128(a), FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "mul(euint256,euint256)")) {
            return FHE.decrypt(FHE.mul(FHE.asEuint256(a), FHE.asEuint256(b))); 
        } else if (Utils.cmp(test, "euint8.mul(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).mul(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.mul(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).mul(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.mul(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).mul(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint64.mul(euint64)")) {
            return FHE.decrypt(FHE.asEuint64(a).mul(FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "euint128.mul(euint128)")) {
            return FHE.decrypt(FHE.asEuint128(a).mul(FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "euint256.mul(euint256)")) {
            return FHE.decrypt(FHE.asEuint256(a).mul(FHE.asEuint256(b)));
        } else if (Utils.cmp(test, "euint8 * euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) * FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 * euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) * FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 * euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) * FHE.asEuint32(b));
        } else if (Utils.cmp(test, "euint64 * euint64")) {
            return FHE.decrypt(FHE.asEuint64(a) * FHE.asEuint64(b));
        } else if (Utils.cmp(test, "euint128 * euint128")) {
            return FHE.decrypt(FHE.asEuint128(a) * FHE.asEuint128(b));
        } else if (Utils.cmp(test, "euint256 * euint256")) {
            return FHE.decrypt(FHE.asEuint256(a) * FHE.asEuint256(b));
        }
    
        revert TestNotFound(test);
    }
}
