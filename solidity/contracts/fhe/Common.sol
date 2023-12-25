// SPDX-License-Identifier: BSD-3-Clause-Clear
pragma solidity >=0.8.20 <0.9.0;

library Common {
    // Values used to communicate types to the runtime.
    uint8 internal constant EBOOL_TFHE_GO = 0;
    uint8 internal constant EUINT8_TFHE_GO = 0;
    uint8 internal constant EUINT16_TFHE_GO = 1;
    uint8 internal constant EUINT32_TFHE_GO = 2;

    function bigIntToBool(uint256 i) internal pure returns (bool) {
        return (i > 0);
    }

    function bigIntToUint8(uint256 i) internal pure returns (uint8) {
        return uint8(i);
    }

    function bigIntToUint16(uint256 i) internal pure returns (uint16) {
        return uint16(i);
    }

    function bigIntToUint32(uint256 i) internal pure returns (uint32) {
        return uint32(i);
    }

    function bigIntToUint64(uint256 i) internal pure returns (uint64) {
        return uint64(i);
    }

    function bigIntToUint128(uint256 i) internal pure returns (uint128) {
        return uint128(i);
    }

    function bigIntToUint256(uint256 i) internal pure returns (uint256) {
        return i;
    }
}