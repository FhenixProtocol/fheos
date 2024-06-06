// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";

contract DivBench {
    private euint8 a8;
    private euint16 a16;
    private euint32 a32;
    private euint64 a64;
    private euint128 a128;
    private euint256 a256;
  
    function div(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "div(euint8,euint8)")) {
            return FHE.decrypt(FHE.div(FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "div(euint16,euint16)")) {
            return FHE.decrypt(FHE.div(FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "div(euint32,euint32)")) {
            return FHE.decrypt(FHE.div(FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.div(euint8)")) {
            return FHE.decrypt(FHE.asEuint8(a).div(FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.div(euint16)")) {
            return FHE.decrypt(FHE.asEuint16(a).div(FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.div(euint32)")) {
            return FHE.decrypt(FHE.asEuint32(a).div(FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8 / euint8")) {
            return FHE.decrypt(FHE.asEuint8(a) / FHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 / euint16")) {
            return FHE.decrypt(FHE.asEuint16(a) / FHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 / euint32")) {
            return FHE.decrypt(FHE.asEuint32(a) / FHE.asEuint32(b));
        }
    
        revert TestNotFound(test);
    }
}
