// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";

contract MinBench {
    private euint8 a8;
    private euint16 a16;
    private euint32 a32;
    private euint64 a64;
    private euint128 a128;
    private euint256 a256;
  
    function min(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "min(euint8,euint8)")) {
            return FHE.decrypt(FHE.min(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "min(euint16,euint16)")) {
            return FHE.decrypt(FHE.min(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "min(euint32,euint32)")) {
            return FHE.decrypt(FHE.min(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "min(euint64,euint64)")) {
            return FHE.decrypt(FHE.min(FHE.asEuint64(a), FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "min(euint128,euint128)")) {
            return FHE.decrypt(FHE.min(FHE.asEuint128(a), FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "euint8.min(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).min(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.min(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).min(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.min(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).min(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint64.min(euint64)")) {
            return FHE.decrypt(FHE.asEuint64(a).min(FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "euint128.min(euint128)")) {
            return FHE.decrypt(FHE.asEuint128(a).min(FHE.asEuint128(b)));
        }
    
        revert TestNotFound(test);
    }
}
