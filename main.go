package main

import (
	"io"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/pelletier/go-toml"
	"github.com/zserge/lorca"
)

// ====== WinAPI相关区域 ======
// Windows API相关的Go代码已移至 winapi.go
// ====== WinAPI相关区域结束 ======

// ====== 全局变量和常量区 ======
var (
	mainWindow lorca.UI
	isRunning  bool
	urlStr     string
	BwPath     string
	lorcaname  string
	logFile    *os.File
	t          *time.Ticker
)

const (
	title = "ClassPaper"
)

// ====== 工具函数区 ======
// 路径中的中文字符编码
func encodeChineseCharacters(path string) string {
	re := regexp.MustCompile("[\u4e00-\u9fa5]")
	return re.ReplaceAllStringFunc(path, func(s string) string {
		return url.QueryEscape(s)
	})
}

// 获取本地文件路径的 file:// URL
func getFilePathURL(relativePath string) (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	absPath := filepath.Join(exeDir, relativePath)
	absPath = encodeChineseCharacters(absPath)
	absPath = filepath.ToSlash(absPath)
	fileURL := "file:///" + absPath
	return fileURL, nil
}

// 生成随机字符串
func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// 配置结构体及读取
// 使用 TOML 格式
// config.toml 示例：
// [default]
// url = "..."
// browser_path = "..."
type Config struct {
	URL         string `toml:"url"`
	BrowserPath string `toml:"browser_path"`
}

func ParseConfig() (*Config, error) {
	// 读取 config.toml
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return nil, err
	}
	config := &Config{}
	// 直接解析 toml 到结构体
	err = toml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// URL标准化
func NormalizeURL(url string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	url, _ = getFilePathURL(url)
	return url
}

// ====== 桌面壁纸/窗口相关函数区 ======
func setWallpaper() {
	ret := SetupWallpaper(lorcaname)
	log.Printf("[桌面穿透] SetupWallpaper(%s) 返回: %v", lorcaname, ret)
	t = time.NewTicker(time.Second)
	go func() {
		for range t.C {
			log.Println("[桌面穿透] 定时调用 RemoveFromTaskbar 保持窗口状态")
			if hwnd := FindWindowByTitle(lorcaname); hwnd != 0 {
				RemoveFromTaskbar(hwnd)
			}
		}
	}()
}

func runLorcaUI() {
	ui, err := lorca.New(urlStr, "", BwPath, 0, 0, "--kiosk")
	if err != nil {
		log.Printf("[Lorca] 创建UI失败: %v", err)
		return
	}
	mainWindow = ui
	lorcaname = "classpaper" + generateRandomString(6)
	ui.Eval("document.title='" + lorcaname + "'")
	log.Printf("[Lorca] 设置窗口标题: %s", lorcaname)

	// 绑定前端可调用的Go函数
	ui.Bind("getWidth", GetScreenWidth)
	ui.Bind("getHeight", GetScreenHeight)
	ui.Bind("readFile", func(path string) (string, error) {
		data, err := os.ReadFile(path)
		return string(data), err
	})
	ui.Bind("writeFile", func(path, content string) error {
		return os.WriteFile(path, []byte(content), 0644)
	})
	ui.Bind("readDir", func(dir string) ([]string, error) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		return names, nil
	})

	time.Sleep(time.Millisecond * 300)
	setWallpaper()
	log.Println("[Lorca] 等待窗口关闭...")
	<-ui.Done()
	mainWindow = nil
	log.Println("[Lorca] 窗口已关闭")
}

// ====== 托盘菜单相关函数区 ======
func reloadPage() {
	log.Println("[托盘] 触发网页重载")
	mainWindow.Eval("location.reload(true)")
}

func restartWebpageDisplayProgram() {
	log.Println("[托盘] 重启网页显示程序...")
	if mainWindow != nil {
		mainWindow.Close()
	}
	go runLorcaUI()
}

func onReady() {
	systray.SetTemplateIcon(IconData, IconData)
	systray.SetTitle(title)
	systray.SetTooltip(title)
	reloadMenuItem := systray.AddMenuItem("重载网页", "Reload Page")
	setPenetrationMenuItem := systray.AddMenuItem("设置程序桌面穿透", "Set Window Penetration")
	restartWebpageMenuItem := systray.AddMenuItem("重启网页显示程序", "Restart Webpage Display")
	settingsMenuItem := systray.AddMenuItem("设置", "Open Settings")
	restartMenuItem := systray.AddMenuItem("重启程序", "Restart Application")
	quitMenuItem := systray.AddMenuItem("退出程序", "Quit Application")
	log.Println("[托盘] 托盘菜单已初始化")
	go runLorcaUI()
	go func() {
		for {
			select {
			case <-reloadMenuItem.ClickedCh:
				log.Println("[托盘] 手动触发网页重载")
				reloadPage()
			case <-setPenetrationMenuItem.ClickedCh:
				log.Println("[托盘] 手动触发桌面穿透")
				setWallpaper()
			case <-restartWebpageMenuItem.ClickedCh:
				log.Println("[托盘] 手动重启网页显示程序")
				restartWebpageDisplayProgram()
			case <-settingsMenuItem.ClickedCh:
				log.Println("[托盘] 打开设置窗口")
				openSettings()
			case <-restartMenuItem.ClickedCh:
				log.Println("[托盘] 重启主程序...")
				restartProgram()
				systray.Quit()
			case <-quitMenuItem.ClickedCh:
				log.Println("[托盘] 退出程序")
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	log.Println("[退出] 程序退出，清理资源...")
	if mainWindow != nil {
		mainWindow.Close()
	}
	if logFile != nil {
		logFile.Sync()
	}
	if t != nil {
		t.Stop()
	}
}

func openSettings() {
	execPath, err := os.Executable()
	if err != nil {
		log.Printf("[设置] 获取可执行文件路径失败: %v", err)
		return
	}
	execDir := filepath.Dir(execPath)
	cmd := exec.Command("./setting")
	cmd.Dir = execDir
	err = cmd.Start()
	if err != nil {
		log.Printf("[设置] 启动设置窗口失败: %v", err)
	}
}

func restartProgram() {
	execPath, err := os.Executable()
	if err != nil {
		log.Printf("[重启] 获取可执行文件路径失败: %v", err)
		return
	}
	execDir := filepath.Dir(execPath)
	cmd := exec.Command(execPath)
	cmd.Dir = execDir
	err = cmd.Start()
	if err != nil {
		log.Printf("[重启] 启动新进程失败: %v", err)
	}
}

func endup() {
	log.Println("[退出] 执行endup，关闭窗口和资源...")
	if mainWindow != nil {
		mainWindow.Close()
	}
	systray.Quit()
	if logFile != nil {
		logFile.Sync()
	}
	if t != nil {
		t.Stop()
	}
}

// ====== 主程序入口和初始化 ======
func main() {
	defer endup()
	var err error
	logFile, err = os.Create("app.log")
	if err != nil {
		log.Fatalf("[启动] 创建日志文件失败: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	config, err := ParseConfig()
	if err != nil {
		log.Printf("[启动] 读取配置失败: %v", err)
		return
	}
	log.Printf("[启动] 加载配置URL: %s", config.URL)
	urlStr = NormalizeURL(config.URL)
	BwPath = config.BrowserPath
	log.Printf("[启动] 标准化URL: %s", urlStr)
	log.Printf("[启动] 浏览器路径: %s", BwPath)
	systray.Run(onReady, onExit)
}

func init() {
	if runtime.GOOS == "windows" {
		os.Setenv("WINGUI_NO_CONSOLE", "1")
		exec.Command("cmd", "/c", "chcp", "65001").Run()
	}
}

// ====== END ======
