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
rm -f Dockerfile
curl -sSLO "https://raw.githubusercontent.com/carsonsx/hfs/master/docker/Dockerfile"

# Stop and remove image if exists
docker rm -f hfs

mv -f hfs_linux_amd64 hfs

set -e

#VERSION="$(git describe --tags --always)"
VERSION="$(./hfs --getversion)"

# Build image
docker build -t carsonsx/hfs .
docker tag carsonsx/hfs carsonsx/hfs:${VERSION}

# Push image
# Please run docker login first
docker push carsonsx/hfs
docker push carsonsx/hfs:${VERSION}

# Clean
docker rmi carsonsx/hfs
docker rmi carsonsx/hfs:${VERSION}
rm -f Dockerfile hfs