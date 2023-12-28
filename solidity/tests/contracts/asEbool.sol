// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {TFHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEboolTest {
    using Utils for *;
    
    
    function castFromEuint8ToEbool(uint256 val, string calldata test) public pure returns (bool) {
        if (Utils.cmp(test, "bound")) {
            return TFHE.decrypt(TFHE.asEuint8(val).toBool());
        } else if (Utils.cmp(test, "regular")) {
            return TFHE.decrypt(TFHE.asEbool(TFHE.asEuint8(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint16ToEbool(uint256 val, string calldata test) public pure returns (bool) {
        if (Utils.cmp(test, "bound")) {
            return TFHE.decrypt(TFHE.asEuint16(val).toBool());
        } else if (Utils.cmp(test, "regular")) {
            return TFHE.decrypt(TFHE.asEbool(TFHE.asEuint16(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint32ToEbool(uint256 val, string calldata test) public pure returns (bool) {
        if (Utils.cmp(test, "bound")) {
            return TFHE.decrypt(TFHE.asEuint32(val).toBool());
        } else if (Utils.cmp(test, "regular")) {
            return TFHE.decrypt(TFHE.asEbool(TFHE.asEuint32(val)));
        }
        revert TestNotFound(test);
    }

    function castFromPlaintextToEbool(uint256 val) public pure returns (bool) {
        return TFHE.decrypt(TFHE.asEbool(val));
    }

    function castFromPreEncryptedToEbool(bytes memory val) public pure returns (bool) {
        return TFHE.decrypt(TFHE.asEbool(val));
    }
}
