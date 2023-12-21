#!/usr/bin/env bash

set -e

buildPath=deploy/functions/build

buildTags=()
if [[ "$TARGET" == "lambda" ]]; then
	buildTags+=(lambda.norpc)
	export GOOS=linux
	export GOARCH=arm64
	export CGO_ENABLED=0
fi

# generate -tags
buildTags="${buildTags[@]}"
if [[ "$buildTags" != "" ]]; then
	buildTags="-tags ${buildTags// /,}"
fi

# DRY to build functions
gobuild() {
	svcPkg=$1
	svcName=$2
	rm -rf $buildPath/$svcName
	go build -trimpath -buildvcs=true $buildTags -ldflags "-s -w" -o $buildPath/$svcName/bootstrap $svcPkg
	zip -j $buildPath/$svcName.zip $buildPath/$svcName/bootstrap
}

gobuild ./functions/migration migration
gobuild ./functions/seed seed
