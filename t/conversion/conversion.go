package main

import "fmt"

func f1() {
	var vbytes []byte
	var s []byte
	s = []byte(vbytes)
	fmt.Printf("%s0\n", s)           // 0
	fmt.Printf("%d\n", len(vbytes)+1) // 1
	fmt.Printf("%d\n", len(s)+2)     // 2
}

func f2() {
	var s []byte
	fmt.Printf("%s3\n", []byte(s))       // 3
	fmt.Printf("%d\n", len(s)+4) // 4
}

func f3() {
	var s string = ""
	fmt.Printf("%s5\n", s)       // 5
	fmt.Printf("%d\n", len(s)+6) // 6
}

func f4() {
	var s []byte
	var vbytes []byte
	vbytes = []byte(s)
	fmt.Printf("%s7\n", []byte(vbytes)) // 7
	fmt.Printf("%d\n", len(vbytes)+8)   // 8
}

func f5() {
	var s []byte
	var bs []byte
	bs = []byte(s)
	if bs == nil {
		fmt.Printf("9\n")
	} else {
		fmt.Printf("ERROR")
	}
}

func main() {
	f1()
	f2()
	f3()
	f4()
	f5()
}
