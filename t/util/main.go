package main

func f1() {
	baseName := getBaseNameFromImport([]byte("foo/bar"))
	fmtPrintf("%s\n", baseName) // bar
}

func main() {
	f1()
}
