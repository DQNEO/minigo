#!/bin/bash
set -eu

prog_name=minigo

function compile {
    local basename=$1
    local src=t/$basename/*.go
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

if [[ `uname` == "Darwin" ]];then
    # for MacOS
    docker run --cap-add=SYS_PTRACE --security-opt='seccomp=unconfined' --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential ./linux_test.sh
else
    # for Linux
    ./linux_test.sh
fi

./testerror.sh

echo "All tests passed"
