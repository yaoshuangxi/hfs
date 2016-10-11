#!/bin/sh
set -x
set -e
#
# This script is meant for quick & easy build and push hfs image via:
#   'curl -sSL https://raw.githubusercontent.com/carsonsx/hfs/master/docker/build.sh | sh'
# or:
#   'wget -qO- https://raw.githubusercontent.com/carsonsx/hfs/master/docker/build.sh | sh'
#

# Download Dockerfile
curl -sSLO "https://raw.githubusercontent.com/carsonsx/hfs/master/docker/Dockerfile"

# Stop and remove image if exists
docker stop hfs
docker rm hfs

# Build image
docker build -t carsonsx/hfs .
docker tag carsonsx/hfs carsonsx/hfs:0.1

# Push image
docker login carsonsx
docker push carsonsx/hfs
docker push carsonsx/hfs:0.1
docker logout

# Clean
docker rmi carsonsx/hfs:0.1
docker rmi carsonsx/hfs
rm -rf Dockerfile