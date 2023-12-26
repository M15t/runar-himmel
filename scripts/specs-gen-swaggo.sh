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

swaggerui_path="./cmd/api/main.go"

# Generate swagger.json file
cd $swaggerui_path
swag i -g swaggerui_path -o cmd/docs
swag fmt

