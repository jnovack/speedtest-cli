#!/bin/bash

GOOS=linux GOARCH=arm GOARM=7 go build -v -o speedtest-cli-raspi3 speedtest.go
scp speedtest-cli-raspi3 pi@10.0.0.38:/tmp

