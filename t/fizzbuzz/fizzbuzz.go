package main

import "fmt"

func main() {
	begin, end := 1, 15

	for i := begin; i <= end; i++ {
		if i%15 == 0 {
			fmt.Printf("%s\n", "FizzBuzz")
		} else if i%3 == 0 {
			fmt.Printf("%s\n", "Fizz")
		} else if i%5 == 0 {
			fmt.Printf("%s\n", "Buzz")
		} else {
			fmt.Printf("%d\n", i)
		}
	}
}
