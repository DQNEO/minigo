package main

type gostring []byte

func f1() {
	baseName := getBaseNameFromImport(gostring("foo/bar"))
	fmtPrintf(S("%s\n"), string(baseName)) // bar
}

func main() {
	f1()
}
