#!/usr/bin/env bash

set -e

# Default values
FHE_OPS_DEST="../precompiles"
OUTPUT="chains/arbitrum"
GEN_FHEOPS="false"
NITRO_OVERRIDE="false"

# Parse flags
while getopts "d:o:g:n:" opt; do
  case $opt in
    d) FHE_OPS_DEST=$OPTARG ;;
    o) OUTPUT=$OPTARG ;;
    g) GEN_FHEOPS=$OPTARG ;;
    n) NITRO_OVERRIDE=$OPTARG ;;
    \?) echo "Invalid option -$OPTARG" >&2
        exit 1
    ;;
  esac
done

# Set FHE_OPS_DEST to "nitro-overrides/precompiles" if the nitro override flag is true
if [ "$NITRO_OVERRIDE" = "true" ]; then
    FHE_OPS_DEST="nitro-overrides/precompiles"
fi

# Rest of the script remains the same
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
    echo "Generating FheOps.go... in $FHE_OPS_DEST"
    cp FheOps_gen.go "$FHE_OPS_DEST"/FheOps.go
fi

rm *_gen*