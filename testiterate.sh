#!/bin/bash
set -u

progname=minigo
mkdir -p /tmp/$progname

echo iterate by $progname
for testfile in t/expected/*.txt
do
    testname=$(basename -s .txt $testfile)
    echo $progname $testname
    testtarget=t/$testname/*.go
    outfile=/tmp/$progname/$testname.ast
    ./$progname --parse-only -d -a $testtarget 2> $outfile
    if [[ $? -ne 0 ]]; then
        cat $outfile
        echo test $testname failed
        exit 1
    fi
done

progname=minigo2
mkdir -p /tmp/$progname

for testfile in t/expected/*.txt
do
    testname=$(basename -s .txt $testfile)
    echo $progname $testname
    testtarget=t/$testname/*.go
    outfile=/tmp/$progname/$testname.ast
    ./$progname --parse-only -d -a $testtarget 2> $outfile
    if [[ $? -ne 0 ]]; then
        echo "[FAIL]" $progname $testname
    fi
done


echo "All iteration done"
