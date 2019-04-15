package os

var Stdout *File

var SStdout = File{
	id:1,
}

var Stderr *File

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

var Args []string

func Exit(i int) {
}

func Init() {
	Stdout = &SStdout
	Stderr = &SStderr
}
