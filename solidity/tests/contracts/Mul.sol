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
        } else if (Utils.cmp(test, "euint8.mul(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).mul(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.mul(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).mul(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.mul(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).mul(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8 * euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) * FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 * euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) * FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 * euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) * FHE.asEuint32(b));
        }
    
        revert TestNotFound(test);
    }
}
