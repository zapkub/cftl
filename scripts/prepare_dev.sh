#!/bin/bash
set -e
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
PWD=$(pwd)

[ ! -f $(pwd)/go.mod ] && echo "Please run prepare_dev.sh in root workspace" && exit 2
mkdir -p $PWD/bin

$DIR/third_party.sh

echo "config file copying..."
if [ ! -f $PWD/.cftl/app.conf ]; then
    mkdir -p ./.cftl
    cp ./internal/conf/app.example.conf ./.cftl/app.conf
else
    echo "config file exists in ./.cftl skipped"
fi