#!/bin/bash
set -eu

prog_name=minigo

function compile {
    local basename=$1
    local src=t/$basename/${basename}.go
    local as_file=out/${basename}.s
    echo -n "compile $src  > $as_file ... "
    ./${prog_name} $src > $as_file
    echo ok
}

[[ -d  ./out ]] || mkdir ./out

for testfile in t/expected/*.txt
do
    name=$(basename -s .txt $testfile)
    compile $name
done

echo "All tests passed"
