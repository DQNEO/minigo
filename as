#!/usr/bin/env bash
file=$1
set -u

if [[ $1 == "--verbose" ]];then
  set -x
fi

[[ $file == ""  ]] &&  { echo "not input" ; exit 1 ;}
[[ -e $file ]] || { echo "file not found: $file"; exit 1; }

set -e
as -o a.o $file
ld -o a.out a.o
./a.out
#|| gdb --batch --eval-command=run ./a.out

