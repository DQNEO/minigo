package main


func main() {
	begin, end := 1, 15

	for i := begin; i <= end; i++ {
		if i%15 == 0 {
			fmtPrintf("%s\n", "FizzBuzz")
		} else if i%3 == 0 {
			fmtPrintf("%s\n", "Fizz")
		} else if i%5 == 0 {
			fmtPrintf("%s\n", "Buzz")
		} else {
			fmtPrintf("%d\n", i)
		}
	}
}
