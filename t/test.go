package main

func fa() {
}

func fb() {
	printf("%d\n", 0)
	printf("%d\n", 2 - 1)
	printf("%d\n", 1 + 1)
	printf("%d\n", 1 + 1 + 1)
	printf("%d\n", 2 * 2)
	printf("%d\n", 2 * 3 - 1)
	printf("%d\n", 1 + 1 * 5)
}

func fc() {
	var i int
	i = 3
	printf("%d\n", i + 4)
}

func main() {
	fa()
	fb()
	fc()
	puts("hello world")
}

