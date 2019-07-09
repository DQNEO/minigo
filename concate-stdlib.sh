#!/usr/bin/env bash

set -eu

echo "package main"
echo ""

for dir in stdlib/*
do
    basename=${dir##*/}
    echo "var ${basename}Code gostring = gostring(\`"
    cat stdlib/$basename/$basename.go
    echo "\`)"
    echo ""
done

echo "func makeStdLib() map[identifier]gostring {"
echo "    var mp map[identifier]gostring = map[identifier]gostring{"

for dir in stdlib/*
do
    basename=${dir##*/}
    echo -e "        \"$basename\": ${basename}Code,"
done


echo "    }"
echo "    return mp"
echo "}"
echo ""
