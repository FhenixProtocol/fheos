// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";

contract ReqBench {

    private euint8 a8;
    private euint16 a16;
    private euint32 a32;
    private euint64 a64;
    private euint128 a128;
    private euint256 a256;
  
    function load32(inEuint32 _a) public {
        a32 = FHE.asEuint32(_a);
    }
    
    function benchReq32() public view {
        FHE.req(a32);
    }
}
