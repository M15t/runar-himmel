#!/usr/bin/env bash

set -e

# Only generate specs for development environments
if [[ "$SWAGGER" != "true" ]]; then
	echo "no specs generated"
	exit 0
fi

# Workaround for difference between sed command on linux & mac
sed_cmd=( sed -i )
if [[ "$(uname)" == "Darwin" ]]; then
	sed_cmd=( sed -i '' )
fi

set -x

swaggerui_path="cmd/api/swagger-ui"

# Generate swagger.json file
swagger -q generate spec --scan-models --compact -t swagger -w cmd/api/ -o $swaggerui_path/swagger.json

# Replace version by corresponding env var
COMMITID=$(git rev-parse --short HEAD)
"${sed_cmd[@]}" -e "s#%{VERSION}#$COMMITID#g"  $swaggerui_path/swagger.json
