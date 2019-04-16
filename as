#!/usr/bin/env bash
set -ux
file=out/a.s

[[ -e $file ]] || echo "$file ooes not exist"

cmd="gcc -g -no-pie $file && ./a.out || gdb --batch --eval-command=run ./a.out"
if [[ `uname` == "Darwin" ]];then
    # for MacOS
    docker run -it --cap-add=SYS_PTRACE --security-opt='seccomp=unconfined' --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential bash -c "$cmd"
else
    # for Linux
    gcc -g -no-pie $file && ./a.out || gdb --batch --eval-command=run ./a.out
fi

