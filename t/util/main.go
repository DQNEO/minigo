package main

func f1() {
	baseName := getBaseNameFromImport(bytes("foo/bar"))
	fmtPrintf("%s\n", baseName) // bar
}

func main() {
	f1()
}
