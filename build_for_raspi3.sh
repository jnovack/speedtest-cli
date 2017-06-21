#!/bin/bash

if [ -d target ]; then
   rm -rf target
fi
mkdir -p target

GOOS=linux GOARCH=arm GOARM=7 go build -v -o target/speedtest-cli-raspi3 speedtest.go
#scp speedtest-cli-raspi3 pi@10.0.0.38:/tmp

