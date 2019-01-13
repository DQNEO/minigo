
all: minigo.linux minigo.darwin

minigo.linux: *.go stdlib.go
	GOOS=linux   GOARCH=amd64 go build -o minigo.linux

minigo.darwin: *.go stdlib.go
	GOOS=darwin  GOARCH=amd64 go build -o minigo.darwin

stdlib.go: stdlib/*/*.go
	./concate-stdlib.sh > stdlib.go


test: all
	./resolve *.go
	docker run --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go ./linux_test.sh

parse: all
	./parse *.go

clean:
	rm -f minigo*
	rm -f a.s a.out

fmt:
	gofmt -w *.go t/*/*.go
