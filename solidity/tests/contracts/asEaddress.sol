// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEaddressTest {
    using Utils for *;
    
    
    function castFromEboolToEaddress(uint256 val, string calldata test) public pure returns (address) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEbool(val).toEaddress().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEaddress(FHE.asEbool(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint8ToEaddress(uint256 val, string calldata test) public pure returns (address) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint8(val).toEaddress().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEaddress(FHE.asEuint8(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint16ToEaddress(uint256 val, string calldata test) public pure returns (address) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint16(val).toEaddress().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEaddress(FHE.asEuint16(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint32ToEaddress(uint256 val, string calldata test) public pure returns (address) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint32(val).toEaddress().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEaddress(FHE.asEuint32(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint64ToEaddress(uint256 val, string calldata test) public pure returns (address) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint64(val).toEaddress().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEaddress(FHE.asEuint64(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint128ToEaddress(uint256 val, string calldata test) public pure returns (address) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint128(val).toEaddress().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEaddress(FHE.asEuint128(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint256ToEaddress(uint256 val, string calldata test) public pure returns (address) {
        if (Utils.cmp(test, "bound")) {
            return FHE.asEuint256(val).toEaddress().decrypt();
        } else if (Utils.cmp(test, "regular")) {
            return FHE.decrypt(FHE.asEaddress(FHE.asEuint256(val)));
        }
        revert TestNotFound(test);
    }

    function castFromPlaintextToEaddress(uint256 val) public pure returns (address) {
        return FHE.decrypt(FHE.asEaddress(val));
    }

    function castFromPreEncryptedToEaddress(bytes memory val) public pure returns (address) {
        return FHE.decrypt(FHE.asEaddress(val));
    }
}
