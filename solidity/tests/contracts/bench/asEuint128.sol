// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
import {
    ebool, inEbool,
    euint8, inEuint8,
    euint16, inEuint16,
    euint32, inEuint32,
    euint64, inEuint64,
    euint256, inEuint256,
    eaddress, inEaddress,
    inEuint128
} from "../../../FHE.sol";

contract AsEuint128Bench {
    ebool internal aBool;
    euint8 internal a8;
    euint16 internal a16;
    euint32 internal a32;
    euint64 internal a64;
    euint256 internal a256;
    eaddress internal aAddress;
    uint256 internal aUint256;
    inEuint128 internal aPreEncrypted;

    function loadBool(inEbool calldata _a) public {
        aBool = FHE.asEbool(_a);
    }
    function load8(inEuint8 calldata _a) public {
        a8 = FHE.asEuint8(_a);
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
    function load256(inEuint256 calldata _a) public {
        a256 = FHE.asEuint256(_a);
    }
    function loadAddress(inEaddress calldata _a) public {
        aAddress = FHE.asEaddress(_a);
    }
    function loadUint256(uint256 _a) public {
        aUint256 = _a;
    }
    function loadPreEncrypted(inEuint128 calldata _a) public {
        aPreEncrypted = _a;
    }

    function benchCastEboolToEuint128() public view {
        FHE.asEuint128(aBool);
    }
    function benchCastEuint8ToEuint128() public view {
        FHE.asEuint128(a8);
    }
    function benchCastEuint16ToEuint128() public view {
        FHE.asEuint128(a16);
    }
    function benchCastEuint32ToEuint128() public view {
        FHE.asEuint128(a32);
    }
    function benchCastEuint64ToEuint128() public view {
        FHE.asEuint128(a64);
    }
    function benchCastEuint256ToEuint128() public view {
        FHE.asEuint128(a256);
    }
    function benchCastEaddressToEuint128() public view {
        FHE.asEuint128(aAddress);
    }
    function benchCastUint256ToEuint128() public view {
        FHE.asEuint128(aUint256);
    }
    function benchCastPreEncryptedToEuint128() public view {
        FHE.asEuint128(aPreEncrypted);
    }
}
