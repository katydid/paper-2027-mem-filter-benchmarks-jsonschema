#!/bin/sh

set -o errexit
set -o nounset

jq --raw-output '.dependencies["ajv"]' < implementations/ajv/package.json
