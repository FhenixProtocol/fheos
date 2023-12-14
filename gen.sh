#!/usr/bin/env bash

set -e

FHE_OPS_DEST=${1:-"../precompiles"}
OUTPUT=${2:-"chains/arbitrum"}

go run gen.go 1
if [ ! -e precompiles/contracts ]; then
    mkdir precompiles/contracts

fi
cp FheOs_gen.sol precompiles/contracts/FheOs.sol
mv FheOs_gen.sol solidity/FheOS.sol
cd precompiles
rm -rf artifacts
yarn
yarn build
rm -r ./contracts/
cd ../
go run gen.go 2 $OUTPUT
cp FheOps_gen.go "$FHE_OPS_DEST"/FheOps.go
rm *_gen*