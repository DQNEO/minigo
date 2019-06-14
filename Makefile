# All the commands are supposed to run on Linux.
# I use Docker like below.
# docker run -it --cap-add=SYS_PTRACE --security-opt='seccomp=unconfined' --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go bash

all: minigo

internal_runtime.go: internal/runtime/*.go
	./cp-internalcode.sh

internal_universe.go: internal/universe/*.go
	./cp-internalcode.sh

stdlib.go: stdlib/*/*.go
	./concate-stdlib.sh > stdlib.go

# 1st gen compiler
minigo: *.go internal_runtime.go internal_universe.go stdlib.go
	go build -o minigo *.go

# 2nd gen assembly
minigo.s: minigo
	./minigo --position *.go > /tmp/minigo.s
	cp /tmp/minigo.s minigo.s

# 2nd gen compiler
minigo2: minigo.s
	gcc -g -no-pie -o minigo2 minigo.s

minigo2.s: minigo2 minigo *.go
	./minigo2 *.go > /tmp/minigo2.s
	cp /tmp/minigo2.s minigo2.s

selfhost: minigo2.s
	sed -e 's|^/\*.*)\*/||' minigo.s > /tmp/minigo.stripped.s
	diff minigo2.s /tmp/minigo.stripped.s && echo ok

test: minigo minigo2
	make vet
	./test1gen.sh
	./test2gen.sh
	./comparison-test.sh

clean:
	rm -f minigo minigo2
	rm -f minigo*.s
	rm -f a.s a.out
	rm -f /tmp/out/*
	rm -rf /tmp/minigo*
	rm -f stdlib.go
	rm -f internal_*.go

fmt:
	gofmt -w *.go t/*/*.go

vet:
	go vet *.go
