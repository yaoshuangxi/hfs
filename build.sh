#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -a -ldflags "-s -w -X main.GitCommit=`git rev-parse --short HEAD`" -o bin/hfs_darwin_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -ldflags "-s -w -X main.GitCommit=`git rev-parse --short HEAD`" -o bin/hfs_linux_amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -a -ldflags "-s -w -X main.GitCommit=`git rev-parse --short HEAD`" -o bin/hfs_windows_amd64.exe

