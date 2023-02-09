#!/bin/bash

ENVV=$1

if [ $ENVV = 'arm64' ];then
  CGO_ENABLED=on GOARCH=arm64 GOOS=linux go build -o bin/arm64/ipicture  main.go
elif [ $ENVV = 'windows' ]; then
  CGO_ENABLED=on GOARCH=amd64 GOOS=windows go build -o bin/windows/ipicture main.go
elif [ $ENVV = 'linux' ]; then
  CGO_ENABLED=on GOARCH=ubuntu GOOS=linux go build -o bin/linux/ipicture main.go
fi

