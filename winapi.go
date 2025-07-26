package main

import (
	"context"
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                    = syscall.NewLazyDLL("user32.dll")
	findWindow                = user32.NewProc("FindWindowW")
	findWindowEx              = user32.NewProc("FindWindowExW")
	sendMessageTimeout        = user32.NewProc("SendMessageTimeoutW")
	enumWindows               = user32.NewProc("EnumWindows")
	getWindowText             = user32.NewProc("GetWindowTextW")
	getClassName              = user32.NewProc("GetClassNameW")
	getWindowThreadProcessId  = user32.NewProc("GetWindowThreadProcessId")
	setParent                 = user32.NewProc("SetParent")
	showWindow                = user32.NewProc("ShowWindow")
	getSystemMetrics          = user32.NewProc("GetSystemMetrics")
	setWindowLongPtr          = user32.NewProc("SetWindowLongPtrW")
	getWindowLongPtr          = user32.NewProc("GetWindowLongPtrW")
	setWindowPos              = user32.NewProc("SetWindowPos")
	getWindow                 = user32.NewProc("GetWindow")
	moveWindow                = user32.NewProc("MoveWindow")
	setWindowPlacement        = user32.NewProc("SetWindowPlacement")
	getWindowRect             = user32.NewProc("GetWindowRect")
	adjustWindowRect          = user32.NewProc("AdjustWindowRect")
	setLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")

	kernel32                  = syscall.NewLazyDLL("kernel32.dll")
	openProcess               = kernel32.NewProc("OpenProcess")
	getModuleBaseName         = kernel32.NewProc("GetModuleBaseNameW")
	closeHandle               = kernel32.NewProc("CloseHandle")

	psapi                     = syscall.NewLazyDLL("psapi.dll")
	getModuleBaseNameW        = psapi.NewProc("GetModuleBaseNameW")

	dwmapi                    = syscall.NewLazyDLL("dwmapi.dll")
	dwmExtendFrameIntoClientArea = dwmapi.NewProc("DwmExtendFrameIntoClientArea")
	dwmEnableBlurBehindWindow    = dwmapi.NewProc("DwmEnableBlurBehindWindow")

	gdi32                     = syscall.NewLazyDLL("gdi32.dll")
	createRectRgn             = gdi32.NewProc("CreateRectRgn")

	ole32            = syscall.NewLazyDLL("ole32.dll")
	coInitialize     = ole32.NewProc("CoInitialize")
	coCreateInstance = ole32.NewProc("CoCreateInstance")
)

const (
	CLSID_TaskbarList = "{56FDF344-FD6D-11d0-958A-006097C9A090}"
	IID_ITaskbarList  = "{56FDF31B-FD6D-11d0-958A-006097C9A090}"

	// Window show states
	SW_HIDE       = 0
	SW_SHOW       = 5
	SW_SHOWNORMAL = 1
	SMTO_NORMAL   = 0

	// System metrics
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1

	// Window long pointer indices
	GWL_STYLE   = ^uintptr(16 - 1)  // -16 as uintptr
	GWL_EXSTYLE = ^uintptr(20 - 1)  // -20 as uintptr

	// Window styles
	WS_OVERLAPPED  = 0x00000000
	WS_POPUP       = 0x80000000
	WS_CHILD       = 0x40000000
	WS_MINIMIZE    = 0x20000000
	WS_VISIBLE     = 0x10000000
	WS_DISABLED    = 0x08000000
	WS_CLIPSIBLINGS = 0x04000000
	WS_CLIPCHILDREN = 0x02000000
	WS_MAXIMIZE    = 0x01000000
	WS_CAPTION     = 0x00C00000
	WS_BORDER      = 0x00800000
	WS_DLGFRAME    = 0x00400000
	WS_VSCROLL     = 0x00200000
	WS_HSCROLL     = 0x00100000
	WS_SYSMENU     = 0x00080000
	WS_THICKFRAME  = 0x00040000
	WS_GROUP       = 0x00020000
	WS_TABSTOP     = 0x00010000
	WS_MINIMIZEBOX = 0x00020000
	WS_MAXIMIZEBOX = 0x00010000

	// Extended window styles
	WS_EX_DLGMODALFRAME = 0x00000001
	WS_EX_NOPARENTNOTIFY = 0x00000004
	WS_EX_TOPMOST       = 0x00000008
	WS_EX_ACCEPTFILES   = 0x00000010
	WS_EX_TRANSPARENT   = 0x00000020
	WS_EX_MDICHILD      = 0x00000040
	WS_EX_TOOLWINDOW    = 0x00000080
	WS_EX_WINDOWEDGE    = 0x00000100
	WS_EX_CLIENTEDGE    = 0x00000200
	WS_EX_CONTEXTHELP   = 0x00000400
	WS_EX_RIGHT         = 0x00001000
	WS_EX_LEFT          = 0x00000000
	WS_EX_RTLREADING    = 0x00002000
	WS_EX_LTRREADING    = 0x00000000
	WS_EX_LEFTSCROLLBAR = 0x00004000
	WS_EX_RIGHTSCROLLBAR = 0x00000000
	WS_EX_CONTROLPARENT = 0x00010000
	WS_EX_STATICEDGE    = 0x00020000
	WS_EX_APPWINDOW     = 0x00040000
	WS_EX_LAYERED       = 0x00080000

	// SetWindowPos flags
	SWP_NOSIZE      = 0x0001
	SWP_NOMOVE      = 0x0002
	SWP_NOZORDER    = 0x0004
	SWP_NOREDRAW    = 0x0008
	SWP_NOACTIVATE  = 0x0010
	SWP_FRAMECHANGED = 0x0020
	SWP_SHOWWINDOW  = 0x0040
	SWP_HIDEWINDOW  = 0x0080
	SWP_NOCOPYBITS  = 0x0100
	SWP_NOOWNERZORDER = 0x0200
	SWP_NOSENDCHANGING = 0x0400
	SWP_DRAWFRAME   = SWP_FRAMECHANGED

	// SetWindowPos Z-order
	HWND_TOP       = 0
	HWND_BOTTOM    = 1
	HWND_TOPMOST   = ^uintptr(0)
	HWND_NOTOPMOST = ^uintptr(1)

	// GetWindow constants
	GW_HWNDFIRST = 0
	GW_HWNDLAST  = 1
	GW_HWNDNEXT  = 2
	GW_HWNDPREV  = 3
	GW_OWNER     = 4
	GW_CHILD     = 5

	// Process access rights
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_READ          = 0x0010

	// LayeredWindow attributes
	LWA_COLORKEY = 0x00000001
	LWA_ALPHA    = 0x00000002

	// DWM Blur Behind flags
	DWM_BB_ENABLE     = 0x00000001
	DWM_BB_BLURREGION = 0x00000002
	DWM_BB_TRANSITIONONMAXIMIZED = 0x00000004

	// Window placement flags
	WPF_SETMINPOSITION = 0x0001
	WPF_RESTORETOMAXIMIZED = 0x0002
	WPF_ASYNCWINDOWPLACEMENT = 0x0004

	// Desktop messages (from test.cpp)
	WM_DESKTOP_RAISE = 0x052C
)

// Windows structures
type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type POINT struct {
	X int32
	Y int32
}

type WINDOWPLACEMENT struct {
	Length           uint32
	Flags            uint32
	ShowCmd          uint32
	PtMinPosition    POINT
	PtMaxPosition    POINT
	RcNormalPosition RECT
}

type MARGINS struct {
	CxLeftWidth    int32
	CxRightWidth   int32
	CyTopHeight    int32
	CyBottomHeight int32
}

type DWM_BLURBEHIND struct {
	DwFlags                uint32
	FEnable                int32
	HRgnBlur               uintptr
	FTransitionOnMaximized int32
}

var (
	workerw     uintptr
	target      uintptr
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

// RaiseDesktop 提升桌面 - 基于test.cpp的实现，适用于Win10/11
func RaiseDesktop(hProgmanWnd uintptr) bool {
	if hProgmanWnd == 0 {
		return false
	}

	var res0, res1, res2, res3 uintptr

	// Call CDesktopBrowser::_IsDesktopWallpaperInitialized
	ret0, _, _ := sendMessageTimeout.Call(hProgmanWnd, WM_DESKTOP_RAISE, 0xA, 0, SMTO_NORMAL, 1000, uintptr(unsafe.Pointer(&res0)))
	if ret0 == 0 || res0 != 0 {
		log.Printf("[桌面穿透] 桌面壁纸初始化检查失败: ret=%d, res=%d", ret0, res0)
		return false
	}

	// Prepare to generate wallpaper window
	sendMessageTimeout.Call(hProgmanWnd, WM_DESKTOP_RAISE, 0xD, 0, SMTO_NORMAL, 1000, uintptr(unsafe.Pointer(&res1)))
	sendMessageTimeout.Call(hProgmanWnd, WM_DESKTOP_RAISE, 0xD, 1, SMTO_NORMAL, 1000, uintptr(unsafe.Pointer(&res2)))
	// "Animate desktop", which will make sure the wallpaper window is there
	sendMessageTimeout.Call(hProgmanWnd, WM_DESKTOP_RAISE, 0, 0, SMTO_NORMAL, 1000, uintptr(unsafe.Pointer(&res3)))

	success := res1 == 0 && res2 == 0 && res3 == 0
	log.Printf("[桌面穿透] RaiseDesktop 结果: res1=%d, res2=%d, res3=%d, success=%v", res1, res2, res3, success)
	return success
}

// IsExplorerWorker 检查窗口是否为Explorer的WorkerW窗口
func IsExplorerWorker(hwnd uintptr) bool {
	var className [256]uint16
	ret, _, _ := getClassName.Call(hwnd, uintptr(unsafe.Pointer(&className[0])), 256)
	if ret == 0 {
		return false
	}

	classNameStr := syscall.UTF16ToString(className[:])
	if classNameStr != "WorkerW" {
		return false
	}

	var pid uint32
	getWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid == 0 {
		return false
	}

	hProc, _, _ := openProcess.Call(PROCESS_QUERY_INFORMATION|PROCESS_VM_READ, 0, uintptr(pid))
	if hProc == 0 {
		return false
	}
	defer closeHandle.Call(hProc)

	var exeName [260]uint16 // MAX_PATH
	ret, _, _ = getModuleBaseNameW.Call(hProc, 0, uintptr(unsafe.Pointer(&exeName[0])), 260)
	if ret == 0 {
		return false
	}

	exeNameStr := syscall.UTF16ToString(exeName[:])
	// 转换为小写进行比较
	exeNameLower := ""
	for _, r := range exeNameStr {
		if r >= 'A' && r <= 'Z' {
			exeNameLower += string(r + 32)
		} else {
			exeNameLower += string(r)
		}
	}

	return len(exeNameLower) > 0 && (exeNameLower == "explorer.exe" || 
		(len(exeNameLower) >= 12 && exeNameLower[len(exeNameLower)-12:] == "explorer.exe"))
}

// ConfigureWindowForDesktop 配置窗口样式以适合桌面嵌入
func ConfigureWindowForDesktop(hEmbedWnd uintptr) bool {
	if hEmbedWnd == 0 {
		return false
	}

	// 获取当前窗口样式
	styleTw, _, _ := getWindowLongPtr.Call(hEmbedWnd, GWL_STYLE)
	exstyleTw, _, _ := getWindowLongPtr.Call(hEmbedWnd, GWL_EXSTYLE)

	log.Printf("[桌面穿透] 原始样式: style=0x%x, exstyle=0x%x", styleTw, exstyleTw)

	// 修改扩展样式
	if (exstyleTw & WS_EX_LAYERED) == 0 {
		exstyleTw |= WS_EX_LAYERED
	}
	if (exstyleTw & WS_EX_TOOLWINDOW) != 0 {
		exstyleTw &= ^uintptr(WS_EX_TOOLWINDOW)
	}

	// 修改基本样式
	if (styleTw & WS_CHILD) != 0 {
		styleTw &= ^uintptr(WS_CHILD)
	}
	if (styleTw & WS_POPUP) != 0 {
		styleTw &= ^uintptr(WS_POPUP)
	}
	if (styleTw & WS_OVERLAPPED) != 0 {
		styleTw &= ^uintptr(WS_OVERLAPPED)
	}
	if (styleTw & WS_CAPTION) != 0 {
		styleTw &= ^uintptr(WS_CAPTION)
	}
	if (styleTw & WS_BORDER) != 0 {
		styleTw &= ^uintptr(WS_BORDER)
	}
	if (styleTw & WS_SYSMENU) != 0 {
		styleTw &= ^uintptr(WS_SYSMENU)
	}
	if (styleTw & WS_THICKFRAME) != 0 {
		styleTw &= ^uintptr(WS_THICKFRAME)
	}

	// 应用新样式
	setWindowLongPtr.Call(hEmbedWnd, GWL_STYLE, styleTw)
	setWindowLongPtr.Call(hEmbedWnd, GWL_EXSTYLE, exstyleTw)

	log.Printf("[桌面穿透] 新样式: style=0x%x, exstyle=0x%x", styleTw, exstyleTw)
	return true
}

// SetupDesktopTransparency 设置桌面透明效果
func SetupDesktopTransparency(hEmbedWnd uintptr) bool {
	if hEmbedWnd == 0 {
		return false
	}

	// 设置为顶级窗口
	setParent.Call(hEmbedWnd, 0)

	// DWM扩展框架到客户区域实现透明
	margins := MARGINS{0, 0, -1, -1}
	ret1, _, _ := dwmExtendFrameIntoClientArea.Call(hEmbedWnd, uintptr(unsafe.Pointer(&margins)))

	// 创建模糊效果
	hRgn, _, _ := createRectRgn.Call(0, 0, ^uintptr(0), ^uintptr(0))
	bb := DWM_BLURBEHIND{
		DwFlags:  DWM_BB_ENABLE | DWM_BB_BLURREGION,
		FEnable:  1,
		HRgnBlur: hRgn,
	}
	ret2, _, _ := dwmEnableBlurBehindWindow.Call(hEmbedWnd, uintptr(unsafe.Pointer(&bb)))

	// 设置分层窗口属性
	ret3, _, _ := setLayeredWindowAttributes.Call(hEmbedWnd, 0, 0xFF, LWA_ALPHA)

	log.Printf("[桌面穿透] 透明效果设置: DwmExtend=%d, DwmBlur=%d, LayeredAttr=%d", ret1, ret2, ret3)
	return ret1 == 0 && ret2 == 0 && ret3 != 0
}

// SetDesktopFullscreen 设置窗口为全屏桌面大小
func SetDesktopFullscreen(hEmbedWnd uintptr) bool {
	if hEmbedWnd == 0 {
		return false
	}

	// 获取屏幕尺寸
	rcx, _, _ := getSystemMetrics.Call(SM_CXSCREEN)
	rcy, _, _ := getSystemMetrics.Call(SM_CYSCREEN)

	if rcx <= 0 || rcy <= 0 {
		log.Printf("[桌面穿透] 获取屏幕尺寸失败: %dx%d", rcx, rcy)
		return false
	}

	// 获取当前窗口样式用于调整窗口矩形
	styleTw, _, _ := getWindowLongPtr.Call(hEmbedWnd, GWL_STYLE)
	
	rcFullScreen := RECT{0, 0, int32(rcx), int32(rcy)}
	adjustWindowRect.Call(uintptr(unsafe.Pointer(&rcFullScreen)), styleTw, 0)

	rcfx := uint32(rcFullScreen.Right - rcFullScreen.Left)
	rcfy := uint32(rcFullScreen.Bottom - rcFullScreen.Top)

	// 移动并调整窗口大小
	ret1, _, _ := moveWindow.Call(hEmbedWnd, uintptr(rcFullScreen.Left), uintptr(rcFullScreen.Top), uintptr(rcfx), uintptr(rcfy), 1)

	// 设置窗口位置信息
	wp := WINDOWPLACEMENT{
		Length:           uint32(unsafe.Sizeof(WINDOWPLACEMENT{})),
		Flags:            WPF_SETMINPOSITION,
		ShowCmd:          SW_SHOWNORMAL,
		PtMinPosition:    POINT{rcFullScreen.Left, rcFullScreen.Top},
		PtMaxPosition:    POINT{rcFullScreen.Left, rcFullScreen.Top},
		RcNormalPosition: rcFullScreen,
	}
	ret2, _, _ := setWindowPlacement.Call(hEmbedWnd, uintptr(unsafe.Pointer(&wp)))

	log.Printf("[桌面穿透] 全屏设置: 尺寸=%dx%d, MoveWindow=%d, SetPlacement=%d", rcx, rcy, ret1, ret2)
	return ret1 != 0 && ret2 != 0
}

// AdvancedSetDesktop 高级桌面设置 - 基于test.cpp实现
func AdvancedSetDesktop(hEmbedWnd uintptr) bool {
	if hEmbedWnd == 0 {
		log.Printf("[桌面穿透] 无效的窗口句柄")
		return false
	}

	// 1. 配置窗口样式
	if !ConfigureWindowForDesktop(hEmbedWnd) {
		log.Printf("[桌面穿透] 配置窗口样式失败")
		return false
	}

	// 2. 查找桌面窗口
	hTopDeskWnd, _, _ := findWindow.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Progman"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Program Manager"))),
	)
	if hTopDeskWnd == 0 {
		log.Printf("[桌面穿透] 未找到桌面顶级窗口")
		return false
	}

	// 3. 提升桌面
	if !RaiseDesktop(hTopDeskWnd) {
		log.Printf("[桌面穿透] 提升桌面失败")
		return false
	}

	// 4. 查找SHELLDLL_DefView
	var hShellDefView uintptr
	hShellDefView, _, _ = findWindowEx.Call(hTopDeskWnd, 0, 
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("SHELLDLL_DefView"))), 0)

	var hWorker2 uintptr
	if hShellDefView == 0 {
		// 回退到23H2搜索模式
		log.Printf("[桌面穿透] 使用23H2兼容模式搜索")
		hWorker_p1, _, _ := getWindow.Call(hTopDeskWnd, GW_HWNDPREV)
		if hWorker_p1 != 0 {
			hShellDefView, _, _ = findWindowEx.Call(hWorker_p1, 0,
				uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("SHELLDLL_DefView"))), 0)
			if hShellDefView == 0 {
				hWorker2 = hWorker_p1
				hWorker_p2, _, _ := getWindow.Call(hWorker_p1, GW_HWNDPREV)
				if hWorker_p2 != 0 {
					hShellDefView, _, _ = findWindowEx.Call(hWorker_p2, 0,
						uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("SHELLDLL_DefView"))), 0)
					// hWorker1 was unused, so we removed it
				}
			}
		}
	}

	if hShellDefView == 0 {
		log.Printf("[桌面穿透] 未找到桌面shell defview窗口")
		return false
	}

	// 5. 查找Worker窗口
	hWorker, _, _ := findWindowEx.Call(hTopDeskWnd, 0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("WorkerW"))), 0)
	if hWorker == 0 {
		hWorker, _, _ = findWindowEx.Call(hTopDeskWnd, 0,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("WorkerA"))), 0)
	}

	// 6. 确定版本和Worker窗口
	bIsVersion1_2 := false
	if hWorker == 0 {
		if hWorker2 != 0 {
			hWorker = hWorker2
		} else {
			hWorker = hTopDeskWnd
		}
		bIsVersion1_2 = true
		log.Printf("[桌面穿透] 使用版本1.2兼容模式")
	}

	// 7. 设置透明效果
	if !SetupDesktopTransparency(hEmbedWnd) {
		log.Printf("[桌面穿透] 设置透明效果失败")
	}

	// 8. 设置父窗口和窗口层次
	var parentWnd uintptr
	if bIsVersion1_2 {
		parentWnd = hWorker
	} else {
		parentWnd = hTopDeskWnd
	}
	setParent.Call(hEmbedWnd, parentWnd)

	// 9. 调整窗口Z序
	setWindowPos.Call(hEmbedWnd, HWND_TOP, 0, 0, 0, 0,
		SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE|SWP_DRAWFRAME)
	setWindowPos.Call(hShellDefView, HWND_TOP, 0, 0, 0, 0,
		SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE)
	setWindowPos.Call(hWorker, HWND_BOTTOM, 0, 0, 0, 0,
		SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE|SWP_DRAWFRAME)

	// 10. 设置全屏
	if !SetDesktopFullscreen(hEmbedWnd) {
		log.Printf("[桌面穿透] 设置全屏失败")
	}

	// 11. 显示窗口
	showWindow.Call(hTopDeskWnd, SW_SHOW)
	showWindow.Call(hEmbedWnd, SW_SHOW)
	showWindow.Call(hWorker, SW_SHOW)

	log.Printf("[桌面穿透] 高级桌面设置完成: Parent=0x%x, ShellDefView=0x%x, Worker=0x%x, Version1.2=%v", 
		parentWnd, hShellDefView, hWorker, bIsVersion1_2)

	return true
}

// SetDesktop 将窗口设置为桌面壁纸 - 保持向后兼容
func SetDesktop(hwnd uintptr) {
	// 使用新的高级设置方法
	AdvancedSetDesktop(hwnd)
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

	var pTaskbar **struct {
		vtbl *struct {
			QueryInterface uintptr
			AddRef         uintptr
			Release        uintptr
			HrInit         uintptr
			AddTab         uintptr
			DeleteTab      uintptr
			ActivateTab    uintptr
			SetActiveAlt   uintptr
		}
	}

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

// EnsureEmbedWindowBelow 确保嵌入窗口在ShellDefView下方 - 基于test.cpp实现
func EnsureEmbedWindowBelow(hShellDefView, hEmbedWnd uintptr) (bool, uintptr) {
	if hShellDefView == 0 || hEmbedWnd == 0 {
		return false, 0
	}

	prev, _, _ := getWindow.Call(hEmbedWnd, GW_HWNDPREV)
	if prev == hShellDefView {
		// 顺序已正确，无需操作
		return true, 0
	}

	// 修复Z序：将嵌入窗口移动到ShellDefView下方
	setWindowPos.Call(hEmbedWnd, hShellDefView, 0, 0, 0, 0,
		SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE)

	if IsExplorerWorker(prev) {
		return true, 0
	}

	return false, prev
}

// StartZOrderMonitoring 启动Z序监控 - 基于test.cpp的24H2兼容性实现
func StartZOrderMonitoring(hShellDefView, hEmbedWnd uintptr, ctx context.Context) {
	if hShellDefView == 0 || hEmbedWnd == 0 {
		log.Printf("[Z序监控] 无效的窗口句柄，跳过监控")
		return
	}

	log.Printf("[Z序监控] 开始监控嵌入窗口和ShellDefView之间的Z序...")
	
	const maxConsecutiveFixes = 5
	consecutiveFixCount := 0
	var lastConflictHwnd uintptr

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("[Z序监控] 监控已停止")
				return
			case <-ticker.C:
				ok, conflictHwnd := EnsureEmbedWindowBelow(hShellDefView, hEmbedWnd)

				if !ok {
					if conflictHwnd == 0 {
						consecutiveFixCount = 0
					} else if conflictHwnd == lastConflictHwnd {
						consecutiveFixCount++
					} else {
						lastConflictHwnd = conflictHwnd
						consecutiveFixCount = 1
					}

					if conflictHwnd != 0 {
						log.Printf("[Z序监控] 检测到冲突窗口: 0x%x, 连续修复次数: %d", conflictHwnd, consecutiveFixCount)
					}

					if consecutiveFixCount >= maxConsecutiveFixes {
						log.Printf("[Z序监控] 检测到重复的Z序冲突！退出监控以保护系统")
						// 这里可以选择发送关闭消息或其他保护措施
						return
					}
				} else {
					// 没冲突，重置计数器
					consecutiveFixCount = 0
					lastConflictHwnd = 0
				}
			}
		}
	}()
}

// SetupAdvancedWallpaper 设置高级壁纸功能 - 整合所有增强功能
func SetupAdvancedWallpaper(windowTitle string) bool {
	hwnd := FindWindowByTitle(windowTitle)
	if hwnd == 0 {
		log.Printf("[桌面穿透] 未找到窗口: %s", windowTitle)
		return false
	}

	log.Printf("[桌面穿透] 找到窗口: %s (0x%x)", windowTitle, hwnd)

	// 1. 从任务栏移除
	if err := RemoveFromTaskbar(hwnd); err != nil {
		log.Printf("[桌面穿透] 从任务栏移除失败: %v", err)
	}

	// 2. 使用高级桌面设置
	if !AdvancedSetDesktop(hwnd) {
		log.Printf("[桌面穿透] 高级桌面设置失败")
		return false
	}

	// 3. 启动Z序监控（仅在非版本1.2模式下）
	// 查找ShellDefView用于监控
	hTopDeskWnd, _, _ := findWindow.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Progman"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Program Manager"))),
	)
	if hTopDeskWnd != 0 {
		hShellDefView, _, _ := findWindowEx.Call(hTopDeskWnd, 0, 
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("SHELLDLL_DefView"))), 0)
		
		if hShellDefView != 0 {
			// 创建监控上下文
			monitorCtx, _ := context.WithCancel(context.Background())
			StartZOrderMonitoring(hShellDefView, hwnd, monitorCtx)
		} else {
			log.Printf("[桌面穿透] 未找到ShellDefView，跳过Z序监控")
		}
	}

	log.Printf("[桌面穿透] 高级壁纸设置完成")
	return true
}

// SetupWallpaper 设置壁纸 - 保持向后兼容
func SetupWallpaper(windowTitle string) bool {
	return SetupAdvancedWallpaper(windowTitle)
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

// SetDPIAware 设置DPI感知，返回是否成功
func SetDPIAware() bool {
	const (
		DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = ^uintptr(3)
	)

	proc := user32.NewProc("SetProcessDpiAwarenessContext")
	ret, _, _ := proc.Call(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2)
	return ret != 0
}
