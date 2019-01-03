
all: minigo.linux minigo.darwin
	make minigo.linux minigo.darwin

minigo.linux: *.go
	GOOS=linux   GOARCH=amd64 go build -o minigo.linux

minigo.darwin: *.go
	GOOS=darwin  GOARCH=amd64 go build -o minigo.darwin

test: minigo.linux minigo.darwin
	docker run --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go ./linux_test.sh
	#make parse

parse: all
	for src in *.go; do  ./minigo.darwin --parse-only $$src && echo ok  ; done

clean:
	rm -f minigo*
	rm -f a.s a.out

fmt:
	gofmt -w *.go t/*/*.go
