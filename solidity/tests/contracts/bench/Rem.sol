// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
import {
	ebool, inEbool,
	euint8, inEuint8,
	euint16, inEuint16,
	euint32, inEuint32
} from "../../../FHE.sol";

contract RemBench {
	ebool internal aBool;
	euint8 internal a8;
	euint16 internal a16;
	euint32 internal a32;

	ebool internal bBool;
	euint8 internal b8;
	euint16 internal b16;
	euint32 internal b32;

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

    function benchRemBool() public view {
        FHE.rem(aBool, bBool);
    }
    function benchRem8() public view {
        FHE.rem(a8, b8);
    }
    function benchRem16() public view {
        FHE.rem(a16, b16);
    }
    function benchRem32() public view {
        FHE.rem(a32, b32);
    }
}
