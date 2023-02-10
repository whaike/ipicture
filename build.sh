#!/bin/bash

ENVV=$1

if [ $ENVV = 'arm64' ];then
  go build -o bin/arm64/ipicture  main.go
elif [ $ENVV = 'windows' ]; then
  CGO_ENABLED=1 GOARCH=amd64 GOOS=windows GOARCH=amd64 CC=/opt/local/bin/x86_64-w64-mingw32-gcc go build -o bin/windows/ipicture.exe -ldflags "-H windowsgui" main.go
fi
