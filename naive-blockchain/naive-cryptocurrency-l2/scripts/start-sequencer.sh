#!/bin/bash

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

NAIVE_SEQUENCER_DATA="$HOME/naive-sequencer-data"

go run ${BASE_DIR}/../cmd/sequencer/
