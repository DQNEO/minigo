package main


var messages = [2]bytes{bytes("hello"), bytes("world")}

func main() {
	fmtPrintln(string(messages[1]))
}
