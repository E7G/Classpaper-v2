package main

import (
	"os"
	"os/exec"
	"net/url"
	"regexp"
	"path/filepath"
	"runtime"
	"math/rand"
	"time"
	"strings" 
 //   "syscall" 
	"github.com/getlantern/systray"
	"github.com/zserge/lorca"
//	"github.com/AllenDang/w32"
//	"github.com/hnakamur/w32syscall"
	"io"
	"log"
//	"fmt"

	"gopkg.in/ini.v1"
)



var (
	mainWindow lorca.UI
	isRunning  bool
	urlStr string
	BwPath string
//	lorcaui w32.HWND
	lorcaname string
	logFile *os.File
)

const (
	title      = "ClassPaper"
	//relativePath     = "./res/index.html" // 在这里指定您的本地网页文件路径

)

func encodeChineseCharacters(path string) string {
	// 使用正则表达式匹配中文字符，并进行编码
	re := regexp.MustCompile("[\u4e00-\u9fa5]")
	encodedPath := re.ReplaceAllStringFunc(path, func(s string) string {
		return url.QueryEscape(s)
	})
	return encodedPath
}

func getFilePathURL(relativePath string) (string, error) {
	// 获取程序执行路径
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// 获取程序所在目录
	exeDir := filepath.Dir(exePath)

	// 使用filepath.Join确保正确拼接路径
	absPath := filepath.Join(exeDir, relativePath)

	// 将路径中的中文字符进行编码
	absPath = encodeChineseCharacters(absPath)

	// 将路径中的正斜杠转换为反斜杠
	absPath = filepath.ToSlash(absPath)

	// 转换为file:// URL
	fileURL := "file:///" + absPath

	return fileURL, nil
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
/*
func getWindowHandlesWithSubstring(substring string) ([]w32.HWND, error) {
	var windowHandles []w32.HWND

	// Enumerate all top-level windows
	err := w32syscall.EnumWindows(func(hwnd syscall.Handle, lparam uintptr) bool { 
		h:=w32.HWND(hwnd)
		// Get the length of the window title
		length := w32.GetWindowTextLength(h)
		if length > 0 {
			// Get the window title
			title := w32.GetWindowText(h)
			// Check if the title contains the specified substring
			if strings.Contains(title, substring) {
				windowHandles = append(windowHandles, h)
			}
		}
		return true // 继续枚举
	}, 0)

	if err != nil {
		return nil, err
	}

	return windowHandles, nil
}

func setWallpaper(hwnd w32.HWND) {
*/
func setWallpaper(){
/*
	// 获取窗口原有样式
	style := w32.GetWindowLong(hwnd, w32.GWL_STYLE)

	// 设置窗口透明度为半透明
	//w32.SetLayeredWindowAttributes(hwnd, 0, 128, w32.LWA_ALPHA)

	// 设置窗口样式为 WS_EX_LAYERED|WS_EX_TRANSPARENT
	w32.SetWindowLong(hwnd, w32.GWL_EXSTYLE, w32.WS_POPUP|w32.WS_EX_LAYERED|w32.WS_EX_TRANSPARENT)

	// 将 WS_EX_TOOLWINDOW 样式添加到窗口扩展样式中
	style = w32.GetWindowLong(hwnd, w32.GWL_STYLE)
	w32.SetWindowLong(hwnd, w32.GWL_EXSTYLE, uint32(style | w32.WS_EX_TOOLWINDOW))

	// 移动窗口到桌面的底层
	//w32.SetWindowPos(hwnd, w32.HWND_BOTTOM, 0, 0, 0, 0, w32.SWP_NOMOVE|w32.SWP_NOSIZE)

	// 设置窗口样式为 WS_CHILD，即子窗口
	//style = w32.GetWindowLong(hwnd, w32.GWL_STYLE)
	//style = style & ^int32(w32.WS_CHILD) // 清除WS_CHILD和WS_POPUP样式
	//w32.SetWindowLong(hwnd, w32.GWL_STYLE, uint32(style))
	w32.SetParent(hwnd, 0) // 设置父窗口为桌面
	//style = w32.GetWindowLong(hwnd, w32.GWL_STYLE)
	//style = style | 2147483648 // 设置WS_CHILD样式
	//w32.SetWindowLong(hwnd, w32.GWL_STYLE, uint32(style))

	style = w32.GetWindowLong(hwnd, w32.GWL_STYLE)
	style = style | w32.WS_CHILD // 设置WS_CHILD样式
	w32.SetWindowLong(hwnd, w32.GWL_STYLE, uint32(style))

*/	

	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		log.Println("Failed to get executable path:", err)
		return
	}

	// 获取可执行文件的所在目录
	execDir := filepath.Dir(execPath)

	log.Println("./deltaskbar "+lorcaname)
	// 执行外部程序 
	cmd := exec.Command("./deltaskbar",lorcaname)
	
	cmd.Dir = execDir
	err = cmd.Run()
	if err != nil {
		log.Println("Error starting deltaskbar:", err)
	}

	log.Println("./setwallpaper "+lorcaname)
	// 执行外部程序 
	
	//cmd := exec.Command("./setwallpaper.exe",fmt.Sprintf("%v", hwnd))
	cmd = exec.Command("./setwallpaper",lorcaname)
	
	cmd.Dir = execDir
	err = cmd.Run()
	if err != nil {
		log.Println("Error starting setwallpaper:", err)
	}


	// 等待以允许窗口更新
	//time.Sleep(time.Second)
}

func runLorcaUI() {
	ui, err := lorca.New(urlStr, "",BwPath,0,0 ,"--kiosk") // 设置窗口大小
	if err != nil {
		log.Println("Failed to create Lorca UI:", err)
		return
	}
	mainWindow = ui

	//log.Println("Loading url :"+urlStr)

	// 设置全屏和桌面穿透
	//ui.SetBounds(lorca.Bounds{
	//	WindowState: lorca.WindowStateFullscreen,
	//})
/*
	var finalUniqueHandlesMap = make(map[w32.HWND]int)
	var currentUniqueHandlesMap = make(map[w32.HWND]bool)
	var circle=1

	// 等待以允许窗口更新
	//time.Sleep(time.Second)

	// 重复该过程3次
	for i := 0; i < circle; i++ {
		randomString := "classpaper" + generateRandomString(6)

		// 使用生成的随机字符串设置文档标题
		ui.Eval("document.title='" + randomString + "'")
		
		// 等待以允许窗口更新
		time.Sleep(time.Second)

		// 获取包含指定子字符串的窗口句柄
		handles, err := getWindowHandlesWithSubstring(randomString)

		if err != nil {
			log.Println("错误:", err)
			return
		}

		// 将句柄添加到当前循环的唯一句柄映射中
		currentUniqueHandlesMap = make(map[w32.HWND]bool)
		for _, handle := range handles {
			currentUniqueHandlesMap[handle] = true
		}

		// 将当前循环的唯一句柄映射与最终结果的映射合并
		for handle := range currentUniqueHandlesMap {
			finalUniqueHandlesMap[handle]++
		}

		// 等待以允许窗口更新
		//time.Sleep(time.Second)
	}

	// 筛选出在三次循环中都存在的唯一句柄
	var resultHandles []w32.HWND
	for handle, count := range finalUniqueHandlesMap {
		if count == circle {
			resultHandles = append(resultHandles, handle)
		}
	}

	// 打印最终结果的唯一窗口句柄
	log.Println("最终结果的唯一窗口句柄:")
	for _, handle := range resultHandles {
		log.Println("-", handle)
	}

	// 取出第一个 handle
	//var lorcaui w32.HWND
	if len(resultHandles) > 0 {
		lorcaui = resultHandles[0]
	}

	ui.Eval("document.title='ClassPaper'")

	// 清理无用的变量
	finalUniqueHandlesMap = nil
	currentUniqueHandlesMap = nil
	resultHandles = nil
*/
	// 使用 wallpaper 变量实现窗口作为壁纸
//	setWallpaper(lorcaui)
	
	lorcaname = "classpaper" + generateRandomString(6)

	// 使用生成的随机字符串设置文档标题
	ui.Eval("document.title='" + lorcaname + "'")
// 等待以允许窗口更新
//	time.Sleep(time.Millisecond*10)	
	
//	time.Sleep(time.Second*5)
	
for range time.Tick( time.Second ) {
 
	setWallpaper()
	time.Sleep(time.Millisecond*10)
}
	
	//setWallpaper()

	//time.Sleep(time.Millisecond*100)	

	//setWallpaper()

	<-ui.Done()
	


	mainWindow = nil
}


func reloadPage() {
	mainWindow.Eval("location.reload(true)")
}

func restartWebpageDisplayProgram(){
	if mainWindow != nil {
		mainWindow.Close()
	}
	go runLorcaUI()
}

func onReady() {
	// 从本地文件加载图标
	//iconPath := "./icon/light/icon.ico"
	
	systray.SetTemplateIcon(IconData, IconData)

	systray.SetTitle(title)
	systray.SetTooltip(title)
	reloadMenuItem := systray.AddMenuItem("重载网页", "Reload Page")
	setPenetrationMenuItem := systray.AddMenuItem("设置程序桌面穿透", "Set Window Penetration")
    restartWebpageMenuItem := systray.AddMenuItem("重启网页显示程序", "Restart Webpage Display")
	settingsMenuItem := systray.AddMenuItem("设置", "Open Settings") // 添加 "设置" 菜单项
	restartMenuItem := systray.AddMenuItem("重启程序", "Restart Application")
	quitMenuItem := systray.AddMenuItem("退出程序", "Quit Application")

	go runLorcaUI()

	go func() {
		for {
			select {
			case <-reloadMenuItem.ClickedCh:
				reloadPage()
			case <-setPenetrationMenuItem.ClickedCh:
    			setWallpaper() // 调用设置桌面穿透的函数
            case <-restartWebpageMenuItem.ClickedCh:
                log.Println("Restarting Webpage Display Program...")
                // 执行重启网页显示程序的命令
                restartWebpageDisplayProgram()
			case <-settingsMenuItem.ClickedCh:
				openSettings() // 打开设置
			case <-restartMenuItem.ClickedCh:
				log.Println("Restarting Program...")
				// 执行重启程序的命令
				restartProgram()
				systray.Quit()

			case <-quitMenuItem.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	if mainWindow != nil {
		mainWindow.Close()
	}

	logFile.Sync()
}

func openSettings() {
	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		log.Println("Failed to get executable path:", err)
		return
	}

	// 获取可执行文件的所在目录
	execDir := filepath.Dir(execPath)


	// 执行外部程序 setting.exe
	cmd := exec.Command("./setting")
	cmd.Dir = execDir
	err = cmd.Start()
	if err != nil {
		log.Println("Error starting setting:", err)
	}
}

func restartProgram() {
	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		log.Println("Failed to get executable path:", err)
		return
	}

	// 获取可执行文件的所在目录
	execDir := filepath.Dir(execPath)

	// 启动一个新的进程来替代当前进程
	cmd := exec.Command(execPath)
	cmd.Dir = execDir
	cmd.Start()

	// 退出当前进程
	//os.Exit(0)
	//endup()
}

func endup(){
	if mainWindow != nil {
		mainWindow.Close()
	}
	systray.Quit()
	
	logFile.Sync()
}

// Config represents the structure of the configuration file
type Config struct {
	URL          string `ini:"url"`
	BrowserPath  string `ini:"browser_path"`
}

// ParseConfig reads the configuration from config.ini file
func ParseConfig() (*Config, error) {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := cfg.Section("default").MapTo(config); err != nil {
		return nil, err
	}

	return config, nil
}

// NormalizeURL converts the URL to a format that the browser can load
func NormalizeURL(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}

	// Assuming it's a local path, convert it to a file URL
	url,_ =getFilePathURL(url)
	return url
}


func main() {
	// 在程序退出时运行 cleanup 函数
	defer endup()


	// 创建日志文件
	logFile, err := os.Create("app.log")
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}
	defer logFile.Close()

	// 设置日志输出到文件和控制台
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetOutput(logFile)

	config, err := ParseConfig()
	if err != nil {
		log.Println("Error reading config:", err)
		return
	}

	log.Println("Loading : ",config.URL)
	urlStr =NormalizeURL(config.URL)
	BwPath=config.BrowserPath

	log.Printf("Normalized URL: %s\n", urlStr)
	log.Printf("Browser Path: %s\n", BwPath)


	//logFile.Sync()
	systray.Run(onReady, onExit)

	// 在程序结束时刷新并关闭日志文件
	//logFile.Sync()
}

func init() {
	// 设置 Windows 环境下 systray 需要的一些参数
	if runtime.GOOS == "windows" {
		os.Setenv("WINGUI_NO_CONSOLE", "1")
		exec.Command("cmd", "/c", "chcp", "65001").Run()
	}
}