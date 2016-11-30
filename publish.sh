#!/usr/bin/env bash

# Login docker
docker login

# pull golang image for build
docker pull golang

# build in docker
docker run --rm -v "$PWD":/hfs -w /hfs golang ./build.sh

# copy build result
\cp -f bin/hfs_linux_amd64 docker/hfs

# Stop and remove image if exists
docker rm -f hfs

set -e

# fetch version
cd docker
VERSION="$(./hfs --getversion)"

# Build image
docker build -t carsonsx/hfs .
docker tag carsonsx/hfs carsonsx/hfs:${VERSION}

# Push image
docker push carsonsx/hfs
docker push carsonsx/hfs:${VERSION}

# Clean
docker rmi carsonsx/hfs
docker rmi carsonsx/hfs:${VERSION}
rm -f hfs