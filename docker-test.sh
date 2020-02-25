#!/bin/bash
# run `make test` on docker
docker run \
    -it\
    --rm\
    --cap-add=SYS_PTRACE\
    --security-opt='seccomp=unconfined'\
    -w /mnt\
    -v `pwd`:/mnt\
    --env PATH=/usr/lib/go-1.14/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\
    --tmpfs=/tmp/tmpfs:rw,size=500m,mode=1777\
     dqneo/ubuntu-build-essential:go make test


