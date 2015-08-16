// +build linux
package	input

import	(
	"syscall"
	"time"
	"unsafe"
)

//	linux kernel constants for klogctl()
const	(
	KLOG_CLOSE		int = iota
	KLOG_OPEN
	KLOG_READ
	KLOG_READ_ALL
	KLOG_READ_CLEAR
	KLOG_CLEAR
	KLOG_CONSOLE_OFF
	KLOG_CONSOLE_ON
	KLOG_CONSOLE_LEVEL
	KLOG_SIZE_UNREAD
	KLOG_SIZE_BUFFER
)

//	linux kernel constants for clock_gettime()
const (
	CLOCK_REALTIME 			int = iota	// Identifier for system-wide realtime clock.
	CLOCK_MONOTONIC					// Monotonic system-wide clock.
	CLOCK_PROCESS_CPUTIME_ID			// High-resolution timer from the CPU.
	CLOCK_THREAD_CPUTIME_ID				// Thread-specific CPU-time clock.
	CLOCK_MONOTONIC_RAW				// Monotonic system-wide clock, not adjusted for frequency scaling.
	CLOCK_REALTIME_COARSE				// Identifier for system-wide realtime clock, updated only on ticks.
	CLOCK_MONOTONIC_COARSE				// Monotonic system-wide clock, updated only on ticks.
	CLOCK_BOOTTIME					// Monotonic system-wide clock that includes time spent in suspension.
)



//	linux sysfs parameter for boot option printk.time
const	KRNL_PRINTK_TIME	= "/sys/module/printk/parameters/time"

//	configure and open access to printk() message buffer
func klog_open()  {
	res	:= make([]byte,8)
	syscall.Klogctl( KLOG_OPEN, nil)
	syscall.Klogctl( KLOG_CONSOLE_LEVEL, res)
	syscall.Klogctl( KLOG_CONSOLE_ON, nil)

	if(file_exists(KRNL_PRINTK_TIME)) {
		switch file_read(KRNL_PRINTK_TIME) {
			case "0", "N", "n":
				file_write(KRNL_PRINTK_TIME,"Y")
		}
	}
}


//	close access to printk() message buffer
func klog_close()  {
	syscall.Klogctl( KLOG_CONSOLE_ON, nil)
	syscall.Klogctl( KLOG_CLOSE, nil)
}


//	get the total size of printk() message buffer
func klog_buffer_size() (size int,err error)  {
	res	:= make([]byte,1)
	size,err = syscall.Klogctl( KLOG_SIZE_BUFFER, res)

	return
}


//	get the size of the unread part of printk() message buffer
func klog_unread_size() (size int,err error)  {
	res	:= make([]byte,1)
	size,err = syscall.Klogctl( KLOG_SIZE_UNREAD, res)

	return
}


//	read a part of the total printk() message buffer
func klog_fullread(size int) ([]byte,error)  {
	res	:= make([]byte, size)
	_,err	:= syscall.Klogctl( KLOG_READ_ALL, res)
	if err != nil {
		return []byte{},err
	}

	return res,nil
}


//	read a part of the unread printk() message buffer
func klog_read(size int) ([]byte,error)  {
	res	:= make([]byte, size)
	_,err:= syscall.Klogctl( KLOG_READ, res)
	if err != nil {
		return []byte{},err
	}

	return res,nil
}


//	golang version of Clock_gettime SYSCALL
func clock_gettime(clockid int, ts *syscall.Timespec) (error) {
	_, _, e1 := syscall.RawSyscall(syscall.SYS_CLOCK_GETTIME, uintptr(clockid), uintptr(unsafe.Pointer(ts)), 0)
	if e1 != 0 {
		return e1
	}
	return nil
}


func boot_time() (time.Time,error)  {
	uptime,err := uptime()
	if err == nil {
		return time.Now().Add( -uptime ),nil
	}

	return time.Now(),err
}

//	return the uptime of the machine
func uptime() (time.Duration,error)  {
	var ts syscall.Timespec
	err := clock_gettime(CLOCK_BOOTTIME, &ts)
	if err != nil {
		return time.Duration(0),err
	}

	return time.Duration(-ts.Nsec)*time.Nanosecond +time.Duration(-ts.Sec)*time.Second,nil
}
