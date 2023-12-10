#!/usr/bin/env bash

set -e

PRECOMPILES=${1:-"precompiles"}
FHE_OPS_DEST=${2:-".."}
OUTPUT=${3:-"chains/arbitrum"}

go run gen.go 1
cp FheOps_gen.sol "$PRECOMPILES"/contracts/FheOps.sol
cd "$PRECOMPILES"
rm -rf artifacts
yarn
yarn build
cd ..
go run gen.go 2 $OUTPUT
cp FheOps_gen.go "$FHE_OPS_DEST"/FheOps.go
rm *_gen*