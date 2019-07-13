package main

func f1() {
	var lmap map[string]bytes = map[string]bytes{
		string("key1"):S("value1"),
	}
	key11, ok := lmap[string("key1")]
	if ok {
		fmtPrintf(S("%s\n"), key11)
	} else {
		fmtPrintf(S("ERROR\n"))
	}

	lmap[string("key2")] = S("value2")
	key2, ok := lmap[string("key2")]
	if ok {
		fmtPrintf(S("%s\n"), key2)
	} else {
		fmtPrintf(S("ERROR\n"))
	}

}

func main() {
	f1()
}
