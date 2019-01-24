# Mini Go Compiler

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
root@3d19df4a8fa6:/mnt# ./minigo t/hello/hello.go  > out/a.s
root@3d19df4a8fa6:/mnt# ./as
hello world
```
