all: minigo

minigo: *.go stdlib.go
	go build -o minigo

stdlib.go: stdlib/*/*.go
	./concate-stdlib.sh > stdlib.go

test: all
	./resolve *.go
	./compile.sh
	docker run --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go ./linux_test.sh

parse: all
	./parse *.go

clean:
	rm -f minigo*
	rm -f a.s a.out
	rm -f out/*
	rm -f stdlib.go

fmt:
	gofmt -w *.go t/*/*.go
