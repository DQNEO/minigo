#!/bin/bash
set -e
PATH="/usr/lib/go-1.10/bin:$PATH"

as_file="./out/a.s"
executable="./out/a.out"
prog_name="minigo.linux"

function test_file {
    local source=$1
    local expected=$2
    ./${prog_name} $source > $as_file
    gcc -no-pie -o $executable $as_file
    $executable > out/actual.txt
    diff -u out/actual.txt $expected
}

function test_main {
    local code="$1"
    local expected=$2
    rm -f $as_file
    echo -e "
package main

func main() {
  $code
}

  " |  ./${prog_name} - > $as_file
    gcc -no-pie -o $executable $as_file
    local actual=`$executable`
    if [[ "$actual" -eq "$expected" ]];then
        echo "ok"
    else
        echo "not ok"
        exit 1
    fi

}

test_file t/test.go t/test.txt
test_file t/hello.go t/hello.txt

test_main "var i int; i = 3; printf(\"%d\",i)" 3
test_main "printf(\"%d\",1); printf(\"%d\",7)" 17
test_main 'printf("%d", 2 + 5)' 7
test_main 'printf("%d", 2 * 3)' 6
test_main 'printf("%d", 3 -2)' 1

echo "All tests passed"
