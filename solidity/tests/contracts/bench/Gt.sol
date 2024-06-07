// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";

contract GtBench {
	private ebool aBool;
	private euint8 a8;
	private euint16 a16;
	private euint32 a32;
	private euint64 a64;
	private euint128 a128;

	private ebool bBool;
	private euint8 b8;
	private euint16 b16;
	private euint32 b32;
	private euint64 b64;
	private euint128 b128;

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

    function benchGtBool() public view {
        FHE.gt(aBool, bBool);
    }
    function benchGt8() public view {
        FHE.gt(a8, b8);
    }
    function benchGt16() public view {
        FHE.gt(a16, b16);
    }
    function benchGt32() public view {
        FHE.gt(a32, b32);
    }
    function benchGt64() public view {
        FHE.gt(a64, b64);
    }
    function benchGt128() public view {
        FHE.gt(a128, b128);
    }
}
