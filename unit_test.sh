#!/bin/bash

progname=$1
basename=$2
suffix=$3

out_dir=/tmp/out

if [[ -z $suffix ]];then
    as_file=$out_dir/${basename}.s
else
    as_file=$out_dir/${basename}.${suffix}.s
fi
src=t/$basename/*.go
expected=t/expected/${basename}.txt
bin_file=$out_dir/${basename}.bin
obj_file=$out_dir/${basename}.o

actual=$out_dir/actual.txt
# for os.Args
ARGS=t/data/sample.txt

function compile {
    echo -n "compile $src  > $as_file ... "
    ./${progname} $src > $as_file
    echo ok
}

function as_run {
    rm -f $actual
    echo -n "as_run $as_file  ... "
    as -o $obj_file $as_file
    # gave up direct invocation of "ld"
    # https://stackoverflow.com/questions/33970159/bash-a-out-no-such-file-or-directory-on-running-executable-produced-by-ld
    gcc -no-pie -o $bin_file $obj_file
    $bin_file $ARGS > $actual
    diff -u $expected $actual
}

function run_unit_test {
    compile
    as_run
    if [[ $? -ne 0 ]];then
        echo failed
        return 1
    else
        echo ok
        return 0
    fi
}

run_unit_test
