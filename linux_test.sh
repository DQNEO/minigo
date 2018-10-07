#!/bin/bash
set -ex
export PATH="/usr/lib/go-1.10/bin:$PATH"

as_file="./out/a.s"

function run_case {
    local code="$1"
    local expected=$2
    rm -f $as_file
    echo  "$code" | go run main.go > $as_file
    gcc $as_file
    local actual=`./a.out ; echo $?`
    if [[ actual -eq $expected ]];then
        echo "ok"
    else
        echo "not ok"
        exit 1
    fi

}
run_case 0 0
run_case 7 7
run_case '2 + 5' 7
run_case '2 * 3' 6
run_case '3 -2' 1

echo "All tests passed"
