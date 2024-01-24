// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {FHE} from "../../FHE.sol";
import {ebool, euint8} from "../../FHE.sol";
import {Utils} from "./utils/Utils.sol";

error TestNotFound(string test);

contract SealoutputTest {
    using Utils for *;
    
    function sealoutput(string calldata test, uint256 a, bytes32 pubkey) public pure returns (bytes memory reencrypted) {
        if (Utils.cmp(test, "sealoutput(euint8)")) {
            return FHE.sealoutput(FHE.asEuint8(a), pubkey);
        } else if (Utils.cmp(test, "sealoutput(euint16)")) {
            return FHE.sealoutput(FHE.asEuint16(a), pubkey);
        } else if (Utils.cmp(test, "sealoutput(euint32)")) {
            return FHE.sealoutput(FHE.asEuint32(a), pubkey);
        } else if (Utils.cmp(test, "sealoutput(ebool)")) {
            bool b = true;
            if (a == 0) {
                b = false;
            }

            return FHE.sealoutput(FHE.asEbool(b), pubkey);
        } else if (Utils.cmp(test, "seal(euint8)")) {
            euint8 aEnc = FHE.asEuint8(a);
            return aEnc.seal(pubkey);
        }
        revert TestNotFound(test);
    }
}
