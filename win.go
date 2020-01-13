package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32 = syscall.NewLazyDLL("user32.dll")
	dwmapi = syscall.NewLazyDLL("Dwmapi.dll")
	gdi32 = syscall.NewLazyDLL("Gdi32.dll")

	procGetModuleHandle    = kernel32.NewProc("GetModuleHandleW")
	procGetWindowLong                 = user32.NewProc("GetWindowLongW")
	procSetWindowLong                 = user32.NewProc("SetWindowLongW")
	procCreateWindowEx                = user32.NewProc("CreateWindowExW")
	procShowWindow                    = user32.NewProc("ShowWindow")
	procUpdateWindow                  = user32.NewProc("UpdateWindow")
	procGetDC                         = user32.NewProc("GetDC")
	procDefWindowProc                 = user32.NewProc("DefWindowProcW")
	procReleaseDC                     = user32.NewProc("ReleaseDC")
	procSetWindowPos                  = user32.NewProc("SetWindowPos")
	procPostQuitMessage               = user32.NewProc("PostQuitMessage")
	procSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	procDwmExtendFrameIntoClientArea = dwmapi.NewProc("DwmExtendFrameIntoClientArea")
	procSetPixelV = gdi32.NewProc("SetPixelV")
	procRegisterClassEx               = user32.NewProc("RegisterClassExW")
	procGetMessage                    = user32.NewProc("GetMessageW")
	procTranslateMessage              = user32.NewProc("TranslateMessage")
	procDispatchMessage               = user32.NewProc("DispatchMessageW")
	procCreateBrushIndirect       = gdi32.NewProc("CreateBrushIndirect")
)


func PostQuitMessage(exitCode int) {
	procPostQuitMessage.Call(
		uintptr(exitCode))
}

func DefWindowProc(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procDefWindowProc.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam)

	return ret
}

func GetDC(hwnd uintptr) uintptr {
	ret, _, _ := procGetDC.Call(
		uintptr(hwnd))

	return uintptr(ret)
}

func ShowWindow(hwnd uintptr, cmdshow int) bool {
	ret, _, _ := procShowWindow.Call(
		uintptr(hwnd),
		uintptr(cmdshow))

	return ret != 0

}

func UpdateWindow(hwnd uintptr) bool {
	ret, _, _ := procUpdateWindow.Call(
		uintptr(hwnd))
	return ret != 0
}

func ReleaseDC(hwnd uintptr, hDC uintptr) bool {
	ret, _, _ := procReleaseDC.Call(uintptr(hwnd), uintptr(hDC))
	return ret != 0
}

func SetWindowPos(hwnd, hWndInsertAfter uintptr, x, y, cx, cy int, uFlags uint) bool {
	ret, _, _ := procSetWindowPos.Call(
		uintptr(hwnd),
		uintptr(hWndInsertAfter),
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(uFlags))

	return ret != 0
}

func SetWindowLong(hwnd uintptr, index int, value uint32) uint32 {
	ret, _, _ := procSetWindowLong.Call(
		uintptr(hwnd),
		uintptr(index),
		uintptr(value))

	return uint32(ret)
}

func GetWindowLong(hwnd uintptr, index int) int32 {
	ret, _, _ := procGetWindowLong.Call(
		uintptr(hwnd),
		uintptr(index))

	return int32(ret)
}

func CreateWindowEx(exStyle uint, className, windowName *uint16,
	style uint, x, y, width, height int, parent uintptr, menu uintptr,
	instance uintptr, param unsafe.Pointer) uintptr {
	ret, _, _ := procCreateWindowEx.Call(
		uintptr(exStyle),
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		uintptr(style),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(parent),
		uintptr(menu),
		uintptr(instance),
		uintptr(param))

	return uintptr(ret)
}

func GetModuleHandle(modulename string) uintptr {
	var mn uintptr
	if modulename == "" {
		mn = 0
	} else {
		mn = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(modulename)))
	}
	ret, _, _ := procGetModuleHandle.Call(mn)
	return uintptr(ret)
}

func CreateBrushIndirect(lplb uintptr) uintptr {
	ret, _, _ := procCreateBrushIndirect.Call(
		uintptr(unsafe.Pointer(lplb)))

	return uintptr(ret)
}

func RegisterClassEx(wndClassEx uintptr) uintptr {
	ret, _, _ := procRegisterClassEx.Call(uintptr(unsafe.Pointer(wndClassEx)))
	return uintptr(ret)
}

func GetMessage(msg *MSG, hwnd uintptr, msgFilterMin, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))

	return int(ret)
}

func TranslateMessage(msg *MSG) bool {
	ret, _, _ := procTranslateMessage.Call(
		uintptr(unsafe.Pointer(msg)))

	return ret != 0

}

func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := procDispatchMessage.Call(
		uintptr(unsafe.Pointer(msg)))

	return ret

}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms633577.aspx
type WNDCLASSEX struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd145035.aspx
type LOGBRUSH struct {
	LbStyle uint32
	LbColor uint32
	LbHatch uintptr
}

// Ternary raster operations
const (
	SRCCOPY        = 0x00CC0020
	SRCPAINT       = 0x00EE0086
	SRCAND         = 0x008800C6
	SRCINVERT      = 0x00660046
	SRCERASE       = 0x00440328
	NOTSRCCOPY     = 0x00330008
	NOTSRCERASE    = 0x001100A6
	MERGECOPY      = 0x00C000CA
	MERGEPAINT     = 0x00BB0226
	PATCOPY        = 0x00F00021
	PATPAINT       = 0x00FB0A09
	PATINVERT      = 0x005A0049
	DSTINVERT      = 0x00550009
	BLACKNESS      = 0x00000042
	WHITENESS      = 0x00FF0062
	NOMIRRORBITMAP = 0x80000000
	CAPTUREBLT     = 0x40000000
)

// Extended window style constants
const (
	WS_EX_DLGMODALFRAME    = 0X00000001
	WS_EX_NOPARENTNOTIFY   = 0X00000004
	WS_EX_TOPMOST          = 0X00000008
	WS_EX_ACCEPTFILES      = 0X00000010
	WS_EX_TRANSPARENT      = 0X00000020
	WS_EX_MDICHILD         = 0X00000040
	WS_EX_TOOLWINDOW       = 0X00000080
	WS_EX_WINDOWEDGE       = 0X00000100
	WS_EX_CLIENTEDGE       = 0X00000200
	WS_EX_CONTEXTHELP      = 0X00000400
	WS_EX_RIGHT            = 0X00001000
	WS_EX_LEFT             = 0X00000000
	WS_EX_RTLREADING       = 0X00002000
	WS_EX_LTRREADING       = 0X00000000
	WS_EX_LEFTSCROLLBAR    = 0X00004000
	WS_EX_RIGHTSCROLLBAR   = 0X00000000
	WS_EX_CONTROLPARENT    = 0X00010000
	WS_EX_STATICEDGE       = 0X00020000
	WS_EX_APPWINDOW        = 0X00040000
	WS_EX_OVERLAPPEDWINDOW = 0X00000100 | 0X00000200
	WS_EX_PALETTEWINDOW    = 0X00000100 | 0X00000080 | 0X00000008
	WS_EX_LAYERED          = 0X00080000
	WS_EX_NOINHERITLAYOUT  = 0X00100000
	WS_EX_LAYOUTRTL        = 0X00400000
	WS_EX_NOACTIVATE       = 0X08000000
)

// Window style constants
const (
	WS_OVERLAPPED       = 0X00000000
	WS_POPUP            = 0X80000000
	WS_CHILD            = 0X40000000
	WS_MINIMIZE         = 0X20000000
	WS_VISIBLE          = 0X10000000
	WS_DISABLED         = 0X08000000
	WS_CLIPSIBLINGS     = 0X04000000
	WS_CLIPCHILDREN     = 0X02000000
	WS_MAXIMIZE         = 0X01000000
	WS_CAPTION          = 0X00C00000
	WS_BORDER           = 0X00800000
	WS_DLGFRAME         = 0X00400000
	WS_VSCROLL          = 0X00200000
	WS_HSCROLL          = 0X00100000
	WS_SYSMENU          = 0X00080000
	WS_THICKFRAME       = 0X00040000
	WS_GROUP            = 0X00020000
	WS_TABSTOP          = 0X00010000
	WS_MINIMIZEBOX      = 0X00020000
	WS_MAXIMIZEBOX      = 0X00010000
	WS_TILED            = 0X00000000
	WS_ICONIC           = 0X20000000
	WS_SIZEBOX          = 0X00040000
	WS_OVERLAPPEDWINDOW = 0X00000000 | 0X00C00000 | 0X00080000 | 0X00040000 | 0X00020000 | 0X00010000
	WS_POPUPWINDOW      = 0X80000000 | 0X00800000 | 0X00080000
	WS_CHILDWINDOW      = 0X40000000
)

const (
	CW_USEDEFAULT = ^0x7fffffff
)

// ChangeDisplaySettings
const (
	CDS_UPDATEREGISTRY  = 0x00000001
	CDS_TEST            = 0x00000002
	CDS_FULLSCREEN      = 0x00000004
	CDS_GLOBAL          = 0x00000008
	CDS_SET_PRIMARY     = 0x00000010
	CDS_VIDEOPARAMETERS = 0x00000020
	CDS_RESET           = 0x40000000
	CDS_NORESET         = 0x10000000

	DISP_CHANGE_SUCCESSFUL  = 0
	DISP_CHANGE_RESTART     = 1
	DISP_CHANGE_FAILED      = -1
	DISP_CHANGE_BADMODE     = -2
	DISP_CHANGE_NOTUPDATED  = -3
	DISP_CHANGE_BADFLAGS    = -4
	DISP_CHANGE_BADPARAM    = -5
	DISP_CHANGE_BADDUALVIEW = -6
)

// GetWindowLong and GetWindowLongPtr constants
const (
	GWL_EXSTYLE     = -20
	GWL_STYLE       = -16
	GWL_WNDPROC     = -4
	GWLP_WNDPROC    = -4
	GWL_HINSTANCE   = -6
	GWLP_HINSTANCE  = -6
	GWL_HWNDPARENT  = -8
	GWLP_HWNDPARENT = -8
	GWL_ID          = -12
	GWLP_ID         = -12
	GWL_USERDATA    = -21
	GWLP_USERDATA   = -21
)

// ShowWindow constants
const (
	SW_HIDE            = 0
	SW_NORMAL          = 1
	SW_SHOWNORMAL      = 1
	SW_SHOWMINIMIZED   = 2
	SW_MAXIMIZE        = 3
	SW_SHOWMAXIMIZED   = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}