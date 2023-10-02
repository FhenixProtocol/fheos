// SPDX-License-Identifier: BSD-3-Clause-Clear

pragma solidity >=0.8.13 <0.9.0;

import "https://github.com/FhenixProtocol/fhevm/lib/TFHE.sol";

contract Lior {
    uint private counter;

    function add() public {
        counter = TFHE.lior(4, 6);
    }

    function getCounter() public view returns (uint) {
        return counter;
    }
}
