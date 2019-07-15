package main

import "fmt"

var message string = "hello"

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

var gfoo string = string("foo")

func f3() {
	foo := "foo"
	if "foo" == "foo" {
		fmt.Printf("1\n")
	}
	if foo  == foo {
		fmt.Printf("2\n")
	}
	if "foo"  ==  foo {
		fmt.Printf("3\n")
	}
	if foo ==  "foo" {
		fmt.Printf("4\n")
	}
	if foo ==  gfoo {
		fmt.Printf("5\n")
	}
}

func f4() {
	s1 := "aaa"
	if s1 !=  "bbb" {
		fmt.Printf("6\n")
	}

	if s1 !=  "" {
		fmt.Printf("7\n")
	}
}

type mystring []byte

func f5() {
	s := "8"
	ms := mystring(s)
	fmt.Printf("%s\n", ms) // 8
}

func f6() {
	fmt.Printf("%d\n", len("123456789")) // 9
	s := "0123456789"
	fmt.Printf("%d\n", len(s)) // 10
}

func f7() {
	s := `11
12`
	fmt.Printf("%s\n", s) // 11 12
}

func f8() {
	s := "dummy"
	if s ==  "" {
		fmt.Printf("ERROR\n")
	}
}

func f9() {
	var s2 string
	fmt.Printf("%s", s2)
	if s2 ==  "" {
		fmt.Printf("13\n")
	}
}

func f10() {
	var str0 string
	fmt.Printf("1%s4\n", str0)

	var str1 string = ""
	fmt.Printf("1%s5\n", str1)

	str2 := ""
	fmt.Printf("1%s6\n", str2)

	fmt.Printf("%d\n", len(str0) + len(str1) + 17) // 17

	str3 := "abc\n"

	fmt.Printf("%d\n", len(str3) + 14) // 18
}

func main() {
	f1()
	f2()
	f3()
	f4()
	f5()
	f6()
	f7()
	f8()
	f9()
	f10()
}
