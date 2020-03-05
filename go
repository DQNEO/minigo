#!/usr/bin/env bash
# An alternative "go" command

set -e
subcmd=$1
set -u

shift
if [[ $subcmd == "run" ]]; then
    make --silent
    ./minigo "$@" > /tmp/a.s
    cp /tmp/a.s a.s
    ./as /tmp/a.s
fi


