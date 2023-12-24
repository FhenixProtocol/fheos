// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {TFHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract MinTest {
    using Utils for *;
    
    function min(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "min(euint8,euint8)")) {
            return TFHE.decrypt(TFHE.min(TFHE.asEuint8(a), TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "min(euint16,euint16)")) {
            return TFHE.decrypt(TFHE.min(TFHE.asEuint16(a), TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "min(euint32,euint32)")) {
            return TFHE.decrypt(TFHE.min(TFHE.asEuint32(a), TFHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.min(euint8)")) {
            return TFHE.decrypt(TFHE.asEuint8(a).min(TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.min(euint16)")) {
            return TFHE.decrypt(TFHE.asEuint16(a).min(TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.min(euint32)")) {
            return TFHE.decrypt(TFHE.asEuint32(a).min(TFHE.asEuint32(b)));
        } else {
            revert TestNotFound(test);
        }
    }
}
