package main


func main() {
	var heapHead *int
	var heapTail *int

	var address *int

	address = malloc(8)
	*address = 1
	fmtPrintf("%d\n", *address)
	address = malloc(8)
	*address = 2
	fmtPrintf("%d\n", *address)
	address = malloc(8)
	*address = 3
	fmtPrintf("%d\n", *address)

	heapA := malloc(8)
	heapB := malloc(0)

	fmtPrintf("%d\n", (heapB-heapA)-4) // 4
}
