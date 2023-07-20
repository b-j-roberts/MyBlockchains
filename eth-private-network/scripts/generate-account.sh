#!/bin/bash
#
# Generate a new account key pair for geth usage

ACCOUNT_DIR="${HOME}/.eth-accounts"

display_help() {
  echo "Usage: $0 [option...] " >&2
  echo
  echo "   -h, --help                 display help"
  echo "   -d, --data-dir             data directory"
  echo "   -x, --clear                clear data directory & accounts"
  echo "   -a, --account-dir          account directory"
  echo
  echo "Example: $0 --data-dir ${HOME}/l1-miner-data/"
}

while getopts ":hxd:a:" opt; do
  case $opt in
    h)
      display_help
      exit 0
      ;;
    x)
      clear_data=true
      ;;
    d)
      data_dir=$OPTARG
      ;;
    a)
      ACCOUNT_DIR=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      display_help
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      display_help
      exit 1
      ;;
  esac
done

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORK_DIR=$SCRIPT_DIR/..

if [ -z "$data_dir" ]; then
  echo "Data directory is required"
  display_help
  exit 1
fi

if [ "$clear_data" = true ]; then
  echo "Clearing data directory: $data_dir"
  rm -rf $data_dir
  rm -rf $ACCOUNT_DIR

  mkdir -p $data_dir
  mkdir -p $ACCOUNT_DIR
fi

PASSWORD_FILE="$data_dir/password.txt"
if [ ! -f "$PASSWORD_FILE" ]; then
  echo "Generating password file: $PASSWORD_FILE"
  ACCOUNT_PASS=${ACCOUNT_PASS:-"password"}
  echo $ACCOUNT_PASS > $PASSWORD_FILE
fi

echo "Generating account key pair"
${WORK_DIR}/go-ethereum/build/bin/geth --datadir $data_dir --password $PASSWORD_FILE account new

echo "Copying account key pair to $ACCOUNT_DIR"
cp $data_dir/keystore/* $ACCOUNT_DIR
cp $data_dir/password.txt $ACCOUNT_DIR

echo "Account key pair:"
ls -l $ACCOUNT_DIR

echo "Done"
