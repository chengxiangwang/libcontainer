#!/usr/bin/env bash

go build && go install

rm -rf /mnt/containers

libcontainer run -name container_test001 -ti  -v /root/backup:/mnt/backup /bin/sh
