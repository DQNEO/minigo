#!/bin/bash

set -eu

test_names="
anytype-switch
anytype
append-int
argv
arith
assign-to-indexexpr
assign
atoi
backquote
byte-cmp
byte
cmp
const
conversion
fizzbuzz
for
forcond
forrangeshort
func
funcref
global-array-string
global-indirection
heap
hello
if
incr
internal
len
map-of-map
map
map2
min
multi
multireturn
newline
open-read
open
os
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
type
var
write
"

for testname in $test_names
do
    ./unit_test.sh  minigo2 $testname 2
done

echo "All 2gen tests passed."
