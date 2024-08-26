// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE, euint8, euint16, euint32, euint64, euint128, euint256, ebool} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract ReqTest {
    using Utils for *;
    
    function req(string calldata test, uint256 a) public {
        if (Utils.cmp(test, "req(euint8)")) {
            FHE.req(FHE.asEuint8(a));
        } else if (Utils.cmp(test, "req(euint16)")) {
            FHE.req(FHE.asEuint16(a));
        } else if (Utils.cmp(test, "req(euint32)")) {
            FHE.req(FHE.asEuint32(a));
        } else if (Utils.cmp(test, "req(euint64)")) {
            FHE.req(FHE.asEuint64(a));
        } else if (Utils.cmp(test, "req(euint128)")) {
            FHE.req(FHE.asEuint128(a));
        } else if (Utils.cmp(test, "req(euint256)")) {
            FHE.req(FHE.asEuint256(a));
        } else if (Utils.cmp(test, "req(ebool)")) {
            bool b = true;
            if (a == 0) {
                b = false;
            }
            FHE.req(FHE.asEbool(b));
        } else {
            revert TestNotFound(test);
        }
    }
}
