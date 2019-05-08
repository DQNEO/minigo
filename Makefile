all: minigo out

out:
	mkdir out

minigo: *.go internalcode.go stdlib.go
	go build -o minigo *.go

internalcode.go: internalcode/runtime.go
	./cp-internalcode.sh

stdlib.go: stdlib/*/*.go
	./concate-stdlib.sh > stdlib.go

minigo.s: *.go minigo
	./minigo *.go > minigo.s

minigo2: minigo.s # 2nd generation
	./compat-run.sh gcc -g -no-pie -o minigo2 minigo.s

test2gen: minigo2
	./compat-run.sh ./minigo2 --version
	./compat-run.sh ./minigo2 t/min/min.go > out/min2.s
	./as out/min2.s

test1gen: all
	./test1gen.sh

test: all
	make test1gen
	make test2gen
	diff --strip-trailing-cr out/min2.s out/min.s

circlecitest: all
	make test1gen
	make test2gen
	diff --strip-trailing-cr out/min2.s out/min.s

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
