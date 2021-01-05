#!/bin/bash
set -e
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
PWD=$(pwd)
[ ! -f $(pwd)/go.mod ] && echo "Please run prepare_dev.sh in root workspace" && exit 2

export os=darwin
export nodejs="node-v14.15.3-"$os"-x64"
nodejspkg=$nodejs".tar.xz"

echo "Downloading NodeJS 14"
mkdir -p ./dist
[ ! -f ./dist/$nodejspkg ] && curl -o ./dist/$nodejspkg https://nodejs.org/dist/v14.15.3/$nodejspkg
echo "Extracting NodeJS bin from "$PWD/dist/$nodejs
tar -xvf $PWD/dist/$nodejspkg -C $PWD/bin --strip-components 2 $nodejs/bin/node 
echo "Done"
