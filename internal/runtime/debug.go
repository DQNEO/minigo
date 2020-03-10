package runtime

func reportMemoryUsage() {
	printstring([]byte("# memory-usage: "))
	i := getMemoryUsage()
	s := itoa(int(i))
	printstring([]byte(s))
	printstring([]byte("\n"))
}

func getMemoryUsage() uintptr {
	return heapPtr - heapHead
}

// This is a copy from stconv
func itoa(i int) []byte {
	var r []byte
	var tmp []byte
	var isMinus bool

	// open(2) returs  0xffffffff 4294967295 on error.
	// I don't understand this yet.
	if i > 2147483648 {
		i = i - 2147483648*2
	}

	if i < 0 {
		i = i * -1
		isMinus = true
	}
	for i > 0 {
		mod := i % 10
		tmp = append(tmp, byte('0')+byte(mod))
		i = i / 10
	}

	if isMinus {
		r = append(r, '-')
	}

	for j := len(tmp) - 1; j >= 0; j-- {
		r = append(r, tmp[j])
	}

	if len(r) == 0 {
		return []byte{'0'}
	}
	return r
}
