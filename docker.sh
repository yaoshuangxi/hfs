#!/usr/bin/env bash

# mkdir ~/hfs in docker machine
# copy Dockerfile, hfs, docker.sh to ~/hfs
# run ./docker.sh

docker stop hfs
docker rm hfs
docker build -t carsonsx/hfs .
docker run -itd --name hfs -v ~/hfs/files:/files:rw -p 80:8011 --restart=always carsonsx/hfs