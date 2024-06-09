// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
import {
    ebool, inEbool,
    euint8, inEuint8,
    euint16, inEuint16,
    euint64, inEuint64,
    euint128, inEuint128,
    euint256, inEuint256,
    eaddress, inEaddress
} from "../../../FHE.sol";

contract AsEuint32Bench {
    ebool internal aBool;
    euint8 internal a8;
    euint16 internal a16;
    euint64 internal a64;
    euint128 internal a128;
    euint256 internal a256;
    eaddress internal aAddress;
    uint256 internal aUint256;
    bytes internal aBytes;

    function loadBool(inEbool calldata _a) public {
        aBool = FHE.asEbool(_a);
    }
    function load8(inEuint8 calldata _a) public {
        a8 = FHE.asEuint8(_a);
    }
    function load16(inEuint16 calldata _a) public {
        a16 = FHE.asEuint16(_a);
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

    function benchCastEboolToEuint32() public view {
        FHE.asEuint32(aBool);
    }
    function benchCastEuint8ToEuint32() public view {
        FHE.asEuint32(a8);
    }
    function benchCastEuint16ToEuint32() public view {
        FHE.asEuint32(a16);
    }
    function benchCastEuint64ToEuint32() public view {
        FHE.asEuint32(a64);
    }
    function benchCastEuint128ToEuint32() public view {
        FHE.asEuint32(a128);
    }
    function benchCastEuint256ToEuint32() public view {
        FHE.asEuint32(a256);
    }
    function benchCastEaddressToEuint32() public view {
        FHE.asEuint32(aAddress);
    }
    function benchCastUint256ToEuint32() public view {
        FHE.asEuint32(aUint256);
    }
    function benchCastBytesToEuint32() public view {
        FHE.asEuint32(aBytes);
    }
}
