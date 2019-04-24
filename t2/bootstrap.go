package runtime

var heapTail *int

func f1(size int) *int {
	if heapTail + size > 0 {
		return 0
	}
}

