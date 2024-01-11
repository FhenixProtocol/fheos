// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract RemTest {
    using Utils for *;
    
    function rem(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "rem(euint8,euint8)")) {
            return FHE.decrypt(FHE.rem(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "rem(euint16,euint16)")) {
            return FHE.decrypt(FHE.rem(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "rem(euint32,euint32)")) {
            return FHE.decrypt(FHE.rem(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.rem(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).rem(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.rem(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).rem(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.rem(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).rem(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8 % euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) % FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 % euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) % FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 % euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) % FHE.asEuint32(b));
        }
    
        revert TestNotFound(test);
    }
}
