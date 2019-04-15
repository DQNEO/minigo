#!/usr/bin/env bash

set -eu

echo "package main"
echo ""

for dir in stdlib/*
do
    basename=${dir##*/}
    echo "const ${basename}Code=\`"
    cat stdlib/$basename/$basename.go
    echo "\`"
    echo ""
done

echo "var stdPkgs map[identifier]string = map[identifier]string{"

for dir in stdlib/*
do
    basename=${dir##*/}
    echo -e "\t\"$basename\": ${basename}Code,"
done


echo "}"
echo ""
