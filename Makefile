# All the commands are supposed to run on Linux.
# I use Docker. Please see README.md

all: minigo /tmp/tmpfs

/tmp/tmpfs:
	mkdir -p /tmp/tmpfs

internal_runtime.go: internal/runtime/*.go cp-internalcode.sh
	./cp-internalcode.sh

internal_universe.go: internal/universe/*.go cp-internalcode.sh
	./cp-internalcode.sh

# 1st gen compiler
minigo: *.go internal_runtime.go internal_universe.go /tmp/tmpfs
	go build -o minigo *.go

# assembly for 2gen
minigo.s: minigo
	./minigo --position [a-z]*.go > /tmp/tmpfs/minigo.s
	cp /tmp/tmpfs/minigo.s minigo.s

# 2gen compiler
minigo2: minigo.s
	as -o minigo.o minigo.s
	ld -o minigo2 minigo.o

# assembly for 3gen
minigo2.s: minigo2
	./minigo2 [a-z]*.go > /tmp/tmpfs/minigo2.s
	cp /tmp/tmpfs/minigo2.s minigo2.s

# 3gen compiler
minigo3: minigo2.s
	as -o minigo2.o minigo2.s
	ld -o minigo3 minigo2.o

# assembly for 4gen
minigo3.s: minigo3
	./minigo3 [a-z]*.go > /tmp/tmpfs/minigo3.s
	cp /tmp/tmpfs/minigo3.s minigo3.s


selfhost: minigo3.s
	diff /tmp/tmpfs/minigo2.s /tmp/tmpfs/minigo3.s && echo ok

test: minigo3.s
	make vet selfhost
	./test1gen.sh
	./test2gen.sh
	./comparison-test.sh

clean:
	rm -f minigo minigo2 minigo3
	rm -f minigo*.s
	rm -f a.s a.out
	rm -rf /tmp/tmpfs/*
	rm -f internal_*.go

fmt:
	gofmt -w *.go t/*/*.go

vet:
	go vet *.go
