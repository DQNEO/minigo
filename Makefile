
all: minigo.linux minigo.darwin
	make minigo.linux minigo.darwin

minigo.linux: *.go
	GOOS=linux   GOARCH=amd64 go build -o minigo.linux

minigo.darwin: *.go
	GOOS=darwin  GOARCH=amd64 go build -o minigo.darwin

test: minigo.linux minigo.darwin
	docker run --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go ./linux_test.sh

clean:
	rm -f minigo*

fmt:
	gofmt -w *.go t/*/*.go
