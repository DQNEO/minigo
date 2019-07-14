package main


// For statements with single condition
func main() {
	var i int = 0
	for i < 5 {
		fmtPrintf("%d\n", i)
		i = i + 1
	}
	var j int = i
	for 10 > j {
		fmtPrintf("%d\n", j)
		j = j + 1
	}
}
