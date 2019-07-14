package main

func f1() {
	var lmap map[string][]byte = map[string][]byte{
		string("key1"):S("value1"),
	}
	key11, ok := lmap[string("key1")]
	if ok {
		fmtPrintf("%s\n", key11)
	} else {
		fmtPrintf("ERROR\n")
	}

	lmap[string("key2")] = S("value2")
	key2, ok := lmap[string("key2")]
	if ok {
		fmtPrintf("%s\n", key2)
	} else {
		fmtPrintf("ERROR\n")
	}

}

func main() {
	f1()
}
