package main

import (
	"syscall"
	"unsafe"
)

var hdc uintptr

var hWnd uintptr

type Margins struct {
	Left, Right, Top, Bottom int32
}

func WndProc(hWnd uintptr, msg uint32, wParam, lParam uintptr) (uintptr) {
	switch msg {
	case 2: // DESTROY
		PostQuitMessage(0)
	default:
		return DefWindowProc(hWnd, msg, wParam, lParam)
	}
	return 0
}

func drawRectangle(fromX, toX, fromY, toY, size, color uintptr) {
	for y:=fromY-size;y<toY+size;y++ {
		for x:=fromX-size;x<toX+size;x++ {
			if x < fromX+size || x > toX-size || y < fromY+size || y > toY-size {
				procSetPixelV.Call(uintptr(hdc), uintptr(x), uintptr(y), color)
			}
		}
	}
}

func WinMain() {

	hInstance := GetModuleHandle("")

	lpszClassName := syscall.StringToUTF16Ptr("WNDclass")

	var wcex WNDCLASSEX
	wcex.Size        = uint32(unsafe.Sizeof(wcex))
	wcex.Style         = BLACKNESS
	wcex.WndProc   = syscall.NewCallback(WndProc)
	wcex.Instance     = hInstance
	wcex.Background = CreateBrushIndirect(uintptr(unsafe.Pointer(&LOGBRUSH{})))
	wcex.ClassName = lpszClassName

	RegisterClassEx(uintptr(unsafe.Pointer(&wcex)))

	hWnd = CreateWindowEx(
		WS_EX_TOPMOST, lpszClassName, syscall.StringToUTF16Ptr(""),
		WS_CLIPSIBLINGS | WS_CLIPCHILDREN | WS_POPUP,
		CW_USEDEFAULT, CW_USEDEFAULT, 400, 400, 0, 0, uintptr(hInstance), nil)

	SetWindowPos(hWnd, 0X00000008, 0, 0, 0 ,0, 0)
	ShowWindow(hWnd, CDS_FULLSCREEN)
	UpdateWindow(hWnd)

	SetWindowLong(hWnd, GWL_EXSTYLE, uint32(GetWindowLong(hWnd, GWL_EXSTYLE) ^ WS_EX_LAYERED ^ WS_EX_TRANSPARENT))

	procSetLayeredWindowAttributes.Call(uintptr(hWnd), 0, 255, 0x00000002)

	marg := Margins{}
	marg.Left = 0
	marg.Right = 1920
	marg.Top = 1080
	marg.Bottom = 0

	procDwmExtendFrameIntoClientArea.Call(uintptr(hWnd), uintptr(unsafe.Pointer(&marg)))

	ShowWindow(hWnd, SW_MAXIMIZE)
	UpdateWindow(hWnd)

	hdc = GetDC(hWnd)

	var msg MSG
	for {
		if GetMessage(&msg, 0, 0, 0) == 0 {
			break
		}
		TranslateMessage(&msg)
		DispatchMessage(&msg)
	}

}

func CloseHDC() {
	ReleaseDC(hWnd, hdc)
}
