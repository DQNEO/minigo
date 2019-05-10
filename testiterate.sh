#!/bin/bash
set -u

function iterate() {
    local progname=$1
    mkdir -p /tmp/$progname
    echo iterate by $progname
    for testfile in t/expected/*.txt
    do
        testname=$(basename -s .txt $testfile)
        echo $progname $testname
        testtarget=t/$testname/*.go
        outfile=/tmp/$progname/$testname.token
        ./$progname --parse-only -d -a $testtarget 2> $outfile
        if [[ $? -ne 0 ]]; then
            #cat $outfile
            echo test $testname failed
            exit 1
        fi
    done
}

iterate minigo
iterate minigo2

echo "All iteration done"
