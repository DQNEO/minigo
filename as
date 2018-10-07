#!/usr/bin/env bash

cat > out/a.s
docker run -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go bash -c 'gcc out/a.s && ./a.out'

