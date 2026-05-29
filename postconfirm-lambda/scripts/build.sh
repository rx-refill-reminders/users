#!/bin/bash -eu

rm -rf dist
mkdir dist

# Build the executable
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build \
    -trimpath \
    -buildvcs=false \
    -ldflags="-s -w" \
    -o cmd/bootstrap \
    ./cmd

# Compress the executable
cd cmd
zip ../dist/api.zip bootstrap > /dev/null

# Cleanup the uncompressed executable
rm bootstrap
