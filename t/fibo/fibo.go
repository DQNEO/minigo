package main


func fibonacci(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	var r int
	r = fibonacci(12)
	fmtPrintf("%d\n", r)
}
