// Translation of http://man7.org/linux/man-pages/man2/getdents64.2.html#top_of_page
package main

import (
	"os"
	"unsafe"
	"syscall"
)
import "fmt"

var _buf [1024]byte
const _x_sys_getdents64  = 217

/*
struct linux_dirent64 {
	ino64_t        d_ino;    // 8 bytes: 64-bit inode number
	off64_t        d_off;    // 8 bytes: 64-bit offset to next structure
	unsigned short d_reclen; // 2 bytes: Size of this dirent
	unsigned char  d_type;   // 1 byte: File type
	char           d_name[]; // Filename (null-terminated)
};
 */
type linux_dirent struct {
	d_ino    int
	d_off    int
	d_reclen1 uint16
	d_type   byte
	d_name   byte
}

func cstring2string(b *byte) string {
	var bs []byte
	for {
		if b == nil || *b == 0 {
			break
		}
		bs = append(bs, *b)
		b = uintptr(b) + 1
	}
	return string(bs)
}

func print_dirp(dirp *linux_dirent) {
	var reclen int = int(dirp.d_reclen1)

	//fmt.Printf("%p  ", uintptr(dirp))
	//fmt.Printf("%d\t", dirp.d_ino)
	//fmt.Printf("%d\t", dirp.d_off)
	//fmt.Printf("%d\t", dirp.d_type)
	//fmt.Printf("%d\t", reclen)
	//reclen := int(dirp.d_reclen1)
	//fmt.Printf("%d  ", dirp.d_type)
	var bp *byte = uintptr(unsafe.Pointer(&dirp.d_name))
	var s string = cstring2string(bp)
	return
	fmt.Printf("%s", s)
	fmt.Printf("\n")
}

func main() {
	f , err := os.Open("t/data")
	if err != nil {
		panic(err)
	}
	fd := f.Fd()
	var buf []byte = _buf[:]
	var counter int
	for {
		nread,_ := syscall.Getdents(int(fd), buf)
		if nread == -1 {
			panic("getdents failed")
		}
		if nread == 0 {
			break
		}

		//fmt.Printf("--------------- nread=%d ---------------\n", nread)
		//fmt.Printf("inode   d_off   d_type  d_reclen    d_name\n")
		for bpos := 0; bpos < nread; 1 {
			var dirp *linux_dirent
			p := uintptr(unsafe.Pointer(&buf[0])) + uintptr(bpos)
			dirp = p
			print_dirp(dirp)
			bpos = bpos + int(dirp.d_reclen1) // 24 is wrong
			counter++
		}
	}

	fmt.Printf("%d\n", counter)
}
