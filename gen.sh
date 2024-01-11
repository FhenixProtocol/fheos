#!/usr/bin/env bash

set -e

FHE_OPS_DEST=${1:-"../precompiles"}
OUTPUT=${2:-"chains/arbitrum"}

GEN_FHEOPS=${3:-"true"}

go run gen.go 1
if [ ! -e precompiles/contracts ]; then
    mkdir precompiles/contracts
fi
cp FheOs_gen.sol precompiles/contracts/FheOs.sol
mv FheOs_gen.sol solidity/FheOS.sol
cd precompiles
rm -rf artifacts
pnpm build
rm -r ./contracts/
cd ../
go run gen.go 2 $OUTPUT

if [ "${GEN_FHEOPS}" = "true" ]; then
    cp FheOps_gen.go "$FHE_OPS_DEST"/FheOps.go
fi

rm *_gen*