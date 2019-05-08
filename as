#!/usr/bin/env bash
file=$1
set -ux

[[ $file == ""  ]] &&  { echo "not input" ; exit 1 ;}
[[ -e $file ]] || { echo "file not found: $file"; exit 1; }

gcc -g -no-pie $file && ./a.out || gdb --batch --eval-command=run ./a.out

