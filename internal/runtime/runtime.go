package runtime

func init() {
	heapInit()
}

func printstring(b []byte) {
	var addr *byte = &b[0]
	write(2, addr, len(b))
}

func panic(msg []byte) {
	printstring([]byte("panic: "))
	printstring(msg)
	printstring([]byte("\n"))
	exit(2)
}

const MiniGo int = 1
