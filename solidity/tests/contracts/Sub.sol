// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract SubTest {
    using Utils for *;
    
    function sub(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "sub(euint8,euint8)")) {
            return FHE.decrypt(FHE.sub(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "sub(euint16,euint16)")) {
            return FHE.decrypt(FHE.sub(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "sub(euint32,euint32)")) {
            return FHE.decrypt(FHE.sub(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "sub(euint64,euint64)")) {
            return FHE.decrypt(FHE.sub(FHE.asEuint64(a), FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "sub(euint128,euint128)")) {
            return FHE.decrypt(FHE.sub(FHE.asEuint128(a), FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "sub(euint256,euint256)")) {
            return FHE.decrypt(FHE.sub(FHE.asEuint256(a), FHE.asEuint256(b))); 
        } else if (Utils.cmp(test, "euint8.sub(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).sub(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.sub(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).sub(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.sub(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).sub(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint64.sub(euint64)")) {
            return FHE.decrypt(FHE.asEuint64(a).sub(FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "euint128.sub(euint128)")) {
            return FHE.decrypt(FHE.asEuint128(a).sub(FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "euint256.sub(euint256)")) {
            return FHE.decrypt(FHE.asEuint256(a).sub(FHE.asEuint256(b)));
        } else if (Utils.cmp(test, "euint8 - euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) - FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 - euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) - FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 - euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) - FHE.asEuint32(b));
        } else if (Utils.cmp(test, "euint64 - euint64")) {
            return FHE.decrypt(FHE.asEuint64(a) - FHE.asEuint64(b));
        } else if (Utils.cmp(test, "euint128 - euint128")) {
            return FHE.decrypt(FHE.asEuint128(a) - FHE.asEuint128(b));
        } else if (Utils.cmp(test, "euint256 - euint256")) {
            return FHE.decrypt(FHE.asEuint256(a) - FHE.asEuint256(b));
        }
    
        revert TestNotFound(test);
    }
}
