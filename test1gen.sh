#!/bin/bash
set -u

prog_name=minigo
out_dir=/tmp/out
actual=$out_dir/actual.txt
# for os.Args
ARGS=t/data/sample.txt
differ=0

function compile {
    local basename=$1
    local src=t/$basename/*.go
    local as_file=$out_dir/${basename}.s
    echo -n "compile $src  > $as_file ... "
    ./${prog_name} $src > $as_file
    echo ok
}


function as_run {
    local basename=$1
    local expected=t/expected/${basename}.txt
    local bin_file=$out_dir/${basename}.bin
    local as_file=$out_dir/${basename}.s
    local obj_file=$out_dir/${basename}.o
    rm -f $actual
    echo -n "as_run $as_file  ... "
    as -o $obj_file $as_file
    # gave up direct invocation of "ld"
    # https://stackoverflow.com/questions/33970159/bash-a-out-no-such-file-or-directory-on-running-executable-produced-by-ld
    gcc -no-pie -o $bin_file $obj_file
    $bin_file $ARGS > $actual
    diff -u $expected $actual
    if [[ $? -ne 0 ]];then
        differ=1
    fi
    echo ok
}

function run_unit_test {
    local ame=$1
    compile $name
    as_run $name
}

for testfile in t/expected/*.txt
do
    name=$(basename -s .txt $testfile)
    run_unit_test $name
done

if [[ $differ -eq 0 ]];then
    echo "All tests passed"
else
    echo "FAILED"
    exit 1
fi

set -e
./testerror.sh

echo "All tests passed"
