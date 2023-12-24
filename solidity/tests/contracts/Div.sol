// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {TFHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract DivTest {
    using Utils for *;
    
    function div(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "div(euint8,euint8)")) {
            return TFHE.decrypt(TFHE.div(TFHE.asEuint8(a), TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "div(euint16,euint16)")) {
            return TFHE.decrypt(TFHE.div(TFHE.asEuint16(a), TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "div(euint32,euint32)")) {
            return TFHE.decrypt(TFHE.div(TFHE.asEuint32(a), TFHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.div(euint8)")) {
            return TFHE.decrypt(TFHE.asEuint8(a).div(TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.div(euint16)")) {
            return TFHE.decrypt(TFHE.asEuint16(a).div(TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.div(euint32)")) {
            return TFHE.decrypt(TFHE.asEuint32(a).div(TFHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8 / euint8")) {
            return TFHE.decrypt(TFHE.asEuint8(a) / TFHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 / euint16")) {
            return TFHE.decrypt(TFHE.asEuint16(a) / TFHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 / euint32")) {
            return TFHE.decrypt(TFHE.asEuint32(a) / TFHE.asEuint32(b));
        }
    
        revert TestNotFound(test);
    }
}
