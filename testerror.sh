#!/bin/bash

./minigo terror/panic/panic.go > /tmp/out/a.s

gcc -nostdlib -g -no-pie /tmp/out/a.s && ./a.out >/dev/null

if [[ $? -ne 1 ]]; then
    echo "FAILED"
    exit 1
fi

echo "ok"
