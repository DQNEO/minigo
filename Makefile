all: minigo

minigo: *.go internalcode.go stdlib.go
	go build -o minigo *.go

internalcode.go: internalcode/runtime.go
	./cp-internalcode.sh

stdlib.go: stdlib/*/*.go
	./concate-stdlib.sh > stdlib.go

minigo2: *.go minigo # 2nd generation
	./minigo *.go > out/a.s
	docker run --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential bash -c 'gcc -g -no-pie -o minigo2 out/a.s'


test: all minigo2
	./compile.sh
	./test_as.sh
	./testerror.sh

circlecitest: all
	./resolve *.go
	./compile.sh
	./test_as.sh
	./testerror.sh

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
