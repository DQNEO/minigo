#!/bin/bash
# run `make test` on docker
docker run \
    -it\
    --rm\
    --cap-add=SYS_PTRACE\
    --security-opt='seccomp=unconfined'\
    -w /mnt\
    -v `pwd`:/mnt\
    --env PATH=/usr/lib/go-1.10/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\
    dqneo/ubuntu-build-essential:go make test


