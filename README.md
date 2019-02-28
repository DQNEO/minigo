# Mini Go Compiler

[![CircleCI](https://circleci.com/gh/DQNEO/minigo.svg?style=svg)](https://circleci.com/gh/DQNEO/minigo)

A Go compiler from scratch.

# Description
`minigo` is yet another Go compiler made from scratch.

* No dependency on external packages/libraries
* No dependency on yacc/lex things

Lexer and parse are written by hand.
Standard libraries are also made from scratch.

# Notes

The design is influenced from

* 8cc (https://github.com/rui314/8cc)
* 8cc.go (https://github.com/DQNEO/8cc.go)

# AUTHOR
[@DQNEO](https://twitter.com/DQNEO)

# LICENSE

MIT License

# How to run

Currently the generated code can only run on Linux.
So you need Docker to run it.

```
$ docker run --rm -it -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go bash
```

After entering container, you can make and run it.

```
root@da1ca5837dac:/mnt# make
root@da1ca5837dac:/mnt# ./minigo t/hello/hello.go > a.s
root@da1ca5837dac:/mnt# gcc -g -no-pie a.s
root@da1ca5837dac:/mnt# ./a.out
hello world
```
