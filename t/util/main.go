package main

func f1() {
	baseName := getBaseNameFromImport(gostring("foo/bar"))
	fmtPrintf(S("%s\n"), baseName) // bar
}

func main() {
	f1()
}
