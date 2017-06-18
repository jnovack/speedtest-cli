#!/bin/bash

go build -v -o speedtest-cli-jose speedtest.go
scp speedtest-cli-jose joseleon@10.0.0.35:/tmp

