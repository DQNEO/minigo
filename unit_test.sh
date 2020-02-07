#!/bin/bash

set -e
progname=$1
basename=$2
suffix=$3

out_dir=/tmp/out
out_dir_as=/tmp/tmpfs/out
mkdir -p $out_dir $out_dir_as

if [[ -z $suffix ]];then
    as_file=$out_dir_as/${basename}.s
else
    as_file=$out_dir_as/${basename}.${suffix}.s
fi
src=t/$basename/*.go
expected=t/expected/${basename}.txt
bin_file=$out_dir/${basename}.bin
obj_file=$out_dir_as/${basename}.o

actual=$out_dir/actual.txt
# for os.Args
ARGS=t/data/sample.txt

function compile {
    ./${progname} $src > $as_file
}

function as_run {
    rm -f $actual
    as -o $obj_file $as_file
    # gave up direct invocation of "ld"
    # https://stackoverflow.com/questions/33970159/bash-a-out-no-such-file-or-directory-on-running-executable-produced-by-ld
    gcc -nostdlib -no-pie -o $bin_file $obj_file
    $bin_file $ARGS > $actual
    diff -u $expected $actual
}

function run_unit_test {
    echo -n "./unit_test.sh $progname $basename ... "
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
