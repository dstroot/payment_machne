#!/bin/sh

# Build program for a windows target
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -tags netgo -ldflags "-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse HEAD` -w" .
