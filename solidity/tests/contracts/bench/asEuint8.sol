// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
import {
	ebool, inEbool,
	euint16, inEuint16,
	euint32, inEuint32,
	euint64, inEuint64,
	euint128, inEuint128,
	euint256, inEuint256,
	eaddress, inEaddress
} from "../../../FHE.sol";

contract AsEuint8Bench {
	ebool internal aBool;
	euint16 internal a16;
	euint32 internal a32;
	euint64 internal a64;
	euint128 internal a128;
	euint256 internal a256;
	eaddress internal aAddress;
	uint256 internal aUint256;
	bytes internal aBytes;

	function loadBool(inEbool calldata _a) public {
        aBool = FHE.asEbool(_a);
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
	function loadAddress(inEaddress calldata _a) public {
        aAddress = FHE.asEaddress(_a);
    }
	function loadUint256(uint256 _a) public {
        aUint256 = _a;
    }
	function loadBytes(bytes memory _a) public {
        aBytes = _a;
    }

	function benchCastEboolToEuint8() public view {
        FHE.asEuint8(aBool);
    }
	function benchCastEuint16ToEuint8() public view {
        FHE.asEuint8(a16);
    }
	function benchCastEuint32ToEuint8() public view {
        FHE.asEuint8(a32);
    }
	function benchCastEuint64ToEuint8() public view {
        FHE.asEuint8(a64);
    }
	function benchCastEuint128ToEuint8() public view {
        FHE.asEuint8(a128);
    }
	function benchCastEuint256ToEuint8() public view {
        FHE.asEuint8(a256);
    }
	function benchCastEaddressToEuint8() public view {
        FHE.asEuint8(aAddress);
    }
	function benchCastUint256ToEuint8() public view {
        FHE.asEuint8(aUint256);
    }
	function benchCastBytesToEuint8() public view {
        FHE.asEuint8(aBytes);
    }
}
