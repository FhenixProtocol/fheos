#!/usr/bin/env bash

set -e

cd precompiles
yarn build
cd ../
go run gen.go