# All the commands are supposed to run on Linux.
# I use Docker like below.
# docker run -it --rm -w /mnt -v `pwd`:/mnt --cap-add=SYS_PTRACE --security-opt='seccomp=unconfined' --tmpfs=/tmp/tmpfs:rw,size=500m,mode=1777 dqneo/ubuntu-build-essential:go bash

all: minigo /tmp/tmpfs

/tmp/tmpfs:
	mkdir -p /tmp/tmpfs

internal_runtime.go: internal/runtime/*.go
	./cp-internalcode.sh

internal_universe.go: internal/universe/*.go
	./cp-internalcode.sh

stdlib.go: stdlib/*/*.go
	./concate-stdlib.sh > stdlib.go

# 1st gen compiler
minigo: *.go internal_runtime.go internal_universe.go stdlib.go /tmp/tmpfs
	go build -o minigo *.go

# 2nd gen assembly
minigo.s: minigo
	./minigo --position [a-z]*.go > /tmp/tmpfs/minigo.s
	cp /tmp/tmpfs/minigo.s minigo.s

# 2nd gen compiler
minigo2: minigo.s
	gcc -g -no-pie -o minigo2 minigo.s

minigo2.s: minigo2 minigo *.go
	./minigo2 [a-z]*.go > /tmp/tmpfs/minigo2.s
	cp /tmp/tmpfs/minigo2.s minigo2.s

selfhost: minigo2.s
	sed -e 's|^/\*.*)\*/||' minigo.s > /tmp/tmpfs/minigo.stripped.s
	diff minigo2.s /tmp/tmpfs/minigo.stripped.s && echo ok

test: minigo minigo2
	make vet
	./test1gen.sh
	./test2gen.sh
	./comparison-test.sh

clean:
	rm -f minigo minigo2
	rm -f minigo*.s
	rm -f a.s a.out
	rm -rf /tmp/tmpfs/*
	rm -f stdlib.go
	rm -f internal_*.go

fmt:
	gofmt -w *.go t/*/*.go

vet:
	go vet *.go
