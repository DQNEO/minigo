#!/bin/bash
set -ex
PATH="/usr/lib/go-1.10/bin:$PATH"

as_file="./out/a.s"

prog_name="minigo"
go build -o minigo *.go

function run_test_go {
    ./${prog_name} < test/test.go > $as_file
    gcc -no-pie -o out/a.out $as_file
    ./out/a.out > out/actual.txt
    diff out/actual.txt test/expected.txt
}

function run_case {
    local code="$1"
    local expected=$2
    rm -f $as_file
    echo  "$code" | ./${prog_name} > $as_file
    gcc -no-pie -o out/a.out $as_file
    local actual=`./out/a.out`
    if [[ "$actual" -eq "$expected" ]];then
        echo "ok"
    else
        echo "not ok"
        exit 1
    fi

}

run_test_go

run_case "printf(\"hello\")
printf(\"world\")" "helloworld"
run_case 'printf("%d",0)' 0
run_case 'printf("%d",7)' 7
run_case 'printf("%d", 2 + 5)' 7
run_case 'printf("%d", 2 * 3)' 6
run_case 'printf("%d", 3 -2)' 1

echo "All tests passed"
