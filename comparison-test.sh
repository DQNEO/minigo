#!/bin/bash

# test tokenizer
make

mkdir -p /tmp/minigo2 /tmp/minigo

# Compare toknizer output
for f in *.go
do
    echo -n "  tokenizing $f ...  "
    ./minigo  --tokenize-only -d -t $f 2> /tmp/${f}.1.token
    ./minigo2 --tokenize-only -d -t $f 2> /tmp/${f}.2.token

    diff /tmp/${f}.1.token /tmp/${f}.2.token || exit 1
    echo "ok"
done
echo "tokinzer ok"

# Compare AST output
for f in  *.go
do
    echo -n "  parsing $f ...  "
    ./minigo  --parse-only -d -a $f 2> /tmp/${f}.1.ast
    ./minigo2 --parse-only -d -a $f 2> /tmp/${f}.2.ast

    diff /tmp/${f}.1.ast /tmp/${f}.2.ast || exit 1
    echo "ok"
done

echo "parser ok"
