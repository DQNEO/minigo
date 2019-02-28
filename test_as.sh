#!/usr/bin/env bash

if [[ `uname` == "Darwin" ]];then
    # for MacOS
    docker run --cap-add=SYS_PTRACE --security-opt='seccomp=unconfined' --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential ./linux_test.sh
else
    # for Linux
    ./linux_test.sh
fi
