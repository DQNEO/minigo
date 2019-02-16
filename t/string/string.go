package main

import "fmt"

var message = "hello"

func f1() {
	fmt.Printf("%s\n", message)
}

func f2() {

	var mybytes []byte
	mybytes = []byte(message)

	fmt.Printf("%c", mybytes[0])
	fmt.Printf("%c", mybytes[1])
	fmt.Printf("%c", mybytes[2])
	fmt.Printf("%c", mybytes[3])
	fmt.Printf("%c", mybytes[4])
	fmt.Printf("\n")
}

var gfoo = "foo"

func f3() {
	foo := "foo"
	if "foo" == "foo" {
		fmt.Printf("1\n")
	}
	if foo == foo {
		fmt.Printf("2\n")
	}
	if "foo" == foo {
		fmt.Printf("3\n")
	}
	if foo == "foo" {
		fmt.Printf("4\n")
	}
	if foo == gfoo {
		fmt.Printf("5\n")
	}
}

func f4() {
	s1 := "aaa"
	if s1 != "bbb" {
		fmt.Printf("6\n")
	}

	if s1 != "" {
		fmt.Printf("7\n")
	}
}

type mystring string
func f5() {
	s := "8"
	ms := mystring(s)
	fmt.Printf("%s\n", ms)
}

func main() {
	f1()
	f2()
	f3()
	f4()
	f5()
}
