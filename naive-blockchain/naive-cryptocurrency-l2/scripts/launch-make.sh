#!/bin/bash

# Script to take argument and run as a make command

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="${SCRIPT_DIR}/.."

cd $WORK_DIR && make $1
