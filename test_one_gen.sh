#!/bin/bash
set -u

differ=0
generation=$1

for testfile in t/expected/*.txt
do
    name=$(basename -s .txt $testfile)
    ./unit_test $name $generation
    if [[ $? -ne 0 ]];then
        differ=1
    fi
done

if [[ $differ -eq 0 ]];then
    :
else
    echo "FAILED"
    exit 1
fi

set -e
./terror/testerror.sh $generation

echo "All tests passed"
