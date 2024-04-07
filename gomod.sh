#!/bin/bash
set -eu

if [ -f ./go.mod ]; then
    exit 0
fi

touch go.mod

CURRENT_DIR=$(basename $(pwd))

CONTENT=$(cat <<-EOD
module github.com/MoreiraVinicius/${CURRENT_DIR}

go 1.21.3

require github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos v0.3.6
EOD
)

echo "$CONTENT" > go.mod
