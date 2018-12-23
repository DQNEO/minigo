#!/bin/bash
set -e
PATH="/usr/lib/go-1.10/bin:$PATH"

as_file="./out/a.s"
executable="./out/a.out"
prog_name="minigo.linux"
actual=out/actual.txt

function do_test {
    ./${prog_name} $src > $as_file
    gcc -no-pie -o $executable $as_file
    $executable > $actual
    diff -u $actual $expected
}

function test_file {
    local basename=$1
    local src=t/$basename/${basename}.go
    local expected=t/expected/${basename}.txt
    rm -f $actual
    echo -n "test_file $src  ... "
    do_test ./${prog_name} $src
    echo ok
}

[[ -d  ./out ]] || mkdir ./out

for testfile in t/expected/*.txt
do
    name=$(basename -s .txt $testfile)
    test_file $name
done

echo "All tests passed"
