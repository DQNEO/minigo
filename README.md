# Minigo🐥

[![CircleCI](https://circleci.com/gh/DQNEO/minigo.svg?style=svg)](https://circleci.com/gh/DQNEO/minigo)

![Go](https://github.com/DQNEO/minigo/workflows/Go/badge.svg)

A Go compiler from scratch.

# Description

`minigo🐥` is a small Go compiler made from scratch. It can compile itself.

* Generates a single static  binary executable
* No dependency on yacc/lex or any external libraries
* Standard libraries are also made from scratch

It depends only on GNU Assembler and GNU ld.

`minigo` supports x86-64 Linux only.
 
# Design

I made this almost without reading the original Go compiler.

`minigo` inherits most of its design from the followings.

* 8cc (https://github.com/rui314/8cc)
* 8cc.go (https://github.com/DQNEO/8cc.go)

There are several steps in the compilation proccess.

[go source] -> byte_stream.go -> [byte stream] -> token.go -> [token stream] -> parser.go -> [AST] -> gen.go -> [assembly code]


# How to run

You need Linux, so I would recommend that you use Docker.

```sh
$ docker run --rm -it -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go bash
```

After entering the container, you can build and run it.

```sh
$ make
$ ./minigo t/hello/hello.go > hello.s
$ as -o hello.o hello.s
$ ld -o hello hello.o
$ ./hello
hello world
```

# How to "self compile"

```sh
$ make
$ ./minigo --version
minigo 0.1.0
Copyright (C) 2019 @DQNEO

$ ./minigo *.go > /tmp/minigo2.s
$ as -o /tmp/minigo2.o /tmp/minigo2.s
$ ld -o minigo2 /tmp/minigo2.o
$ ./minigo2 --version
minigo 0.1.0
Copyright (C) 2019 @DQNEO

$ ./minigo2 *.go > /tmp/minigo3.s
$ as -o /tmp/minigo3.o /tmp/minigo3.s
$ ld -o minigo3 /tmp/minigo3.o
$ ./minigo3 --version
minigo 0.1.0
Copyright (C) 2019 @DQNEO
```

You will see that the contents of 2nd generation compiler and 3rd generation compiler are identical.

```sh
$ diff /tmp/minigo2.s /tmp/minigo3.s
```

# Test

```sh
$ make test
```

# Debug by gdb

Add `--cap-add=SYS_PTRACE --security-opt='seccomp=unconfined'` option to `docker run`.
It will allow you to use `gdb` in the docker image.

```
docker run --cap-add=SYS_PTRACE --security-opt='seccomp=unconfined' -it --rm -w /mnt -v `pwd`:/mnt --tmpfs=/tmp/tmpfs:rw,size=500m,mode=1777 dqneo/ubuntu-build-essential:go bash
```

# AUTHOR

[@DQNEO](https://twitter.com/DQNEO)

# LICENSE

MIT License
