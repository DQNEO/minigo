#!/bin/bash
set -u
program=$1

mkdir -p /tmp/out

./${program} terror/panic/panic.go > /tmp/out/a.s

as -o /tmp/out/a.o /tmp/out/a.s && ld -o a.out /tmp/out/a.o && ./a.out >/dev/null

if [[ $? -eq 0 ]]; then
    echo "FAILED"
    exit 1
fi

./${program} terror/exit/exit.go > /tmp/out/a.s

as -o /tmp/out/a.o /tmp/out/a.s && ld -o a.out /tmp/out/a.o && ./a.out >/dev/null

if [[ $? -eq 0 ]]; then
    echo "FAILED"
    exit 1
fi

echo "ok"
