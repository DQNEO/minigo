#!/bin/bash
set -e
export PATH="/usr/lib/go-1.10/bin:$PATH"


function run_case {
    local code=$1
    local expected=$2
    rm -f a.s
    echo -n "$code" | go run main.go > a.s
    gcc a.s
    ./a.out || true
    if [[ $? -eq 0 ]];then
        echo "ok"
    else
        echo "not ok"
        exit 1
    fi

}
run_case 0 0
run_case 7 7

echo "All tests passed"
