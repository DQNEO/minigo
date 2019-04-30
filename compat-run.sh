#!/usr/bin/env bash
set -ux

cmd="$@"
if [[ `uname` == "Darwin" ]];then
    # for MacOS
    docker run -it --cap-add=SYS_PTRACE --security-opt='seccomp=unconfined' --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential bash -c "$cmd"
else
    # for Linux
    $cmd
fi




