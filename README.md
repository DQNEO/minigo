# Mini Go Compiler

[![CircleCI](https://circleci.com/gh/DQNEO/minigo.svg?style=svg)](https://circleci.com/gh/DQNEO/minigo)

minigo is a Go compiler from scratch.

The design is much influenced from 

* 8cc (https://github.com/rui314/8cc)
* 8cc.go (https://github.com/DQNEO/8cc.go)

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
