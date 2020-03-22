#!/usr/bin/env bash
# An alternative "go" command

set -e
subcmd=$1
set -u

shift
if [[ $subcmd == "run" ]]; then
    make --silent
    ./minigo "$@" > /tmp/tmpfs/a.s
    cp /tmp/tmpfs/a.s a.s
    ./as /tmp/tmpfs/a.s
fi
