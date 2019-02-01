#!/bin/bash
set -e

actual=out/actual.txt

function test_file {
    local basename=$1
    local expected=t/expected/${basename}.txt
    local bin_file=out/${basename}.bin
    local as_file=out/${basename}.s
    local obj_file=out/${basename}.o
    rm -f $actual
    echo -n "test_file $as_file  ... "
    as -o $obj_file $as_file
    # gave up direct invocation of "ld"
    # https://stackoverflow.com/questions/33970159/bash-a-out-no-such-file-or-directory-on-running-executable-produced-by-ld
    gcc -no-pie -o $bin_file $obj_file
    $bin_file /etc/lsb-release > $actual
    diff -u $expected $actual
    echo ok
}

[[ -d  ./out ]] || mkdir ./out

for testfile in t/expected/*.txt
do
    name=$(basename -s .txt $testfile)
    test_file $name
done

echo "All tests passed"
