#!/bin/bash

set -e

PID=".pid"

# Color variables
green='\033[0;32m'
yellow='\033[0;33m'
magenta='\033[0;35m'
# Clear the color after that
clear='\033[0m'

prefix='[watcher] '

# From http://stackoverflow.com/a/12498485
function relative_path {
	# both $1 and $2 are absolute paths beginning with /
	# returns relative path to $2 from $1
	local source=$1
	local target=$2

	local commonPart=$source
	local result=""

	while [[ "${target#$commonPart}" == "${target}" ]]; do
		# no match, means that candidate common part is not correct
		# go up one level (reduce common part)
		commonPart="$(dirname $commonPart)"
		# and record that we went back, with correct / handling
		if [[ -z $result ]]; then
			result=".."
		else
			result="../$result"
		fi
	done

	if [[ $commonPart == "/" ]]; then
		# special case for root (no common path)
		result="$result/"
	fi

	# since we now have identified the common part,
	# compute the non-common part
	local forwardPart="${target#$commonPart}"

	# and now stick all parts together
	if [[ -n $result ]] && [[ -n $forwardPart ]]; then
		result="$result$forwardPart"
	elif [[ -n $forwardPart ]]; then
		# extra slash removal
		result="${forwardPart:1}"
	fi

	echo $result
}
export -f relative_path

function wait_for_changes {
	echo -e "${green}${prefix}watching for file changes...${clear}"
	fswatch -1 -l 1 -e ".*" -i "\\.go$" -i "\\.env$" -i "swagger-ui/.*" --recursive .env ./cmd/ ./config/ ./pkg/ ./internal/  | xargs -I{} bash -c 'relative_path $(pwd) "$@"' _ {} | xargs printf "${magenta}${prefix}file changed: %s${clear}\n"
}

function start_server {
	go run -tags swagger ./cmd/api & echo $! > $PID
}

function kill_server {
	APP_PID="$(lsof -t -i :8083 || true)"
	if [[ "$APP_PID" != "" ]]; then kill -INT $APP_PID 2>&1 > /dev/null || true; fi
}

function reload_server {
	echo -e "${yellow}${prefix}reloading server...${clear}"
	kill -INT $(pgrep -P `cat $PID`) || true
	start_server
}

if [[ "$1" == "terminate" ]]; then
	kill_server
	exit 0
fi

# Exit on ctrl-c (without this, ctrl-c would go to fswatch, causing it to
# reload instead of exit):
trap 'exit 0' SIGINT

start_server

while true; do
	wait_for_changes
	reload_server
done
