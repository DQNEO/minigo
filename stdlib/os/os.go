package os

var SStdout = File{
	id:1,
}
var SStderr = File{
	id:2,
}

var Stderr *File
var Stdout *File


// File represents an open file descriptor.
type File struct {
	id int
}

func (f *File) Write(b []byte) (int, error) {
	var fid int = f.id
	var n int
	n = write(fid, string(b), len(b))
	return n,nil
}

func Exit(i int) {
}

func init() {
	Stdout = &SStdout
	Stderr = &SStderr

	// runtime_args is written in assembly code
	Args = runtime_args()
}

//func runtime_args() []string
