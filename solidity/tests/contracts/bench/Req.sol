// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
import {ebool, euint8} from "../../../FHE.sol";

contract ReqBench {

    euint8 internal a8;
    euint16 internal a16;
    euint32 internal a32;
    euint64 internal a64;
    euint128 internal a128;
    euint256 internal a256;
  
    function load32(inEuint32 _a) public {
        a32 = FHE.asEuint32(_a);
    }
    
    function benchReq32() public view {
        FHE.req(a32);
    }
}
