all: minigo

minigo: *.go internalcode.go stdlib.go
	go build -o minigo

internalcode.go: internalcode/runtime.go
	./cp-internalcode.sh

stdlib.go: stdlib/*/*.go
	./concate-stdlib.sh > stdlib.go

child: minigo # selfhost
	./minigo *.go > out/a.s
	docker run --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential bash -c 'gcc -g -no-pie -o child out/a.s'


test: all
	./resolve *.go
	./compile.sh
	./test_as.sh
	./test2.sh

travistest: all
	./resolve *.go
	./compile.sh
	./test_as.sh
	./test2.sh

parse: all
	./parse *.go

clean:
	rm -f minigo*
	rm -f a.s a.out
	rm -f out/*
	rm -f stdlib.go
	rm -f internalcode.go

fmt:
	gofmt -w *.go t/*/*.go
