#!/usr/bin/env bash

# Get the git commit
GIT_COMMIT="$(git rev-parse --short HEAD)"
GIT_DIRTY="$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)"
GIT_DESCRIBE="$(git describe --tags --always)"

echo "git commit: ${GIT_COMMIT}"
echo "git dirty: ${GIT_DIRTY}"
echo "git describe: ${GIT_DESCRIBE}"

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -ldflags "-s -w -X main.GitCommit='${GIT_COMMIT}${GIT_DIRTY}' -X main.GitDescribe='${GIT_DESCRIBE}'" -o bin/hfs_darwin_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-s -w -X main.GitCommit='${GIT_COMMIT}${GIT_DIRTY}' -X main.GitDescribe='${GIT_DESCRIBE}'" -o bin/hfs_linux_amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags "-s -w -X main.GitCommit='${GIT_COMMIT}${GIT_DIRTY}' -X main.GitDescribe='${GIT_DESCRIBE}'" -o bin/hfs_windows_amd64.exe

