#!/bin/sh

set -o errexit
set -o nounset

jq --raw-output '.dependencies["@hyperjump/json-schema"]' < implementations/hyperjump/package.json
