#!/usr/bin/env bash

set -e

go run gen.go 1
cp FheOs_gen.sol ./precompiles/contracts/FheOs.sol
mv FheOs_gen.sol solidity/FheOS.sol
cd precompiles
rm -rf artifacts
yarn
yarn build
rm ./contracts/FheOs.sol
cd ../
go run gen.go 2
cp FheOps_gen.go ../precompiles/FheOps.go
rm *_gen*