// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {inEaddress} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract AsEaddressTest {
    using Utils for *;
    
    
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

    function castFromPreEncryptedToEaddress(inEaddress calldata val) public pure returns (address) {
        return FHE.decrypt(FHE.asEaddress(val));
    }
}
