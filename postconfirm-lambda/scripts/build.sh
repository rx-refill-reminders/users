#!/bin/bash -eu

# Get the location of the build script
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

# Get the repo root
REPO_ROOT=$(git rev-parse --show-toplevel)

# Validate we are in the right directory
EXPECTED_PWD=$(cd "$SCRIPT_DIR/../src" && pwd)
ACTUAL_PWD=$(pwd)
if [[ "$ACTUAL_PWD" != "$EXPECTED_PWD" ]]; then
    >&2 echo "This should be run from '${EXPECTED_PWD#${REPO_ROOT}/}'"
    >&2 echo "Instead, ran from '${ACTUAL_PWD#${REPO_ROOT}/}'"
    exit 1
fi

rm -rf ../dist
mkdir ../dist

# Build the executable
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build \
    -trimpath \
    -buildvcs=false \
    -ldflags="-s -w" \
    -o ../dist/bootstrap \
    ./internal/cmd

# Compress the executable
cd ../dist
zip api.zip bootstrap > /dev/null

# Cleanup the uncompressed executable
rm bootstrap
