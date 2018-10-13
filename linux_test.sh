#!/bin/bash
set -ex
PATH="/usr/lib/go-1.10/bin:$PATH"

as_file="./out/a.s"
executable="./out/a.out"
prog_name="minigo"

make clean
make

function test_file {
    ./${prog_name} < t/test.go > $as_file
    gcc -no-pie -o $executable $as_file
    $executable > out/actual.txt
    diff out/actual.txt t/expected.txt
}

function test_expr {
    local code="$1"
    local expected=$2
    rm -f $as_file
    echo -e "
package main

func main() {
  $code
}

  " |  ./${prog_name} > $as_file
    gcc -no-pie -o $executable $as_file
    local actual=`$executable`
    if [[ "$actual" -eq "$expected" ]];then
        echo "ok"
    else
        echo "not ok"
        exit 1
    fi

}

test_file

test_expr "var i int\ni = 3\nprintf(\"%d\",i)" 3
test_expr "printf(\"%d\",1)\nprintf(\"%d\",7)" 17
test_expr 'printf("%d", 2 + 5)' 7
test_expr 'printf("%d", 2 * 3)' 6
test_expr 'printf("%d", 3 -2)' 1

echo "All tests passed"
