// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {TFHE} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEboolTest {
    using Utils for *;

    function castFromEuint8ToEbool(uint256 val) public pure returns (bool) {
        return FHE.decrypt(FHE.asEbool(FHE.asEuint8(val)));
    }

    function castFromEuint16ToEbool(uint256 val) public pure returns (bool) {
        return FHE.decrypt(FHE.asEbool(FHE.asEuint16(val)));
    }

    function castFromEuint32ToEbool(uint256 val) public pure returns (bool) {
        return FHE.decrypt(FHE.asEbool(FHE.asEuint32(val)));
    }

    function castFromPlaintextToEbool(uint256 val) public pure returns (bool) {
        return FHE.decrypt(FHE.asEbool(val));
    }

    function castFromPreEncryptedToEbool(bytes memory val) public pure returns (bool) {
        return FHE.decrypt(FHE.asEbool(val));
    }
}
