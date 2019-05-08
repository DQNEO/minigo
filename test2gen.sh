#!/bin/bash

set -eu
for testname in min hello byte-cmp 'if' println
do
    ./unit_test.sh  minigo2 $testname 2
done

echo "All 2gen tests passed."
