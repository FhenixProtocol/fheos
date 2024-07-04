// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {inEuint128} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEuint128Test {
    using Utils for *;
    
    
    function castFromEboolToEuint128(uint256 val, string calldata test) public pure returns (uint128) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEbool(val).toU128().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint128(FHE.asEbool(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint8ToEuint128(uint256 val, string calldata test) public pure returns (uint128) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint8(val).toU128().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint128(FHE.asEuint8(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint16ToEuint128(uint256 val, string calldata test) public pure returns (uint128) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint16(val).toU128().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint128(FHE.asEuint16(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint32ToEuint128(uint256 val, string calldata test) public pure returns (uint128) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint32(val).toU128().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint128(FHE.asEuint32(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint64ToEuint128(uint256 val, string calldata test) public pure returns (uint128) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint64(val).toU128().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint128(FHE.asEuint64(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint256ToEuint128(uint256 val, string calldata test) public pure returns (uint128) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint256(val).toU128().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint128(FHE.asEuint256(val)));
        }
        revert TestNotFound(test);
    }

    function castFromPlaintextToEuint128(uint256 val) public pure returns (uint128) {
        return FHE.decrypt(FHE.asEuint128(val));
    }

    function castFromPreEncryptedToEuint128(inEuint128 calldata val) public pure returns (uint128) {
        return FHE.decrypt(FHE.asEuint128(val));
    }
}
