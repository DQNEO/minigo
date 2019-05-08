#!/bin/bash
set -eu

prog_name=minigo

function compile {
    local basename=$1
    local src=t/$basename/*.go
    local as_file=/tmp/out/${basename}.s
    echo -n "compile $src  > $as_file ... "
    ./${prog_name} $src > $as_file
    echo ok
}

for testfile in t/expected/*.txt
do
    name=$(basename -s .txt $testfile)
    compile $name
done

./linux_test.sh

./testerror.sh

echo "All tests passed"
