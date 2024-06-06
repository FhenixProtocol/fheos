// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";

contract NotBench {
    private euint8 a8;
    private euint16 a16;
    private euint32 a32;
    private euint64 a64;
    private euint128 a128;
    private euint256 a256;
  
    function not(string calldata test, uint256 a) public pure returns (uint256 output) {
        if (Utils.cmp(test, "not(euint8)")) {
            return FHE.decrypt(FHE.not(FHE.asEuint8(a)));
        } else if (Utils.cmp(test, "not(euint16)")) {
            return FHE.decrypt(FHE.not(FHE.asEuint16(a)));
        } else if (Utils.cmp(test, "not(euint32)")) {
            return FHE.decrypt(FHE.not(FHE.asEuint32(a)));
        } else if (Utils.cmp(test, "not(euint64)")) {
            return FHE.decrypt(FHE.not(FHE.asEuint64(a)));
        } else if (Utils.cmp(test, "not(euint128)")) {
            return FHE.decrypt(FHE.not(FHE.asEuint128(a)));
        } else if (Utils.cmp(test, "not(ebool)")) {
            bool aBool = true;
            if (a == 0) {
                aBool = false;
            }

            if (FHE.decrypt(FHE.not(FHE.asEbool(aBool)))) {
                return 1;
            }

            return 0;
        }
        
        revert TestNotFound(test);
    }
}
