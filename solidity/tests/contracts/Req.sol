// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {TFHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract ReqTest {
    using Utils for *;
    
    function req(string calldata test, uint256 a) public {
        if (Utils.cmp(test, "req(euint8)")) {
            TFHE.req(TFHE.asEuint8(a));
        } else if (Utils.cmp(test, "req(euint16)")) {
            TFHE.req(TFHE.asEuint16(a));
        } else if (Utils.cmp(test, "req(euint32)")) {
            TFHE.req(TFHE.asEuint32(a));
        } else if (Utils.cmp(test, "req(ebool)")) {
            bool b = true;
            if (a == 0) {
                b = false;
            }
            TFHE.req(TFHE.asEbool(b));
        }
        
        revert TestNotFound(test);
    }
}
