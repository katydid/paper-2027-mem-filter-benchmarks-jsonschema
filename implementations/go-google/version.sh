#!/bin/bash

set -o errexit
set -o nounset

# Extract the version of the jsonschema module from go.mod
jsonschema_version=$(grep 'github.com/google/jsonschema-go' implementations/go-google/go.mod | sed -E 's/.*\/jsonschema-go v([0-9\.]+)/\1/')

# Output the version
echo "$jsonschema_version"
