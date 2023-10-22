#!/usr/bin/env bash

set -e

cd precompiles
yarn
cd ../
go run gen.go