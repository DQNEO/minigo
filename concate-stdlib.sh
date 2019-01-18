#!/usr/bin/env bash

set -eu

echo "package main"
echo ""

for dir in stdlib/*
do
    basename=${dir##*/}
    code=$(cat stdlib/$basename/$basename.go)
    echo -e "const ${basename}Code=\`\n$code\n\`"
    echo ""
done

echo "var pkgMap map[identifier]string = map[identifier]string{"

for dir in stdlib/*
do
    basename=${dir##*/}
    echo -e "\t\"$basename\": ${basename}Code,"
done


echo "}"
echo ""
