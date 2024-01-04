// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEuint8Test {
    using Utils for *;
    
    
    function castFromEboolToEuint8(uint256 val, string calldata test) public pure returns (uint8) {
        if (Utils.cmp(test, "bound")) {
            return FHE.decrypt(FHE.asEbool(val).toU8());
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint8(FHE.asEbool(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint16ToEuint8(uint256 val, string calldata test) public pure returns (uint8) {
        if (Utils.cmp(test, "bound")) {
            return FHE.decrypt(FHE.asEuint16(val).toU8());
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint8(FHE.asEuint16(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint32ToEuint8(uint256 val, string calldata test) public pure returns (uint8) {
        if (Utils.cmp(test, "bound")) {
            return FHE.decrypt(FHE.asEuint32(val).toU8());
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint8(FHE.asEuint32(val)));
        }
        revert TestNotFound(test);
    }

    function castFromPlaintextToEuint8(uint256 val) public pure returns (uint8) {
        return FHE.decrypt(FHE.asEuint8(val));
    }

    function castFromPreEncryptedToEuint8(bytes memory val) public pure returns (uint8) {
        return FHE.decrypt(FHE.asEuint8(val));
    }
}
