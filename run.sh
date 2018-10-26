#!/usr/bin/env bash

go build && go install

rm -rf /mnt/containers

libcontainer run   -v /root/backup:/mnt/backup top
