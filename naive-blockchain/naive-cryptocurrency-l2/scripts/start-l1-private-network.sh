#!/bin/bash
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

$BASEDIR/../../../eth-private-network/scripts/launch-network-overwrite.sh
