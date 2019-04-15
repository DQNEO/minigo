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

func (f *File) Write(b string) (int, error) {
	var fid int = f.id
	var n int
	n = write(fid, b, len(b))
	return n,nil
}

var Args []string

func Exit(i int) {
}

func init() {
	Stdout = &SStdout
	Stderr = &SStderr

	//runtime_args()
	//Args = runtime_args()
}
