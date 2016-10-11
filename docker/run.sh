#!/bin/sh
set -x
set -e
#
# This script is meant for quick & easy run latest image via:
#   'curl -sSL https://raw.githubusercontent.com/carsonsx/hfs/master/docker/run.sh | sh'
# or:
#   'wget -qO- https://raw.githubusercontent.com/carsonsx/hfs/master/docker/run.sh | sh'
#

# Download Dockerfile
curl -sSLO "https://raw.githubusercontent.com/carsonsx/hfs/master/docker/Dockerfile"

# Stop and remove image if exists
docker stop hfs
docker rm hfs

# Build image
docker rmi carsonsx/hfs
docker build -t carsonsx/hfs .

# Run
docker run -itd --name hfs --restart=always -v ~/hfs/files:/files:rw -p 80:8100 carsonsx/hfs