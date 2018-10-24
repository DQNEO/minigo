#!/bin/bash
set -e
PATH="/usr/lib/go-1.10/bin:$PATH"

as_file="./out/a.s"
executable="./out/a.out"
prog_name="minigo.linux"

function test_file {
    local basename=$1
    local src=t/$basename/${basename}.go
    local expected=t/expected/${basename}.txt
    echo -n "test_file $src  ... "
    ./${prog_name} $src > $as_file
    gcc -no-pie -o $executable $as_file
    $executable > out/actual.txt
    diff -u out/actual.txt $expected
    echo ok
}

test_file min
test_file test
test_file hello
test_file const
test_file var
test_file byte
test_file array
test_file type

echo "All tests passed"
