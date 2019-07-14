#!/usr/bin/env bash

set -eu

echo "package main"
echo ""

for dir in stdlib/*
do
    basename=${dir##*/}
    echo "var ${basename}Code string = string(\`"
    cat stdlib/$basename/$basename.go
    echo "\`)"
    echo ""
done

echo "func makeStdLib() map[identifier]string {"
echo "    var mp map[identifier]string = map[identifier]string{"

for dir in stdlib/*
do
    basename=${dir##*/}
    echo -e "        identifier(\"$basename\"): ${basename}Code,"
done


echo "    }"
echo "    return mp"
echo "}"
echo ""
