#!/usr/bin/env bash
file=$1
set -ux

[[ $file == ""  ]] &&  { echo "not input" ; exit 1 ;}
[[ -e $file ]] || { echo "file not found: $file"; exit 1; }

as -o a.o $file && ld -o a.out a.o && ./a.out || gdb --batch --eval-command=run ./a.out

