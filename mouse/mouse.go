package mouse

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type DWORD int32

var setCursorPos *windows.LazyProc
var getCursorPos *windows.LazyProc
var mouseEvent *windows.LazyProc

func init() {
	libuser32 := windows.NewLazySystemDLL("user32.dll")
	setCursorPos = libuser32.NewProc("SetCursorPos")
	getCursorPos = libuser32.NewProc("GetCursorPos")
	mouseEvent = libuser32.NewProc("mouse_event")
}

type POINT struct {
	X, Y int32
}

func GetCursorPos(lpPoint *POINT) bool {
	ret, _, _ := syscall.Syscall(getCursorPos.Addr(), 1,
		uintptr(unsafe.Pointer(lpPoint)),
		0,
		0)

	return ret != 0
}

func SetCursorPos(X, Y int32) bool {
	ret, _, _ := syscall.Syscall(setCursorPos.Addr(), 2,
		uintptr(X),
		uintptr(Y),
		0)

	return ret != 0
}

const (
	MOUSEEVENTF_ABSOLUTE        = 0x8000
	MOUSEEVENTF_HWHEEL          = 0x1000
	MOUSEEVENTF_MOVE            = 0x0001
	MOUSEEVENTF_MOVE_NOCOALESCE = 0x2000
	MOUSEEVENTF_LEFTDOWN        = 0x0002
	MOUSEEVENTF_LEFTUP          = 0x0004
	MOUSEEVENTF_RIGHTDOWN       = 0x0008
	MOUSEEVENTF_RIGHTUP         = 0x0010
	MOUSEEVENTF_MIDDLEDOWN      = 0x0020
	MOUSEEVENTF_MIDDLEUP        = 0x0040
	MOUSEEVENTF_VIRTUALDESK     = 0x4000
	MOUSEEVENTF_WHEEL           = 0x0800
	MOUSEEVENTF_XDOWN           = 0x0080
	MOUSEEVENTF_XUP             = 0x0100
)

func MouseEvent(MOUSEINPUT uint16, p *POINT) bool {
	ret, _, _ := syscall.Syscall(
		mouseEvent.Addr(),
		3,
		uintptr(MOUSEINPUT),
		uintptr(p.X),
		uintptr(p.Y),
	)

	return ret != 0
}

func Scroll(delta DWORD) bool {
	ret, _, _ := syscall.Syscall6(
		mouseEvent.Addr(),
		5,
		uintptr(MOUSEEVENTF_WHEEL),
		0,
		0,
		uintptr(delta),
		0,
		0,
	)

	return ret != 0
}

// void sendMouseRightclick(Point p)
// {
//     mouse_event(MOUSEEVENTF_RIGHTDOWN | MOUSEEVENTF_RIGHTUP, p.X, p.Y, 0, 0);
// }
