package main

import "fmt"

func main() {
	for i := 1; i <= 15; i = i + 1 {
		if i == 3 || i == 6 || i == 9 || i == 12 {
			fmt.Printf("Fizz\n")
		} else if i == 5 || i == 10 {
			fmt.Printf("Buzz\n")
		} else if i == 15 {
			fmt.Printf("FizzBuzz\n")
		} else {
			fmt.Printf("%d\n", i)
		}
	}
}
