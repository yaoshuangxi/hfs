#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix hfs -ldflags "-s -w" -o hfs

