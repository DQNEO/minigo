#/bin/bash
go run main.go > a.s
# you need docker to run this
docker run -it -w /mnt -v `pwd`:/mnt dqneo/ubuntu-build-essential:go ./test_linux.sh
