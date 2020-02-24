#!/bin/bash
set -u

differ=0

for testfile in t/expected/*.txt
do
    name=$(basename -s .txt $testfile)
    ./unit_test.sh minigo2 $name 2
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
./testerror.sh minigo2

echo "All tests passed"
