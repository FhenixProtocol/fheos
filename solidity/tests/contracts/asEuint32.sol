// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {inEuint32} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEuint32Test {
    using Utils for *;
    
    
    function castFromEboolToEuint32(uint256 val, string calldata test) public pure returns (uint32) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEbool(val).toU32().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint32(FHE.asEbool(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint8ToEuint32(uint256 val, string calldata test) public pure returns (uint32) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint8(val).toU32().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint32(FHE.asEuint8(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint16ToEuint32(uint256 val, string calldata test) public pure returns (uint32) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint16(val).toU32().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint32(FHE.asEuint16(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint64ToEuint32(uint256 val, string calldata test) public pure returns (uint32) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint64(val).toU32().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint32(FHE.asEuint64(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint128ToEuint32(uint256 val, string calldata test) public pure returns (uint32) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint128(val).toU32().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint32(FHE.asEuint128(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint256ToEuint32(uint256 val, string calldata test) public pure returns (uint32) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint256(val).toU32().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint32(FHE.asEuint256(val)));
        }
        revert TestNotFound(test);
    }

    function castFromPlaintextToEuint32(uint256 val) public pure returns (uint32) {
        return FHE.decrypt(FHE.asEuint32(val));
    }

    function castFromPreEncryptedToEuint32(inEuint32 calldata val) public pure returns (uint32) {
        return FHE.decrypt(FHE.asEuint32(val));
    }
}
