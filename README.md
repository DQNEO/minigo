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
root@40823f69850c:/mnt# make
root@40823f69850c:/mnt# ./minigo.linux t/hello.go |./as
hello world
```
