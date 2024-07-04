// SPDX-License-Identifier: BSD-3-Clause-Clear

pragma solidity >=0.8.19 <0.9.0;

import { FHE, euint8, inEuint8 } from "../../FHE.sol";

contract Counter {
  euint8 private counter;
  mapping(uint256 => euint8) private counterMapping;

  function add(inEuint8 calldata encryptedValue) public {
    euint8 value = FHE.asEuint8(encryptedValue);
    counter = FHE.add(counter, value);
  }

  function addMapping(uint8 idx, inEuint8 calldata encryptedValue) public {
    euint8 value = FHE.asEuint8(encryptedValue);
    counterMapping[idx] = FHE.add(counterMapping[idx], value);
  }

  function getCounter(bytes32 publicKey) public view returns (string memory) {
    return FHE.sealoutput(counter, publicKey);
  }

  function getCounterMapping(
    uint8 idx,
    bytes32 publicKey
  ) public view returns (string memory) {
    return FHE.sealoutput(counterMapping[idx], publicKey);
  }
}
