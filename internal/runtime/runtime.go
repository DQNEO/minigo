package runtime

var _argv []*byte

func init() {
	heapInit()
}

func runtime_args() []string {
	var r []string
	for _, a := range _argv {
		// Convert *byte to string
		var s string = cstring2string(a)
		r = append(r, s)
	}
	return r
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
