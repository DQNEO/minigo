#!/bin/bash

set -eu

test_names="
anytype-switch
anytype
arith
assign-to-indexexpr
assign
atoi
backquote
byte-cmp
cmp
const
conversion
fizzbuzz
for
forcond
funcref
heap
hello
if
incr
internal
map-of-map
map
map2
min
multi
newline
open
pointer
println
read-file
slice2
sprintf
string-concat
string-index
string
switch
test
var
write
"

for testname in $test_names
do
    ./unit_test.sh  minigo2 $testname 2
done

echo "All 2gen tests passed."
