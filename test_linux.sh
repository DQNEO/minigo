#!/bin/bash

gcc a.s && ./a.out

if [[ $? -eq 0 ]];then
    echo "ok"
fi

