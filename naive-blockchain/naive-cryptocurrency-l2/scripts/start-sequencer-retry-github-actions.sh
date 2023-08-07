#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="${SCRIPT_DIR}/.."

while true; do
  cd ${WORK_DIR} && make run-sequencer-github-actions
  sleep 20
  if cat ${OUTPUT_FILE} | grep Commit\ new\ sealing\ work | grep -q number=5; then
    break
  fi
  echo "Sequencer failed to start, restarting..."
  cat $OUTPUT_FILE
  ps aux | grep scripts/../build/sequencer | awk '{print $2}' | xargs kill
  sleep 5
done
