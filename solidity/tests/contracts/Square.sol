// This file is auto-generated by solgen/templates/testContracts.ts
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE, euint8, euint16, euint32, euint64, euint128, euint256, ebool} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract SquareTest {
    using Utils for *;
    
    function square(string calldata test, uint256 a, int32 securityZone) public pure returns (uint256 output) {
        if (Utils.cmp(test, "square(euint8)")) {
            return FHE.decrypt(FHE.square(FHE.asEuint8(a, securityZone)));
        } else if (Utils.cmp(test, "square(euint16)")) {
            return FHE.decrypt(FHE.square(FHE.asEuint16(a, securityZone)));
        } else if (Utils.cmp(test, "square(euint32)")) {
            return FHE.decrypt(FHE.square(FHE.asEuint32(a, securityZone)));
        } else if (Utils.cmp(test, "square(euint64)")) {
            return FHE.decrypt(FHE.square(FHE.asEuint64(a, securityZone)));
        } else if (Utils.cmp(test, "euint8.square()")) {
            return FHE.decrypt(FHE.asEuint8(a, securityZone).square());
        } else if (Utils.cmp(test, "euint16.square()")) {
            return FHE.decrypt(FHE.asEuint16(a, securityZone).square());
        } else if (Utils.cmp(test, "euint32.square()")) {
            return FHE.decrypt(FHE.asEuint32(a, securityZone).square());
        } else if (Utils.cmp(test, "euint64.square()")) {
            return FHE.decrypt(FHE.asEuint64(a, securityZone).square());
        }
        
        revert TestNotFound(test);
    }
}
