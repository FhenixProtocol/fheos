// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../../FHE.sol";
import {
	euint8, inEuint8,
	euint16, inEuint16,
	euint32, inEuint32
} from "../../../FHE.sol";

contract DivBench {
	euint8 internal a8;
	euint16 internal a16;
	euint32 internal a32;

	euint8 internal b8;
	euint16 internal b16;
	euint32 internal b32;

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

    function benchDiv8() public view {
        FHE.div(a8, b8);
    }
    function benchDiv16() public view {
        FHE.div(a16, b16);
    }
    function benchDiv32() public view {
        FHE.div(a32, b32);
    }
}
