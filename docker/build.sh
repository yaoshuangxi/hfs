#!/usr/bin/env bash
set -x
#set -e
set -o pipefail
#
# This script is meant for quick & easy build and push hfs image via:
#   'curl -sSL https://raw.githubusercontent.com/carsonsx/hfs/master/docker/build.sh | sh'
# or:
#   'wget -qO- https://raw.githubusercontent.com/carsonsx/hfs/master/docker/build.sh | sh'
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

mv -f hfs_linux_amd64 hfs

#VERSION="$(git describe --tags --always)"
VERSION="$(./hfs -v)"

# Build image
docker build -t carsonsx/hfs .
docker tag carsonsx/hfs carsonsx/hfs:${VERSION}

# Push image
# Please run docker login first
docker push carsonsx/hfs
docker push carsonsx/hfs:${VERSION}

# Clean
docker rmi carsonsx/hfs:${VERSION}
docker rmi carsonsx/hfs
rm -rf Dockerfile