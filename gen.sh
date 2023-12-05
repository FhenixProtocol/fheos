#!/usr/bin/env bash

set -e

go run gen.go 1
mv FheOs_gen.sol solidity/FheOS.sol
cd precompiles
rm -rf artifacts
yarn
yarn build
cd ../
go run gen.go 2
cp FheOps_gen.go ../precompiles/FheOps.go
rm *_gen*