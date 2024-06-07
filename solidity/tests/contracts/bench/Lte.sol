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

contract LteBench {
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

    function loadBool(inEbool calldata _a, inEbool calldata _b) public {
        aBool = FHE.asEbool(_a);
        bBool = FHE.asEbool(_b);
    }
    function load8(inEuint8 calldata _a, inEuint8 calldata _b) public {
        a8 = FHE.asEuint8(_a);
        b8 = FHE.asEuint8(_b);
    }
    function load16(inEuint16 calldata _a, inEuint16 calldata _b) public {
        a16 = FHE.asEuint16(_a);
        b16 = FHE.asEuint16(_b);
    }
    function load32(inEuint32 calldata _a, inEuint32 calldata _b) public {
        a32 = FHE.asEuint32(_a);
        b32 = FHE.asEuint32(_b);
    }
    function load64(inEuint64 calldata _a, inEuint64 calldata _b) public {
        a64 = FHE.asEuint64(_a);
        b64 = FHE.asEuint64(_b);
    }
    function load128(inEuint128 calldata _a, inEuint128 calldata _b) public {
        a128 = FHE.asEuint128(_a);
        b128 = FHE.asEuint128(_b);
    }

    function benchLteBool() public view {
        FHE.lte(aBool, bBool);
    }
    function benchLte8() public view {
        FHE.lte(a8, b8);
    }
    function benchLte16() public view {
        FHE.lte(a16, b16);
    }
    function benchLte32() public view {
        FHE.lte(a32, b32);
    }
    function benchLte64() public view {
        FHE.lte(a64, b64);
    }
    function benchLte128() public view {
        FHE.lte(a128, b128);
    }
}
