#!/usr/bin/env bash

set -eu

echo "package main"
echo ""

for dir in stdlib/*
do
    basename=${dir##*/}
    echo "var ${basename}Code bytes = bytes(\`"
    cat stdlib/$basename/$basename.go
    echo "\`)"
    echo ""
done

echo "func makeStdLib() map[identifier]bytes {"
echo "    var mp map[identifier]bytes = map[identifier]bytes{"

for dir in stdlib/*
do
    basename=${dir##*/}
    echo -e "        identifier(\"$basename\"): ${basename}Code,"
done


echo "    }"
echo "    return mp"
echo "}"
echo ""
