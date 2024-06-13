// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
import {
    ebool, inEbool,
    euint8, inEuint8,
    euint16, inEuint16,
    euint32, inEuint32,
    euint64, inEuint64,
    euint128, inEuint128,
    euint256, inEuint256
} from "../../../FHE.sol";

contract ReqBench {
    ebool internal aBool;
    euint8 internal a8;
    euint16 internal a16;
    euint32 internal a32;
    euint64 internal a64;
    euint128 internal a128;
    euint256 internal a256;

    function loadBool(inEbool calldata _a) public {
        aBool = FHE.asEbool(_a);
    }
    function load8(inEuint8 calldata _a) public {
        a8 = FHE.asEuint8(_a);
    }
    function load16(inEuint16 calldata _a) public {
        a16 = FHE.asEuint16(_a);
    }
    function load32(inEuint32 calldata _a) public {
        a32 = FHE.asEuint32(_a);
    }
    function load64(inEuint64 calldata _a) public {
        a64 = FHE.asEuint64(_a);
    }
    function load128(inEuint128 calldata _a) public {
        a128 = FHE.asEuint128(_a);
    }
    function load256(inEuint256 calldata _a) public {
        a256 = FHE.asEuint256(_a);
    }

    function benchReqBool() public view {
        FHE.req(aBool);
    }
    function benchReq8() public view {
        FHE.req(a8);
    }
    function benchReq16() public view {
        FHE.req(a16);
    }
    function benchReq32() public view {
        FHE.req(a32);
    }
    function benchReq64() public view {
        FHE.req(a64);
    }
    function benchReq128() public view {
        FHE.req(a128);
    }
    function benchReq256() public view {
        FHE.req(a256);
    }
}
