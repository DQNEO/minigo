#!/bin/bash

set -e
# test tokenizer
make

mkdir -p /tmp/minigo2 /tmp/minigo

for f in main.go gen.go parser.go
do
    ./minigo  --tokenize-only -d -t $f 2> /tmp/${f}.1.token
    ./minigo2 --tokenize-only -d -t $f 2> /tmp/${f}.2.token

    diff /tmp/${f}.1.token /tmp/${f}.2.token
done

echo ok