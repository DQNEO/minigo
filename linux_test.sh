#!/bin/bash
set -e
PATH="/usr/lib/go-1.10/bin:$PATH"

prog_name=minigo.linux

as_file=out/a.s
obj_file=out/a.o
bin_file=out/a.out
actual=out/actual.txt

function test_file {
    local basename=$1
    local src=t/$basename/${basename}.go
    local expected=t/expected/${basename}.txt
    rm -f $actual
    echo -n "test_file $src  ... "
    ./${prog_name} $src > $as_file
    as -o $obj_file $as_file
    # gave up direct invocation of "ld"
    # https://stackoverflow.com/questions/33970159/bash-a-out-no-such-file-or-directory-on-running-executable-produced-by-ld
    gcc -no-pie -o $bin_file $obj_file
    $bin_file > $actual
    diff -u $actual $expected
    echo ok
}

[[ -d  ./out ]] || mkdir ./out

for testfile in t/expected/*.txt
do
    name=$(basename -s .txt $testfile)
    test_file $name
done

echo "All tests passed"
