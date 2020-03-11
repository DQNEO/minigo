#!/bin/bash
set -u
generation=$1
progname=""
if [[ $generation == "1" ]];then
  progname="./minigo"
elif [[ $generation == "2" ]]; then
  progname="./minigo2"
elif [[ $generation == "0" ]]; then
  progname="go run"
else
  echo "Invalid argument for " $0
  exit 1
fi

mkdir -p /tmp/out

${progname} terror/panic/panic.go > /tmp/out/a.s 2>/dev/null
status=$?
if [[ $progname != "go run" ]]; then
  as -o /tmp/out/a.o /tmp/out/a.s && ld -o a.out /tmp/out/a.o && ./a.out >/dev/null
  status=$?
fi
if [[ $status -eq 0 ]]; then
    echo "FAILED"
    exit 1
fi

${progname} terror/exit/exit.go > /tmp/out/a.s 2>/dev/null
status=$?
if [[ $progname != "go run" ]]; then
  as -o /tmp/out/a.o /tmp/out/a.s && ld -o a.out /tmp/out/a.o && ./a.out >/dev/null
  status=$?
fi

if [[ $status -eq 0 ]]; then
    echo "FAILED"
    exit 1
fi

echo "ok"
