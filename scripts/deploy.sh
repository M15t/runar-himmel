#!/usr/bin/env bash

set -e

STAGE="$1"
if [[ "$STAGE" == "" ]]; then
	STAGE=dev
fi

ARGS="$@"

# This script must be run from the project root directory
cd ./deploy

# install dependencies when running for the first time
if [[ ! -d "./node_modules" ]]; then bun install; fi

# ! replace npx by bun x for bun testing
bun x serverless deploy --stage $STAGE
bun x serverless functions:invoke --stage $STAGE --function Migration
# ! turn on when needed
# npx serverless functions:invoke --stage $STAGE --function Seed
