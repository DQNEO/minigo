#!/usr/bin/env bash

set -eu

echo "package main"
echo ""

for dir in stdlib/*
do
    basename=${dir##*/}
    echo "const ${basename}Code string = \`"
    cat stdlib/$basename/$basename.go
    echo "\`"
    echo ""
done

echo "func makeStdLib() map[packageName]string {"
echo "    var mp map[packageName]string = map[packageName]string{"

for dir in stdlib/*
do
    basename=${dir##*/}
    echo -e "        \"$basename\": ${basename}Code,"
done


echo "    }"
echo "    return mp"
echo "}"
echo ""
