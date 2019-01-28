#!/usr/bin/env bash
set -eux
file=out/a.s

[[ -e $file ]] || echo "$file ooes not exist"

if [[ `uname` == "Darwin" ]];then
    # for MacOS
    docker run --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential bash -c "gcc -g -no-pie $file && ./a.out"
else
    # for Linux
    gcc -g -no-pie $file && ./a.out
fi

