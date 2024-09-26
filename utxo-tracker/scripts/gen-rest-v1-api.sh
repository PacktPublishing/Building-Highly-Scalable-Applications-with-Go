#!/bin/sh
#
# This generates the stub code for the account server REST API.
#

ROOT=$(git rev-parse --show-toplevel)/utxo-tracker
CONFIG="$ROOT/api/rest-v1-codegen-config.yaml"
SPEC="$ROOT/api/rest-v1-openapi3.yaml"

go run \
  github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 \
  --config="${CONFIG}" "${SPEC}"
