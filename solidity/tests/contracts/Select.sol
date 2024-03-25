// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {ebool, euint8} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract SelectTest {
    using Utils for *;
    
    function select(string calldata test, bool c, uint256 a, uint256 b) public pure returns (uint256 output) {
        ebool condition = FHE.asEbool(c);
        if (Utils.cmp(test, "select: euint8")) {
            return FHE.decrypt(FHE.select(condition, FHE.asEuint8(a), FHE.asEuint8(b)));
        } else if (Utils.cmp(test, "select: euint16")) {
            return FHE.decrypt(FHE.select(condition, FHE.asEuint16(a), FHE.asEuint16(b)));
        } else if (Utils.cmp(test, "select: euint32")) {
            return FHE.decrypt(FHE.select(condition, FHE.asEuint32(a), FHE.asEuint32(b)));
        } else if (Utils.cmp(test, "select: euint64")) {
            return FHE.decrypt(FHE.select(condition, FHE.asEuint64(a), FHE.asEuint64(b)));
        } else if (Utils.cmp(test, "select: euint128")) {
            return FHE.decrypt(FHE.select(condition, FHE.asEuint128(a), FHE.asEuint128(b)));
        } else if (Utils.cmp(test, "select: euint256")) {
            return FHE.decrypt(FHE.select(condition, FHE.asEuint256(a), FHE.asEuint256(b)));
        } else if (Utils.cmp(test, "select: ebool")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }

            if(FHE.decrypt(FHE.select(condition, FHE.asEbool(aBool), FHE.asEbool(bBool)))) {
                return 1;
            }
            return 0;
        } 
        
        revert TestNotFound(test);
    }
}
