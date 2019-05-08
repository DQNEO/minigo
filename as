#!/usr/bin/env bash
file=$1
set -ux

[[ $file == ""  ]] &&  { echo "not input" ; exit 1 ;}
[[ -e $file ]] || { echo "file not found: $file"; exit 1; }

cmd="gcc -g -no-pie $file && ./a.out " # || gdb --batch --eval-command=run ./a.out"
if [[ `uname` == "Darwin" ]];then
    # for MacOS
    docker run -it --cap-add=SYS_PTRACE --security-opt='seccomp=unconfined' --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential bash -c "$cmd"
else
    # for Linux
    gcc -g -no-pie $file && ./a.out || gdb --batch --eval-command=run ./a.out
fi

