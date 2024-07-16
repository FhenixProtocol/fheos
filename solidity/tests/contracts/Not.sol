// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract NotTest {
    using Utils for *;
    
    function not(string calldata test, uint256 a, int32 securityZone) public pure returns (uint256 output) {
        if (Utils.cmp(test, "not(euint8)")) {
            return FHE.decrypt(FHE.not(FHE.asEuint8(a, securityZone)));
        } else if (Utils.cmp(test, "not(euint16)")) {
            return FHE.decrypt(FHE.not(FHE.asEuint16(a, securityZone)));
        } else if (Utils.cmp(test, "not(euint32)")) {
            return FHE.decrypt(FHE.not(FHE.asEuint32(a, securityZone)));
        } else if (Utils.cmp(test, "not(euint64)")) {
            return FHE.decrypt(FHE.not(FHE.asEuint64(a, securityZone)));
        } else if (Utils.cmp(test, "not(euint128)")) {
            return FHE.decrypt(FHE.not(FHE.asEuint128(a, securityZone)));
        } else if (Utils.cmp(test, "not(ebool)")) {
            bool aBool = true;
            if (a == 0) {
                aBool = false;
            }

            if (FHE.decrypt(FHE.not(FHE.asEbool(aBool, securityZone)))) {
                return 1;
            }

            return 0;
        }
        
        revert TestNotFound(test);
    }
}
