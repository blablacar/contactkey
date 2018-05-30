#!/usr/bin/env bash
set -euo pipefail
set -x
IFS=$'\n\t'

if [ "$(which upx)" = "" ]; then
	echo "upx not found"
	exit 0
fi

upx -1 dist/{darwin_amd64,linux_amd64}/cck
