#!/bin/bash

./minigo terror/panic/panic.go > /tmp/out/a.s

as -o /tmp/out/a.o /tmp/out/a.s && ld -o a.out /tmp/out/a.o && ./a.out >/dev/null

if [[ $? -ne 2 ]]; then
    echo "FAILED"
    exit 1
fi

echo "ok"
