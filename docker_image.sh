#!/bin/sh
set -x
set -e
#
# This script is meant for quick & easy build hfs via:
#   'curl -sSL https://raw.githubusercontent.com/carsonsx/hfs/master/docker_image.sh | sh'
# or:
#   'wget -qO- https://raw.githubusercontent.com/carsonsx/hfs/master/docker_image.sh | sh'
#

# Download Dockerfile
curl -sSLO "https://raw.githubusercontent.com/carsonsx/hfs/master/Dockerfile"

# Build and push image
docker build -t carsonsx/hfs .
docker push carsonsx/hfs
docker tag carsonsx/hfs carsonsx/hfs:0.1
docker push carsonsx/hfs:0.1
