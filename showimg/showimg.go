package showimg

import (
	"unsafe"
	"syscall"
	"os"
)

func MyRegisterClass(hInstance winapi.HINSTANCE) (atom uint16, err error) {
	var message MSG
	var hwnd HWND
	var wproc uintptr
	hwnd = CreateWindowEx(
		WS_EX_CLIENTEDGE,
		_TEXT("EDIT"),
		_TEXT("HELLO GUI"),
		WS_OVERLAPPEDWINDOW,
		(GetSystemMetrics(SM_CXSCREEN)-winWidth)>>1,
		(GetSystemMetrics(SM_CYSCREEN)-winHeight)>>1,
		winWidth,
		winHeight,
		0,
		0,
		GetModuleHandle(nil),
		unsafe.Pointer(nil))
	wproc = syscall.NewCallback(WndProc)
	originWndProc = HWND(SetWindowLong(hwnd,GWL_WNDPROC, int32(wproc)))
	ShowWindow(hwnd,SW_SHOW)
	for{
		if GetMessage(&message,0,0,0) == 0{break}
		TranslateMessage(&message)
		DispatchMessage(&message)
	}
	os.Exit(int(message.WParam))
}