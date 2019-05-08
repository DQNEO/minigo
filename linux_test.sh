#!/bin/bash

actual=/tmp/out/actual.txt
differ=0

# for os.Args
sample_file=t/data/sample.txt
function test_file {
    local basename=$1
    local expected=t/expected/${basename}.txt
    local bin_file=/tmp/out/${basename}.bin
    local as_file=/tmp/out/${basename}.s
    local obj_file=/tmp/out/${basename}.o
    rm -f $actual
    echo -n "test_file $as_file  ... "
    as -o $obj_file $as_file
    # gave up direct invocation of "ld"
    # https://stackoverflow.com/questions/33970159/bash-a-out-no-such-file-or-directory-on-running-executable-produced-by-ld
    gcc -no-pie -o $bin_file $obj_file
    $bin_file $sample_file > $actual
    diff -u $expected $actual
    if [[ $? -ne 0 ]];then
        differ=1
    fi
    echo ok
}

for testfile in t/expected/*.txt
do
    name=$(basename -s .txt $testfile)
    test_file $name
done

if [[ $differ -eq 0 ]];then
    echo "All tests passed"
else
    echo "FAILED"
    exit 1
fi
