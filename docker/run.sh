#!/usr/bin/env bash
set -x
#set -e
set -o pipefail
#
# This script is meant for quick & easy run latest image via:
#   'curl -sSL https://raw.githubusercontent.com/carsonsx/hfs/master/docker/run.sh | sh'
# or:
#   'wget -qO- https://raw.githubusercontent.com/carsonsx/hfs/master/docker/run.sh | sh'
#

# Download Dockerfile
curl -sSLO "https://raw.githubusercontent.com/carsonsx/hfs/master/docker/Dockerfile"

# Stop and remove image if exists
result=`docker ps | grep hfs`
if [ -n "$result" ]
then
	docker stop hfs
fi
result=`docker ps -a | grep hfs`
if [ -n "$result" ]
then
	docker rm hfs
fi

# Build image
result=`docker images | grep carsonsx/hfs`
if [ -n "$result" ]
then
	docker rmi carsonsx/hfs
fi
docker build -t carsonsx/hfs .
rm -rf Dockerfile

# Run
docker run -itd --name hfs --restart=always -v ~/hfs/files:/files:rw -p 80:8011 carsonsx/hfs
docker logs -f hfs