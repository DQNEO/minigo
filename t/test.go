package main

func fa() {
	var i int
	i = 0
	printf("%d\n", i)
}

func fb() {
	printf("%d\n", 2 - 1)
	printf("%d\n", 1 + 1)
	printf("%d\n", 1 + 1 + 1)
	printf("%d\n", 2 * 2)
	printf("%d\n", 2 * 3 - 1)
	printf("%d\n", 1 + 1 * 5)
}

func main() {
	fa()
	fb()
	puts("hello world")
}

