# Minigo

[![CircleCI](https://circleci.com/gh/DQNEO/minigo.svg?style=svg)](https://circleci.com/gh/DQNEO/minigo)

A Go compiler from scratch.

# Description
`minigo` is yet another Go compiler made from scratch.

* No dependency on yacc/lex
* No dependency on external libraries
* Standard libraries are also made from scratch.

It depends only on gcc as an assenmbler and linker, and on libc as a runtime.

# Notes

The design is influenced from

* 8cc (https://github.com/rui314/8cc)
* 8cc.go (https://github.com/DQNEO/8cc.go)

# How to run

Currently the generated code can run only on Linux.
So I would recommend you to use Docker.

```
$ docker run --rm -it -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go bash
```

After entering the container, you can build and run it.

```
# make
# ./minigo t/hello/hello.go > a.s
# gcc -g -no-pie a.s
# ./a.out
# hello world
```

# How to do "self compile"

```
# make
# ./minigo --version
minigo 0.1.0
Copyright (C) 2019 @DQNEO

# ./minigo *.go > /tmp/minigo2.s
# gcc -no-pie -o minigo2 /tmp/minigo2.s
# ./minigo2 --version
minigo 0.1.0
Copyright (C) 2019 @DQNEO

```

You will see that the contents of 2nd generation compiler and 3rd generation compiler are identical.

```
# ./minigo2 *.go > /tmp/minigo3.s
# diff /tmp/minigo2.s /tmp/minigo3.s
```

# AUTHOR
[@DQNEO](https://twitter.com/DQNEO)

# LICENSE

MIT License
