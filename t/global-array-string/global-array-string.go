package main


var messages = [2]bytes{bytes("hello"), bytes("world")}

func main() {
	fmtPrintln(messages[1])
}
