// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEuint64Test {
    using Utils for *;
    
    
    function castFromEboolToEuint64(uint256 val, string calldata test) public pure returns (uint64) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEbool(val).toU64().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint64(FHE.asEbool(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint8ToEuint64(uint256 val, string calldata test) public pure returns (uint64) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint8(val).toU64().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint64(FHE.asEuint8(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint16ToEuint64(uint256 val, string calldata test) public pure returns (uint64) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint16(val).toU64().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint64(FHE.asEuint16(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint32ToEuint64(uint256 val, string calldata test) public pure returns (uint64) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint32(val).toU64().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint64(FHE.asEuint32(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint128ToEuint64(uint256 val, string calldata test) public pure returns (uint64) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint128(val).toU64().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint64(FHE.asEuint128(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint256ToEuint64(uint256 val, string calldata test) public pure returns (uint64) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint256(val).toU64().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint64(FHE.asEuint256(val)));
        }
        revert TestNotFound(test);
    }

    function castFromPlaintextToEuint64(uint256 val) public pure returns (uint64) {
        return FHE.decrypt(FHE.asEuint64(val));
    }

    function castFromPreEncryptedToEuint64(bytes memory val) public pure returns (uint64) {
        return FHE.decrypt(FHE.asEuint64(val));
    }
}
