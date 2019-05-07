all: minigo out

out:
	mkdir out

minigo: *.go internalcode.go stdlib.go
	go build -o minigo *.go

internalcode.go: internalcode/runtime.go
	./cp-internalcode.sh

stdlib.go: stdlib/*/*.go
	./concate-stdlib.sh > stdlib.go

minigo2: *.go minigo # 2nd generation
	./minigo *.go > out/a.s
	./compat-run.sh gcc -g -no-pie -o minigo2 out/a.s

test: all
	./compile.sh
	./test_as.sh
	./testerror.sh
	make minigo2
	./compat-run.sh ./minigo2 --version
	./compat-run.sh ./minigo2 t/min/min.go > out/a.s
	./as

circlecitest: all
	make minigo2
	./compat-run.sh ./minigo2 --version
	./compat-run.sh ./minigo2 t/min/min.go > out/a.s
	./as
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
