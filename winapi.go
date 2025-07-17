package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	findWindow             = user32.NewProc("FindWindowW")
	findWindowEx           = user32.NewProc("FindWindowExW")
	sendMessageTimeout     = user32.NewProc("SendMessageTimeoutW")
	enumWindows            = user32.NewProc("EnumWindows")
	getWindowText          = user32.NewProc("GetWindowTextW")
	setParent              = user32.NewProc("SetParent")
	showWindow             = user32.NewProc("ShowWindow")
	getSystemMetrics      = user32.NewProc("GetSystemMetrics")
	
	ole32                  = syscall.NewLazyDLL("ole32.dll")
	coInitialize          = ole32.NewProc("CoInitialize")
	coCreateInstance      = ole32.NewProc("CoCreateInstance")
)

const (
	CLSID_TaskbarList = "{56FDF344-FD6D-11d0-958A-006097C9A090}"
	IID_ITaskbarList  = "{56FDF31B-FD6D-11d0-958A-006097C9A090}"
	
	SW_HIDE     = 0
	SMTO_NORMAL = 0
	
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

var (
	workerw uintptr
	target  uintptr
	searchTitle string // 用于在 EnumWindowsProc2 中传递搜索字符串
)

// GetScreenWidth 获取屏幕宽度
func GetScreenWidth() int {
	ret, _, _ := getSystemMetrics.Call(uintptr(SM_CXSCREEN))
	return int(ret)
}

// GetScreenHeight 获取屏幕高度
func GetScreenHeight() int {
	ret, _, _ := getSystemMetrics.Call(uintptr(SM_CYSCREEN))
	return int(ret)
}

// CLSIDFromString 将 CLSID 字符串转换为 GUID 结构
func CLSIDFromString(str string) (*syscall.GUID, error) {
	var guid syscall.GUID
	str16 := syscall.StringToUTF16Ptr(str)
	ret, _, _ := syscall.NewLazyDLL("ole32.dll").NewProc("CLSIDFromString").Call(
		uintptr(unsafe.Pointer(str16)),
		uintptr(unsafe.Pointer(&guid)))
	if ret != 0 {
		return nil, fmt.Errorf("CLSIDFromString failed with error code %d", ret)
	}
	return &guid, nil
}

// EnumWindowsProc1 是 EnumWindows 的回调函数
func EnumWindowsProc1(hwnd uintptr, lparam uintptr) uintptr {
	defview, _, _ := findWindowEx.Call(
		hwnd,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("SHELLDLL_DefView"))),
		0,
	)
	if defview != 0 {
		workerw, _, _ = findWindowEx.Call(
			0,
			hwnd,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("WorkerW"))),
			0,
		)
	}
	return 1 // 继续枚举
}

// SetDesktop 将窗口设置为桌面壁纸
func SetDesktop(hwnd uintptr) {
	// Find Progman window
	progman, _, _ := findWindow.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Progman"))),
		0,
	)

	// Send 0x052C message to Progman
	sendMessageTimeout.Call(
		progman,
		0x052C,
		0,
		0,
		SMTO_NORMAL,
		0x3E8,
		0,
	)

	// Enumerate windows to find WorkerW
	syscall.NewCallback(EnumWindowsProc1)
	enumWindows.Call(
		syscall.NewCallback(EnumWindowsProc1),
		0,
	)

	// Hide WorkerW
	showWindow.Call(workerw, SW_HIDE)

	// Set parent
	setParent.Call(hwnd, progman)
}

// EnumWindowsProc2 是查找窗口的回调函数
func EnumWindowsProc2(hwnd, lparam uintptr) uintptr {
	var title [256]uint16
	getWindowText.Call(
		hwnd,
		uintptr(unsafe.Pointer(&title[0])),
		256,
	)
	
	// 转换为 Go 字符串
	titleStr := syscall.UTF16ToString(title[:])
	
	if titleStr == searchTitle {
		target = hwnd
		return 0 // 停止枚举
	}
	return 1 // 继续枚举
}

// RemoveFromTaskbar 从任务栏移除窗口
func RemoveFromTaskbar(hwnd uintptr) error {
	// 初始化 COM
	coInitialize.Call(0)

	// 获取 TaskbarList CLSID
	clsid, err := CLSIDFromString(CLSID_TaskbarList)
	if err != nil {
		return fmt.Errorf("获取 CLSID 失败: %v", err)
	}

	// 获取 ITaskbarList IID
	iid, err := CLSIDFromString(IID_ITaskbarList)
	if err != nil {
		return fmt.Errorf("获取 IID 失败: %v", err)
	}

	var pTaskbar **struct{ vtbl *struct{ 
		QueryInterface uintptr
		AddRef        uintptr
		Release       uintptr
		HrInit        uintptr
		AddTab        uintptr
		DeleteTab     uintptr
		ActivateTab   uintptr
		SetActiveAlt  uintptr
	}}

	ret, _, _ := coCreateInstance.Call(
		uintptr(unsafe.Pointer(clsid)),
		0,
		21, // CLSCTX_INPROC_SERVER
		uintptr(unsafe.Pointer(iid)),
		uintptr(unsafe.Pointer(&pTaskbar)),
	)

	if ret != 0 {
		return fmt.Errorf("CoCreateInstance 失败，错误码: %d", ret)
	}

	// 调用 HrInit
	syscall.Syscall((*pTaskbar).vtbl.HrInit, 1, uintptr(unsafe.Pointer(*pTaskbar)), 0, 0)
	
	// 调用 DeleteTab
	syscall.Syscall((*pTaskbar).vtbl.DeleteTab, 2, uintptr(unsafe.Pointer(*pTaskbar)), hwnd, 0)

	return nil
}

// FindWindowByTitle 通过标题查找窗口
func FindWindowByTitle(title string) uintptr {
	target = 0
	searchTitle = title
	enumWindows.Call(
		syscall.NewCallback(EnumWindowsProc2),
		0,
	)
	return target
}

// SetupWallpaper 查找窗口并设置为壁纸
func SetupWallpaper(windowTitle string) bool {
	hwnd := FindWindowByTitle(windowTitle)
	if hwnd != 0 {
		RemoveFromTaskbar(hwnd)
		SetDesktop(hwnd)
		return true
	}
	return false
}

// GetWindowsDarkMode 检测 Windows 是否为暗色模式
func GetWindowsDarkMode() bool {
	regKey, err := syscall.UTF16PtrFromString(`Software\\Microsoft\\Windows\\CurrentVersion\\Themes\\Personalize`)
	if err != nil {
		return false
	}
	var hKey syscall.Handle
	err = syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_USER, regKey, 0, syscall.KEY_READ, &hKey)
	if err != nil {
		return false
	}
	defer syscall.RegCloseKey(hKey)

	var typ uint32
	var data [4]byte
	var dataLen uint32 = 4
	valueName, _ := syscall.UTF16PtrFromString("AppsUseLightTheme")
	err = syscall.RegQueryValueEx(hKey, valueName, nil, &typ, (*byte)(unsafe.Pointer(&data[0])), &dataLen)
	if err != nil || typ != syscall.REG_DWORD {
		return false
	}
	// 0 表示暗色模式，1 表示亮色模式
	return data[0] == 0
}
