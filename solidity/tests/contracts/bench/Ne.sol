// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";

contract NeBench {
	private ebool aBool;
	private euint8 a8;
	private euint16 a16;
	private euint32 a32;
	private euint64 a64;
	private euint128 a128;
	private euint256 a256;
	private eaddress aAddress;

	private ebool bBool;
	private euint8 b8;
	private euint16 b16;
	private euint32 b32;
	private euint64 b64;
	private euint128 b128;
	private euint256 b256;
	private eaddress bAddress;

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

    function benchNeBool() public view {
        FHE.ne(aBool, bBool);
    }
    function benchNe8() public view {
        FHE.ne(a8, b8);
    }
    function benchNe16() public view {
        FHE.ne(a16, b16);
    }
    function benchNe32() public view {
        FHE.ne(a32, b32);
    }
    function benchNe64() public view {
        FHE.ne(a64, b64);
    }
    function benchNe128() public view {
        FHE.ne(a128, b128);
    }
    function benchNe256() public view {
        FHE.ne(a256, b256);
    }
    function benchNeAddress() public view {
        FHE.ne(aAddress, bAddress);
    }
}
