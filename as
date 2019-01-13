#!/usr/bin/env bash
set -eu
file=out/a.s

[[ -e $file ]] || echo "$file ooes not exist"

if [[ `uname` == "Darwin" ]];then
    # for MacOS
    docker run --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go bash -c "gcc -no-pie $file && ./a.out"
else
    # for Linux
    gcc -no-pie $file && ./a.out
fi

