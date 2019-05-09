#!/bin/bash
set -u

progname=minigo

mkdir -p /tmp/$progname

for testfile in t/expected/*.txt
do
    testname=$(basename -s .txt $testfile)
    testtarget=t/$testname/*.go
    outfile=/tmp/$progname/$testname.ast
    ./$progname --parse-only -d -a $testtarget 2> $outfile
    if [[ $? -ne 0 ]]; then
        cat $outfile
        echo test $testname failed
        exit 1
    fi
done

echo "All tests passed"
