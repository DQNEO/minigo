# Mini Go Compiler

[![CircleCI](https://circleci.com/gh/DQNEO/minigo.svg?style=svg)](https://circleci.com/gh/DQNEO/minigo)

A Go compiler from scratch.

# Description
`minigo` is yet another Go compiler made from scratch.

* No dependency on yacc/lex
* No dependency on external packages/libraries

Lexer and parse are written by hand.
Standard libraries are also made from scratch.

It depends only on gcc as an assenmbler and linker, and on libc as a runtime.

# Notes

The design is influenced from

* 8cc (https://github.com/rui314/8cc)
* 8cc.go (https://github.com/DQNEO/8cc.go)

# AUTHOR
[@DQNEO](https://twitter.com/DQNEO)

# LICENSE

MIT License

# How to run

Currently the generated code can run only on Linux.
So I would recommend you to use Docker.

```
$ docker run --rm -it -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go bash
```

After entering the container, you can make and run it.

```
root@da1ca5837dac:/mnt# make
root@da1ca5837dac:/mnt# ./minigo t/hello/hello.go > a.s
root@da1ca5837dac:/mnt# gcc -g -no-pie a.s
root@da1ca5837dac:/mnt# ./a.out
hello world
```

# How to do "self compile"

```
root@ff3a028b3564:/mnt# make
root@ff3a028b3564:/mnt# ./minigo --version
minigo 0.1.0
Copyright (C) 2019 @DQNEO
root@ff3a028b3564:/mnt# ./minigo *.go > /tmp/minigo2.s
root@ff3a028b3564:/mnt# gcc -no-pie -o minigo2 /tmp/minigo2.s
root@ff3a028b3564:/mnt# ./minigo2 --version
minigo 0.1.0
Copyright (C) 2019 @DQNEO
```
