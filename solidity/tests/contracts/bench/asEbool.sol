// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
import {
    euint8, inEuint8,
    euint16, inEuint16,
    euint32, inEuint32,
    euint64, inEuint64,
    euint128, inEuint128,
    euint256, inEuint256,
    eaddress, inEaddress,
    inEbool
} from "../../../FHE.sol";

contract AsEboolBench {
    euint8 internal a8;
    euint16 internal a16;
    euint32 internal a32;
    euint64 internal a64;
    euint128 internal a128;
    euint256 internal a256;
    eaddress internal aAddress;
    uint256 internal aUint256;
    inEbool internal aPreEncrypted;

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
    function loadPreEncrypted(inEbool calldata _a) public {
        aPreEncrypted = _a;
    }

    function benchCastEuint8ToEbool() public view {
        FHE.asEbool(a8);
    }
    function benchCastEuint16ToEbool() public view {
        FHE.asEbool(a16);
    }
    function benchCastEuint32ToEbool() public view {
        FHE.asEbool(a32);
    }
    function benchCastEuint64ToEbool() public view {
        FHE.asEbool(a64);
    }
    function benchCastEuint128ToEbool() public view {
        FHE.asEbool(a128);
    }
    function benchCastEuint256ToEbool() public view {
        FHE.asEbool(a256);
    }
    function benchCastEaddressToEbool() public view {
        FHE.asEbool(aAddress);
    }
    function benchCastUint256ToEbool() public view {
        FHE.asEbool(aUint256);
    }
    function benchCastPreEncryptedToEbool() public view {
        FHE.asEbool(aPreEncrypted);
    }
}
