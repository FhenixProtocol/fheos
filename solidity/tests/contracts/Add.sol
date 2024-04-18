// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AddTest {
    using Utils for *;
    
    function add(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "add(euint8,euint8)")) {
            return FHE.decrypt(FHE.add(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "add(euint16,euint16)")) {
            return FHE.decrypt(FHE.add(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "add(euint32,euint32)")) {
            return FHE.decrypt(FHE.add(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "add(euint64,euint64)")) {
            return FHE.decrypt(FHE.add(FHE.asEuint64(a), FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "add(euint128,euint128)")) {
            return FHE.decrypt(FHE.add(FHE.asEuint128(a), FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "euint8.add(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).add(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.add(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).add(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.add(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).add(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint64.add(euint64)")) {
            return FHE.decrypt(FHE.asEuint64(a).add(FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "euint128.add(euint128)")) {
            return FHE.decrypt(FHE.asEuint128(a).add(FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "euint8 + euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) + FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 + euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) + FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 + euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) + FHE.asEuint32(b));
        }else if (Utils.cmp(test, "euint64 + euint64")) {
            return FHE.decrypt(FHE.asEuint64(a) + FHE.asEuint64(b));
        }else if (Utils.cmp(test, "euint128 + euint128")) {
            return FHE.decrypt(FHE.asEuint128(a) + FHE.asEuint128(b));
        }
    
        revert TestNotFound(test);
    }
}
