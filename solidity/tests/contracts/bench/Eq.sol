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
	euint256, inEuint256,
	eaddress, inEaddress
} from "../../../FHE.sol";

contract EqBench {
	ebool internal aBool;
	euint8 internal a8;
	euint16 internal a16;
	euint32 internal a32;
	euint64 internal a64;
	euint128 internal a128;
	euint256 internal a256;
	eaddress internal aAddress;

	ebool internal bBool;
	euint8 internal b8;
	euint16 internal b16;
	euint32 internal b32;
	euint64 internal b64;
	euint128 internal b128;
	euint256 internal b256;
	eaddress internal bAddress;

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
    function load256(inEuint256 _a, inEuint256 _b) public {
        a32 = FHE.asEuint256(_a);
        b32 = FHE.asEuint256(_b);
    }
    function loadAddress(inEaddress _a, inEaddress _b) public {
        a32 = FHE.asEaddress(_a);
        b32 = FHE.asEaddress(_b);
    }

    function benchEqBool() public view {
        FHE.eq(aBool, bBool);
    }
    function benchEq8() public view {
        FHE.eq(a8, b8);
    }
    function benchEq16() public view {
        FHE.eq(a16, b16);
    }
    function benchEq32() public view {
        FHE.eq(a32, b32);
    }
    function benchEq64() public view {
        FHE.eq(a64, b64);
    }
    function benchEq128() public view {
        FHE.eq(a128, b128);
    }
    function benchEq256() public view {
        FHE.eq(a256, b256);
    }
    function benchEqAddress() public view {
        FHE.eq(aAddress, bAddress);
    }
}
