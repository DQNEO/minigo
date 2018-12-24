#!/usr/bin/env bash

cat > out/a.s

if [[ `uname` == "Darwin" ]];then
    docker run --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go bash -c 'gcc -no-pie out/a.s && ./a.out'
else
    gcc -no-pie out/a.s && ./a.out
fi

