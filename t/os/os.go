package main

var Stdout *File
var Stderr *File

var SStdout = File{
	id:1,
}
var SStderr = File{
	id:2,
}

// File represents an open file descriptor.
type File struct {
	id int
}

func (f *File) Write(b string) (int, error) {
	var fid int = f.id
	var n int
	n = write(fid, b, len(b))
	return n,nil
}

func f1() {
	var b = "hello stdout\n"
	Stdout = &SStdout
	Stdout.Write(b)

	var b = "hello stderr\n"
	Stderr = &SStderr
	Stderr.Write(b)

}

func main() {
	f1()
}
