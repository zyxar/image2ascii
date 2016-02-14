//+build windows

package term

// credit: https://github.com/nsf/termbox-go/blob/master/termbox_windows.go
import (
	"syscall"
	"unsafe"
)

type coord struct {
	x int16
	y int16
}

type console_screen_buffer_info struct {
	size coord
	_    coord
	_    uint16
	_    struct {
		_ int16
		_ int16
		_ int16
		_ int16
	}
	_ coord
}

var (
	consoleInfo                         console_screen_buffer_info
	kernel32                            = syscall.NewLazyDLL("kernel32.dll")
	proc_get_console_screen_buffer_info = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

func get_console_screen_buffer_info(h syscall.Handle, info *console_screen_buffer_info) (err error) {
	r0, _, e1 := syscall.Syscall(proc_get_console_screen_buffer_info.Addr(),
		2, uintptr(h), uintptr(unsafe.Pointer(info)), 0)
	if int(r0) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Width(fd uintptr) int {
	if err := get_console_screen_buffer_info(syscall.Handle(fd), &consoleInfo); err != nil {
		panic(err)
	}
	return int(consoleInfo.size.x)
}
