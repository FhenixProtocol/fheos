// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEboolTest {
    using Utils for *;
    
    
    function castFromEuint8ToEbool(uint256 val, string calldata test) public pure returns (bool) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint8(val).toBool().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEbool(FHE.asEuint8(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint16ToEbool(uint256 val, string calldata test) public pure returns (bool) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint16(val).toBool().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEbool(FHE.asEuint16(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint32ToEbool(uint256 val, string calldata test) public pure returns (bool) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint32(val).toBool().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEbool(FHE.asEuint32(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint64ToEbool(uint256 val, string calldata test) public pure returns (bool) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint64(val).toBool().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEbool(FHE.asEuint64(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint128ToEbool(uint256 val, string calldata test) public pure returns (bool) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint128(val).toBool().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEbool(FHE.asEuint128(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint256ToEbool(uint256 val, string calldata test) public pure returns (bool) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint256(val).toBool().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEbool(FHE.asEuint256(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEaddressToEbool(uint256 val, string calldata test) public pure returns (bool) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEaddress(val).toBool().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEbool(FHE.asEaddress(val)));
        }
        revert TestNotFound(test);
    }

    function castFromPlaintextToEbool(uint256 val) public pure returns (bool) {
        return FHE.decrypt(FHE.asEbool(val));
    }

    function castFromPreEncryptedToEbool(bytes memory val) public pure returns (bool) {
        return FHE.decrypt(FHE.asEbool(val));
    }
}
