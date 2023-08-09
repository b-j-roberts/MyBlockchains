#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR="${SCRIPT_DIR}/.."

if [ -z $SEQUENCER_OUTPUT_FILE ]; then
  SEQUENCER_OUTPUT_FILE="${WORK_DIR}/out-sequencer.txt"
fi

while true; do
  cd ${WORK_DIR} && SEQUENCER_OUTPUT_FILE=${SEQUENCER_OUTPUT_FILE} make run-sequencer
  sleep 20
  if cat ${SEQUENCER_OUTPUT_FILE} | grep Commit\ new\ sealing\ work | grep -q number=5; then
    break
  fi
  echo "Sequencer failed to start, restarting..."
  cat $SEQUENCER_OUTPUT_FILE
  ps aux | grep scripts/../build/sequencer | awk '{print $2}' | xargs kill
  sleep 5
done

echo "Sequencer started successfully"
