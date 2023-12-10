#!/usr/bin/env bash

set -e

go run gen.go 1
if [ ! -e ./precompiles/contracts ]; then
    mkdir ./precompiles/contracts

fi
cp FheOs_gen.sol ./precompiles/contracts/FheOs.sol
mv FheOs_gen.sol solidity/FheOS.sol
cd precompiles
rm -rf artifacts
yarn
yarn build
rm -r ./contracts/
cd ../
go run gen.go 2
cp FheOps_gen.go ../precompiles/FheOps.go
rm *_gen*