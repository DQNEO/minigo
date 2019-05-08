# All the commands are supposed to run on Linux.
# I use Docker like below.
# docker run -it --cap-add=SYS_PTRACE --security-opt='seccomp=unconfined' --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go bash

all: minigo /tmp/out

/tmp/out:
	mkdir /tmp/out

internalcode.go: internalcode/runtime.go
	./cp-internalcode.sh

stdlib.go: stdlib/*/*.go
	./concate-stdlib.sh > stdlib.go

# 1st gen compiler
minigo: *.go internalcode.go stdlib.go
	go build -o minigo *.go


test1gen: all
	./test1gen.sh

# 2nd gen assembly
/tmp/out/minigo.s: *.go minigo
	./minigo *.go > /tmp/out/minigo.s

# 2nd gen compiler
minigo2: /tmp/out/minigo.s
	gcc -g -no-pie -o minigo2 /tmp/out/minigo.s


test2gen: minigo2
	./unit_test.sh  minigo2 min 2
	./unit_test.sh  minigo2 hello 2

test: all
	make test1gen
	make test2gen

clean:
	rm -f minigo*
	rm -f a.s a.out
	rm -f /tmp/out/*
	rm -f stdlib.go
	rm -f internalcode.go

fmt:
	gofmt -w *.go t/*/*.go
