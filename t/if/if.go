package main

func main() {
	var t bool = true
	if t {
		println("1")
	}
	t = false
	if t {
		println("Error")
	}
	println("2")

	t = true
	if t {
		println("3")
	} else {
		println("Error")
	}

	t = false
	if t {
		println("Error")
	} else {
		println("4")
	}

	var i int
	i = 1
	if i == 1 {
		println("5")
	} else if i == 2 {
		println("Error")
	} else {
		println("Error")
	}

	i = 2
	if i == 1 {
		println("Error")
	} else if i == 2 {
		println("6")
	} else {
		println("Error")
	}

	if i = 3; i == 1 {
		println("Error")
	} else if i == 2 {
		println("Error")
	} else {
		println("7")
	}
}
