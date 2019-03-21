package main

func f1() {
	var x byte
	var e byte = 'e'

	if e <= 'z' {
		println("1")
	} else {
		println("NG")
	}
}

func main() {
	f1()
}
