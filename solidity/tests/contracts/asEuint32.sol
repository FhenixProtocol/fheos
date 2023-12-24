// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {TFHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEuint32Test {
    using Utils for *;
    
    
    function castFromEboolToEuint32(uint256 val, string calldata test) public pure returns (uint32) {
        if (Utils.cmp(test, "bound")) {
            return TFHE.decrypt(TFHE.asEbool(val).toU32());
        } else if (Utils.cmp(test, "regular")) {
            return TFHE.decrypt(TFHE.asEuint32(TFHE.asEbool(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint8ToEuint32(uint256 val, string calldata test) public pure returns (uint32) {
        if (Utils.cmp(test, "bound")) {
            return TFHE.decrypt(TFHE.asEuint8(val).toU32());
        } else if (Utils.cmp(test, "regular")) {
            return TFHE.decrypt(TFHE.asEuint32(TFHE.asEuint8(val)));
        }
        revert TestNotFound(test);
    }

    function castFromEuint16ToEuint32(uint256 val, string calldata test) public pure returns (uint32) {
        if (Utils.cmp(test, "bound")) {
            return TFHE.decrypt(TFHE.asEuint16(val).toU32());
        } else if (Utils.cmp(test, "regular")) {
            return TFHE.decrypt(TFHE.asEuint32(TFHE.asEuint16(val)));
        }
        revert TestNotFound(test);
    }

    function castFromPlaintextToEuint32(uint256 val) public pure returns (uint32) {
        return TFHE.decrypt(TFHE.asEuint32(val));
    }

    function castFromPreEncryptedToEuint32(bytes memory val) public pure returns (uint32) {
        return TFHE.decrypt(TFHE.asEuint32(val));
    }
}
