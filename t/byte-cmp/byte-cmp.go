package main

func f1() {
	var x byte = 'a'
	var e byte = 'e'

	if e <= 'z' {
		println("1")
	} else {
		println(x)
	}
}

func f2() {
	var c1 byte = 'p'
	var c2 byte = 'a'

	if 'a' <= c1 && c1 <= 'z' {
		println("2")
	}

	if 'a' <= c2 && c2 <= 'z' {
		println("3")
	}
}

func main() {
	f1()
	f2()
}
