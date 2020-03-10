package runtime

func init() {
	heapInit()
}

func printstring(b []byte) {
	write(2, b)
}

func panic(msg []byte) {
	printstring([]byte("panic: "))
	printstring(msg)
	printstring([]byte("\n"))
	exit(2)
}

const MiniGo int = 1
