#!/bin/bash
set -u
generation=$1
progname=""
if [[ $generation == "1" ]];then
  progname="minigo"
elif [[ $generation == "2" ]]; then
  progname="minigo2"
else
  echo "Invalid argument"
  exit 1
fi

mkdir -p /tmp/out

./${progname} terror/panic/panic.go > /tmp/out/a.s

as -o /tmp/out/a.o /tmp/out/a.s && ld -o a.out /tmp/out/a.o && ./a.out >/dev/null

if [[ $? -eq 0 ]]; then
    echo "FAILED"
    exit 1
fi

./${progname} terror/exit/exit.go > /tmp/out/a.s

as -o /tmp/out/a.o /tmp/out/a.s && ld -o a.out /tmp/out/a.o && ./a.out >/dev/null

if [[ $? -eq 0 ]]; then
    echo "FAILED"
    exit 1
fi

echo "ok"
