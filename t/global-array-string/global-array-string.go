package main


var messages = [2]gostring{gostring("hello"), gostring("world")}

func main() {
	fmtPrintln(messages[1])
}
