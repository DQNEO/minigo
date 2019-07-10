package main

func f1() {
	panic([]byte("Help me!"))
}

func main() {
	f1()
}
