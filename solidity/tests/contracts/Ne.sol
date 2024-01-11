// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract NeTest {
    using Utils for *;
    
    function ne(string calldata test, uint256 a, uint256 b) public pure returns (uint256 output) {
        if (Utils.cmp(test, "ne(euint8,euint8)")) {
            if (FHE.decrypt(FHE.ne(FHE.asEuint8(a), FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ne(euint16,euint16)")) {
            if (FHE.decrypt(FHE.ne(FHE.asEuint16(a), FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ne(euint32,euint32)")) {
            if (FHE.decrypt(FHE.ne(FHE.asEuint32(a), FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint8.ne(euint8)")) {
            if (FHE.decrypt(FHE.asEuint8(a).ne(FHE.asEuint8(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint16.ne(euint16)")) {
            if (FHE.decrypt(FHE.asEuint16(a).ne(FHE.asEuint16(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "euint32.ne(euint32)")) {
            if (FHE.decrypt(FHE.asEuint32(a).ne(FHE.asEuint32(b)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ne(ebool,ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.ne(FHE.asEbool(aBool), FHE.asEbool(bBool)))) {
                return 1;
            }

            return 0;
        } else if (Utils.cmp(test, "ebool.ne(ebool)")) {
            bool aBool = true;
            bool bBool = true;
            if (a == 0) {
                aBool = false;
            }
            if (b == 0) {
                bBool = false;
            }
            if (FHE.decrypt(FHE.asEbool(aBool).ne(FHE.asEbool(bBool)))) {
                return 1;
            }
            return 0;
        }
        revert TestNotFound(test);
    }
}
