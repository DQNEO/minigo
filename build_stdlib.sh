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

echo "var pkgsources []pkgsource = []pkgsource{"
for dir in stdlib/*
do
    basename=${dir##*/}
    echo -e "\tpkgsource{"
    echo -e "\t\tname: \"$basename\","
    echo -e "\t\tcode: ${basename}Code,"
    echo -e "\t},"
done


echo "}"
echo ""
