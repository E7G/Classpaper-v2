package main

import (
	"fmt"
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
	"encoding/json"
)

// ====== WinAPI相关区域 ======
// Windows API相关的Go代码已移至 winapi.go
// ====== WinAPI相关区域结束 ======

// ====== 全局变量和常量区 ======
var (
	mainWindow lorca.UI
	settingsWindow lorca.UI
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
	Default struct {
		URL         string `toml:"url"`
		BrowserPath string `toml:"browser_path"`
	} `toml:"default"`
}

func ParseConfig() (*Config, error) {
	// 读取 config.toml
	data, err := os.ReadFile("config.toml")
	if err != nil {
		log.Printf("[启动] 配置文件不存在，创建默认配置")
		// 创建默认配置
		defaultConfig := &Config{}
		defaultConfig.Default.URL = "./res/index.html"
		defaultConfig.Default.BrowserPath = ""
		
		// 将默认配置写入文件
		configData, err := toml.Marshal(defaultConfig)
		if err != nil {
			return nil, fmt.Errorf("生成默认配置失败: %v", err)
		}
		
		err = os.WriteFile("config.toml", configData, 0644)
	if err != nil {
			return nil, fmt.Errorf("写入默认配置失败: %v", err)
		}
		
		return defaultConfig, nil
	}

	log.Println("[启动] 读取配置文件: config.toml")
	log.Println("[启动] 配置文件内容:", string(data))
	
	config := &Config{}
	// 直接解析 toml 到结构体
	err = toml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if config.Default.URL == "" {
		config.Default.URL = "./res/index.html"
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
		failCount := 0
		for range t.C {
			if hwnd := FindWindowByTitle(lorcaname); hwnd != 0 {
				err := RemoveFromTaskbar(hwnd)
				if err != nil {
					failCount++
					if failCount == 1 || failCount%60 == 0 {
						log.Printf("[桌面穿透] RemoveFromTaskbar 失败: %v（累计 %d 次）", err, failCount)
					}
				} else {
					if failCount > 0 {
						log.Printf("[桌面穿透] RemoveFromTaskbar 恢复正常")
					}
					failCount = 0
				}
			}
		}
	}()
}

func runLorcaUI() {
	ui, err := lorca.New(urlStr, "", BwPath, 0, 0, "--kiosk", "--autoplay-policy=no-user-gesture-required")
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
	// 判断系统是否为夜间模式，选择不同的图标
	var iconLight, iconDark []byte
	iconLight = IconDataLight
	iconDark = IconDataDark

	// 检测夜间模式
	isDarkMode := false
	// Windows下可用注册表检测，macOS可用AppleScript，Linux需根据桌面环境
	isDarkMode = GetWindowsDarkMode()

	if isDarkMode {
		systray.SetTemplateIcon(iconDark, iconDark)
	} else {
		systray.SetTemplateIcon(iconLight, iconLight)
	}

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

func openSettings() {
	if settingsWindow != nil {
		log.Println("[设置] 设置窗口已经打开")
		return
	}

	settingsURL, err := getFilePathURL("res/settings.html")
	if err != nil {
		log.Printf("[设置] 获取设置页面URL失败: %v", err)
		return
	}

	// 获取屏幕尺寸
	screenWidth := GetScreenWidth()
	screenHeight := GetScreenHeight()

	// 设置窗口尺寸（屏幕宽度的60%，高度的80%）
	width := int(float64(screenWidth) * 0.6)
	height := int(float64(screenHeight) * 0.8)

	// 确保窗口尺寸在合理范围内
	if width < 1000 {
		width = 1000 // 最小宽度
	}
	if height < 800 {
		height = 800 // 最小高度
	}

	// 计算窗口位置（居中）
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2

	// 创建窗口并设置位置
	ui, err := lorca.New(settingsURL, "", "", width, height)
	if err != nil {
		log.Printf("[设置] 创建设置窗口失败: %v", err)
		return
	}

	// 设置窗口位置
	ui.Eval(fmt.Sprintf(`
		window.moveTo(%d, %d);
		document.title = "ClassPaper 设置";
	`, x, y))

	settingsWindow = ui

	// 立即绑定函数
	log.Println("[设置] 开始绑定函数...")
	
	// 绑定配置相关函数
	err = ui.Bind("readConfig", func() (interface{}, error) {
		config, err := ParseConfig()
		if err != nil {
			log.Printf("[设置] 读取配置失败: %v", err)
			return nil, err
		}
		log.Printf("[设置] 读取到配置: %+v", config)
		return map[string]interface{}{
			"Default": map[string]interface{}{
				"URL":         config.Default.URL,
				"BrowserPath": config.Default.BrowserPath,
			},
		}, nil
	})
	if err != nil {
		log.Printf("[设置] 绑定readConfig失败: %v", err)
	}

	err = ui.Bind("saveConfig", func(configJSON string) error {
		log.Printf("[设置] 保存配置JSON: %s", configJSON)
		var config Config
		err := json.Unmarshal([]byte(configJSON), &config)
		if err != nil {
			log.Printf("[设置] 解析配置JSON失败: %v", err)
			return fmt.Errorf("解析配置JSON失败: %v", err)
		}

		// 将配置写入文件
		configData, err := toml.Marshal(config)
		if err != nil {
			log.Printf("[设置] 序列化配置失败: %v", err)
			return fmt.Errorf("序列化配置失败: %v", err)
		}

		err = os.WriteFile("config.toml", configData, 0644)
		if err != nil {
			log.Printf("[设置] 写入配置文件失败: %v", err)
			return fmt.Errorf("写入配置文件失败: %v", err)
		}

		log.Printf("[设置] 配置已保存: %+v", config)

		// 更新当前运行的配置
		urlStr = NormalizeURL(config.Default.URL)
		BwPath = config.Default.BrowserPath

		return nil
	})
	if err != nil {
		log.Printf("[设置] 绑定saveConfig失败: %v", err)
	}

	// 绑定文件写入函数
	err = ui.Bind("writeFile", func(path string, content string) error {
		log.Printf("[设置] 写入文件 %s", path)
		return os.WriteFile(path, []byte(content), 0644)
	})
	if err != nil {
		log.Printf("[设置] 绑定writeFile失败: %v", err)
	}

	// 绑定壁纸扫描函数
	err = ui.Bind("scanWallpaperDir", scanWallpaperDir)
	if err != nil {
		log.Printf("[设置] 绑定scanWallpaperDir失败: %v", err)
	}

	// 绑定主窗口刷新函数
	err = ui.Bind("reloadMainWindow", func() error {
		if mainWindow != nil {
			log.Println("[设置] 刷新主窗口")
			mainWindow.Eval("location.reload(true)")
			return nil
		}
		return fmt.Errorf("主窗口未打开")
	})
	if err != nil {
		log.Printf("[设置] 绑定reloadMainWindow失败: %v", err)
	}

	log.Println("[设置] 函数绑定完成")

	// 监听窗口关闭
	go func() {
		<-ui.Done()
		settingsWindow = nil
		log.Println("[设置] 设置窗口已关闭")
	}()
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
	if settingsWindow != nil {
		settingsWindow.Close()
	}
	systray.Quit()
	if logFile != nil {
	logFile.Sync()
	}
	if t != nil {
	t.Stop()
}
}

func onExit() {
	log.Println("[退出] 程序退出，清理资源...")
	if mainWindow != nil {
		mainWindow.Close()
	}
	if settingsWindow != nil {
		settingsWindow.Close()
	}
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
	log.Printf("[启动] 加载配置URL: %s", config.Default.URL)
	urlStr = NormalizeURL(config.Default.URL)
	BwPath = config.Default.BrowserPath
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

// 扫描壁纸文件夹
func scanWallpaperDir() ([]string, error) {
	wallpapers := []string{}
	
	// 扫描res/wallpaper目录
	files, err := os.ReadDir("res/wallpaper")
	if err != nil {
		return nil, fmt.Errorf("读取壁纸目录失败: %v", err)
	}

	// 添加文件到列表，使用相对于index.html的路径
	for _, file := range files {
		if !file.IsDir() {
			// 检查是否是图片文件
			name := file.Name()
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
				// 使用相对于index.html的路径
				wallpapers = append(wallpapers, "wallpaper/"+name)
			}
		}
	}

	return wallpapers, nil
}
