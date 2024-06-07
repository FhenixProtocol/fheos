// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";

contract MulBench {
	private ebool aBool;
	private euint8 a8;
	private euint16 a16;
	private euint32 a32;
	private euint64 a64;

	private ebool bBool;
	private euint8 b8;
	private euint16 b16;
	private euint32 b32;
	private euint64 b64;

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

    function benchMulBool() public view {
        FHE.mul(aBool, bBool);
    }
    function benchMul8() public view {
        FHE.mul(a8, b8);
    }
    function benchMul16() public view {
        FHE.mul(a16, b16);
    }
    function benchMul32() public view {
        FHE.mul(a32, b32);
    }
    function benchMul64() public view {
        FHE.mul(a64, b64);
    }
}
