package main

import "fmt"

var hello = [5]byte{'h', 'e', 'l', 'l', 'o'}

func ghello() {
	fmt.Printf("%c", hello[0])
	fmt.Printf("%c", hello[1])
	fmt.Printf("%c", hello[2])
	fmt.Printf("%c", hello[3])
	fmt.Printf("%c", hello[4])
	fmt.Printf("%s", "\n")

	fmt.Printf("%s\n", hello)

}

func lworld() {
	var world = [5]byte{'w', 'o', 'r', 'l', 'd'}
	fmt.Printf("%c", world[0])
	fmt.Printf("%c", world[1])
	fmt.Printf("%c", world[2])
	fmt.Printf("%c", world[3])
	fmt.Printf("%c", world[4])
	fmt.Printf("%s", "\n")

	fmt.Printf("%s\n", world)
}

func fappend() {
	var chars []byte
	chars = append(chars, '7')
	chars = append(chars, '8')
	fmt.Printf("%d\n", len(chars)+4) // 6
	fmt.Printf("%c\n", chars[0])     // 7
	fmt.Printf("%c\n", chars[1])     // 8
	fmt.Printf("9\n")                // 9

	chars[0] = '1'
	chars[1] = '0'
	fmt.Printf("%s\n", chars) // 10
}

func main() {
	ghello()
	lworld()
	fappend()
}
