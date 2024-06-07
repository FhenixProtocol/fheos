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

contract SelectBench {
	ebool internal control;

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

    function loadBool(inEbool calldata _control, inEbool calldata _a, inEbool calldata _b) public {
        control = FHE.asEbool(_control);
        aBool = FHE.asEbool(_a);
        bBool = FHE.asEbool(_b);
    }
    function load8(inEbool calldata _control, inEuint8 calldata _a, inEuint8 calldata _b) public {
        control = FHE.asEbool(_control);
        a8 = FHE.asEuint8(_a);
        b8 = FHE.asEuint8(_b);
    }
    function load16(inEbool calldata _control, inEuint16 calldata _a, inEuint16 calldata _b) public {
        control = FHE.asEbool(_control);
        a16 = FHE.asEuint16(_a);
        b16 = FHE.asEuint16(_b);
    }
    function load32(inEbool calldata _control, inEuint32 calldata _a, inEuint32 calldata _b) public {
        control = FHE.asEbool(_control);
        a32 = FHE.asEuint32(_a);
        b32 = FHE.asEuint32(_b);
    }
    function load64(inEbool calldata _control, inEuint64 calldata _a, inEuint64 calldata _b) public {
        control = FHE.asEbool(_control);
        a64 = FHE.asEuint64(_a);
        b64 = FHE.asEuint64(_b);
    }
    function load128(inEbool calldata _control, inEuint128 calldata _a, inEuint128 calldata _b) public {
        control = FHE.asEbool(_control);
        a128 = FHE.asEuint128(_a);
        b128 = FHE.asEuint128(_b);
    }
    function load256(inEbool calldata _control, inEuint256 calldata _a, inEuint256 calldata _b) public {
        control = FHE.asEbool(_control);
        a256 = FHE.asEuint256(_a);
        b256 = FHE.asEuint256(_b);
    }
    function loadAddress(inEbool calldata _control, inEaddress calldata _a, inEaddress calldata _b) public {
        control = FHE.asEbool(_control);
        aAddress = FHE.asEaddress(_a);
        bAddress = FHE.asEaddress(_b);
    }

    function benchSelectBool() public view {
        FHE.select(control, aBool, bBool);
    }
    function benchSelect8() public view {
        FHE.select(control, a8, b8);
    }
    function benchSelect16() public view {
        FHE.select(control, a16, b16);
    }
    function benchSelect32() public view {
        FHE.select(control, a32, b32);
    }
    function benchSelect64() public view {
        FHE.select(control, a64, b64);
    }
    function benchSelect128() public view {
        FHE.select(control, a128, b128);
    }
    function benchSelect256() public view {
        FHE.select(control, a256, b256);
    }
    function benchSelectAddress() public view {
        FHE.select(control, aAddress, bAddress);
    }
}
