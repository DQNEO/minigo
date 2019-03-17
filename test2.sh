#!/bin/bash

./minigo t2/panic/panic.go > out/a.s

if [[ `uname` == "Darwin" ]];then
    # for MacOS
    docker run --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential bash -c 'gcc -g -no-pie out/a.s && ./a.out >/dev/null'
else
    # for Linux
    gcc -g -no-pie out/a.s && ./a.out >/dev/null
fi

if [[ $? -ne 1 ]]; then
    echo "FAILED"
    exit 1
fi

echo "ok"
