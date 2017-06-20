#!/bin/bash

GOOS=linux GOARCH=arm GOARM=5 go build -v -o speedtest-cli-raspi speedtest.go
scp speedtest-cli-raspi pi@10.0.0.39:/tmp

