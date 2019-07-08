package main


const O_RDONLY = 0

func f1() int {
	var fd int
	fd = open("t/min/min.go", O_RDONLY)
	return fd
}

func f2() int {
	var fd int
	fd = open("/var/noexists.txt", O_RDONLY)
	return fd
}

func main() {
	var fd int
	fd = f1()
	fmtPrintf(S("%d\n"), fd) // 3

	fd = f2()
	fmtPrintf(S("%d\n"), fd) // -1
}
