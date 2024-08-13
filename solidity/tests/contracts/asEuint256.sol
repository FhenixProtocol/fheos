// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE, euint8, euint16, euint32, euint64, euint128, euint256, ebool} from "../../FHE.sol";
import {inEuint256} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEuint256Test {
    using Utils for *;
    
    
    function castFromEboolToEuint256(uint256 val, string calldata test) public pure returns (uint256) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEbool(val).toU256().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint256(FHE.asEbool(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint8ToEuint256(uint256 val, string calldata test) public pure returns (uint256) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint8(val).toU256().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint256(FHE.asEuint8(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint16ToEuint256(uint256 val, string calldata test) public pure returns (uint256) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint16(val).toU256().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint256(FHE.asEuint16(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint32ToEuint256(uint256 val, string calldata test) public pure returns (uint256) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint32(val).toU256().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint256(FHE.asEuint32(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint64ToEuint256(uint256 val, string calldata test) public pure returns (uint256) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint64(val).toU256().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint256(FHE.asEuint64(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint128ToEuint256(uint256 val, string calldata test) public pure returns (uint256) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint128(val).toU256().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint256(FHE.asEuint128(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEaddressToEuint256(uint256 val, string calldata test) public pure returns (uint256) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEaddress(val).toU256().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEuint256(FHE.asEaddress(val)));
        }
        revert TestNotFound(test);
    }

    function castFromPlaintextToEuint256(uint256 val) public pure returns (uint256) {
        return FHE.decrypt(FHE.asEuint256(val));
    }

    function castFromPreEncryptedToEuint256(inEuint256 calldata val) public pure returns (uint256) {
        return FHE.decrypt(FHE.asEuint256(val));
    }
}
