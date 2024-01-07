// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEuint16Test {
    using Utils for *;
    
    
    function castFromEboolToEuint16(uint256 val, string calldata test) public pure returns (uint16) {
        if (Utils.cmp(test, "bound")) {
            return FHE.decrypt(FHE.asEbool(val).toU16());
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint16(FHE.asEbool(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint8ToEuint16(uint256 val, string calldata test) public pure returns (uint16) {
        if (Utils.cmp(test, "bound")) {
            return FHE.decrypt(FHE.asEuint8(val).toU16());
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint16(FHE.asEuint8(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint32ToEuint16(uint256 val, string calldata test) public pure returns (uint16) {
        if (Utils.cmp(test, "bound")) {
            return FHE.decrypt(FHE.asEuint32(val).toU16());
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint16(FHE.asEuint32(val)));
        }
        revert TestNotFound(test);
    }

    function castFromPlaintextToEuint16(uint256 val) public pure returns (uint16) {
        return FHE.decrypt(FHE.asEuint16(val));
    }

    function castFromPreEncryptedToEuint16(bytes memory val) public pure returns (uint16) {
        return FHE.decrypt(FHE.asEuint16(val));
    }
}
