minigo: *.go
	GOOS=linux   GOARCH=amd64 go build -o minigo.linux
	GOOS=darwin  GOARCH=amd64 go build -o minigo.darwin

test: minigo.linux
	docker run --rm -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go ./linux_test.sh

clean:
	rm -f minigo*
