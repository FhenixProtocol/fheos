// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract MaxTest {
    using Utils for *;
    
    function max(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "max(euint8,euint8)")) {
            return FHE.decrypt(FHE.max(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "max(euint16,euint16)")) {
            return FHE.decrypt(FHE.max(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "max(euint32,euint32)")) {
            return FHE.decrypt(FHE.max(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "max(euint64,euint64)")) {
            return FHE.decrypt(FHE.max(FHE.asEuint64(a), FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "max(euint128,euint128)")) {
            return FHE.decrypt(FHE.max(FHE.asEuint128(a), FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "max(euint256,euint256)")) {
            return FHE.decrypt(FHE.max(FHE.asEuint256(a), FHE.asEuint256(b))); 
        } else if (Utils.cmp(test, "euint8.max(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).max(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.max(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).max(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.max(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).max(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint64.max(euint64)")) {
            return FHE.decrypt(FHE.asEuint64(a).max(FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "euint128.max(euint128)")) {
            return FHE.decrypt(FHE.asEuint128(a).max(FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "euint256.max(euint256)")) {
            return FHE.decrypt(FHE.asEuint256(a).max(FHE.asEuint256(b)));
        }
    
        revert TestNotFound(test);
    }
}
