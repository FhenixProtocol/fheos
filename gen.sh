#!/usr/bin/env bash

set -e

go run gen.go 1
cd precompiles
yarn
cd ../
go run gen.go 2
cp FheOps_gen.go ../precompiles/FheOps.go
cp FheOps_gen.sol ./precompiles/contracts/FheOps.sol
rm *_gen*