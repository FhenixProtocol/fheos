// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {TFHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract MulTest {
    using Utils for *;
function mul(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "mul(euint8,euint8)")) {
            return TFHE.decrypt(TFHE.mul(TFHE.asEuint8(a), TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "mul(euint16,euint16)")) {
            return TFHE.decrypt(TFHE.mul(TFHE.asEuint16(a), TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "mul(euint32,euint32)")) {
            return TFHE.decrypt(TFHE.mul(TFHE.asEuint32(a), TFHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8.mul(euint8)")) {
            return TFHE.decrypt(TFHE.asEuint8(a).mul(TFHE.asEuint8(b)));
        } else if (Utils.cmp(test, "euint16.mul(euint16)")) {
            return TFHE.decrypt(TFHE.asEuint16(a).mul(TFHE.asEuint16(b)));
        } else if (Utils.cmp(test, "euint32.mul(euint32)")) {
            return TFHE.decrypt(TFHE.asEuint32(a).mul(TFHE.asEuint32(b)));
        } else if (Utils.cmp(test, "euint8 * euint8")) {
            return TFHE.decrypt(TFHE.asEuint8(a) * TFHE.asEuint8(b));
        } else if (Utils.cmp(test, "euint16 * euint16")) {
            return TFHE.decrypt(TFHE.asEuint16(a) * TFHE.asEuint16(b));
        } else if (Utils.cmp(test, "euint32 * euint32")) {
            return TFHE.decrypt(TFHE.asEuint32(a) * TFHE.asEuint32(b));
        } else {
            revert TestNotFound(test);
        }
    }
}
