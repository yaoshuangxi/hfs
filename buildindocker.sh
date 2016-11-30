#!/usr/bin/env bash

docker pull golang
docker run --rm -v "$PWD":/hfs -w /hfs golang ./build.sh
