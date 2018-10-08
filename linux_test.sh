#!/bin/bash
set -ex
PATH="/usr/lib/go-1.10/bin:$PATH"

as_file="./out/a.s"
executable="./out/a.out"
prog_name="minigo"
go build -o minigo *.go

function run_test_go {
    ./${prog_name} < t/test.go > $as_file
    gcc -no-pie -o $executable $as_file
    $executable > out/actual.txt
    diff out/actual.txt t/expected.txt
}

run_test_go

echo "All tests passed"
