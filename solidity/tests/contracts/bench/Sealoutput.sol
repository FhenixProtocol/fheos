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

contract SealoutputBench {
	bytes32 internal pubkey;
	ebool internal aBool;
	euint8 internal a8;
	euint16 internal a16;
	euint32 internal a32;
	euint64 internal a64;
	euint128 internal a128;
	euint256 internal a256;
	eaddress internal aAddress;

    function loadBool(inEbool calldata _a, bytes32 _pubkey) public {
        aBool = FHE.asEbool(_a);
        pubkey = _pubkey;
    }
    function load8(inEuint8 calldata _a, bytes32 _pubkey) public {
        a8 = FHE.asEuint8(_a);
        pubkey = _pubkey;
    }
    function load16(inEuint16 calldata _a, bytes32 _pubkey) public {
        a16 = FHE.asEuint16(_a);
        pubkey = _pubkey;
    }
    function load32(inEuint32 calldata _a, bytes32 _pubkey) public {
        a32 = FHE.asEuint32(_a);
        pubkey = _pubkey;
    }
    function load64(inEuint64 calldata _a, bytes32 _pubkey) public {
        a64 = FHE.asEuint64(_a);
        pubkey = _pubkey;
    }
    function load128(inEuint128 calldata _a, bytes32 _pubkey) public {
        a128 = FHE.asEuint128(_a);
        pubkey = _pubkey;
    }
    function load256(inEuint256 calldata _a, bytes32 _pubkey) public {
        a256 = FHE.asEuint256(_a);
        pubkey = _pubkey;
    }
    function loadAddress(inEaddress calldata _a, bytes32 _pubkey) public {
        aAddress = FHE.asEaddress(_a);
        pubkey = _pubkey;
    }

    function benchSealoutputBool() public view {
        FHE.sealoutput(aBool, pubkey);
    }
    function benchSealoutput8() public view {
        FHE.sealoutput(a8, pubkey);
    }
    function benchSealoutput16() public view {
        FHE.sealoutput(a16, pubkey);
    }
    function benchSealoutput32() public view {
        FHE.sealoutput(a32, pubkey);
    }
    function benchSealoutput64() public view {
        FHE.sealoutput(a64, pubkey);
    }
    function benchSealoutput128() public view {
        FHE.sealoutput(a128, pubkey);
    }
    function benchSealoutput256() public view {
        FHE.sealoutput(a256, pubkey);
    }
    function benchSealoutputAddress() public view {
        FHE.sealoutput(aAddress, pubkey);
    }
}
