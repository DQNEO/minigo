GO_OBJS := $(*.go)

minigo: $(GO_OBJS)
	go build -o minigo $(GO_OBJS)

clean:
	rm minigo
