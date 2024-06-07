// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
import {
	ebool, inEbool,
	euint8, inEuint8,
	euint16, inEuint16,
	euint32, inEuint32,
	euint64, inEuint64,
	euint128, inEuint128
} from "../../../FHE.sol";

contract AndBench {
	ebool internal aBool;
	euint8 internal a8;
	euint16 internal a16;
	euint32 internal a32;
	euint64 internal a64;
	euint128 internal a128;

	ebool internal bBool;
	euint8 internal b8;
	euint16 internal b16;
	euint32 internal b32;
	euint64 internal b64;
	euint128 internal b128;

    function loadBool(inEbool _a, inEbool _b) public {
        a32 = FHE.asEbool(_a);
        b32 = FHE.asEbool(_b);
    }
    function load8(inEuint8 _a, inEuint8 _b) public {
        a32 = FHE.asEuint8(_a);
        b32 = FHE.asEuint8(_b);
    }
    function load16(inEuint16 _a, inEuint16 _b) public {
        a32 = FHE.asEuint16(_a);
        b32 = FHE.asEuint16(_b);
    }
    function load32(inEuint32 _a, inEuint32 _b) public {
        a32 = FHE.asEuint32(_a);
        b32 = FHE.asEuint32(_b);
    }
    function load64(inEuint64 _a, inEuint64 _b) public {
        a32 = FHE.asEuint64(_a);
        b32 = FHE.asEuint64(_b);
    }
    function load128(inEuint128 _a, inEuint128 _b) public {
        a32 = FHE.asEuint128(_a);
        b32 = FHE.asEuint128(_b);
    }

    function benchAndBool() public view {
        FHE.and(aBool, bBool);
    }
    function benchAnd8() public view {
        FHE.and(a8, b8);
    }
    function benchAnd16() public view {
        FHE.and(a16, b16);
    }
    function benchAnd32() public view {
        FHE.and(a32, b32);
    }
    function benchAnd64() public view {
        FHE.and(a64, b64);
    }
    function benchAnd128() public view {
        FHE.and(a128, b128);
    }
}
