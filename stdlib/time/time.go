package time

import "syscall"

func Sleep(msec int) {
	if msec > 1000 {
		panic("unsupported second")
	}
	ts := syscall.Timespec{
		Sec:  0,
		Nsec: msec * 1000000,
	}

	syscall.Nanosleep(&ts, nil)
}
