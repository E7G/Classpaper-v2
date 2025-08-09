package main

import (
	"context"
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

	"encoding/json"
	"sync"

	"github.com/getlantern/systray"
	"github.com/pelletier/go-toml"
	"github.com/zserge/lorca"
)

// ====== WinAPI相关区域 ======
// Windows API相关的Go代码已移至 winapi.go
// ====== WinAPI相关区域结束 ======

// ====== 全局变量和常量区 ======
var (
	mainWindow      lorca.UI
	settingsWindow  lorca.UI
	isRunning       bool
	urlStr          string
	BwPath          string
	lorcaname       string
	logFile         *os.File
	t               *time.Ticker
	windowMu        sync.Mutex // 新增互斥锁保护窗口关闭
	wallpaperCtx    context.Context
	wallpaperCancel context.CancelFunc
	settingsCtx     context.Context
	settingsCancel  context.CancelFunc
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
	log.Println("[启动] 开始读取配置文件: config.toml")
	// 读取 config.toml
	data, err := os.ReadFile("config.toml")
	if err != nil {
		log.Printf("[启动] 配置文件不存在，创建默认配置: %v", err)
		// 创建默认配置
		defaultConfig := &Config{}
		defaultConfig.Default.URL = "./res/index.html"
		defaultConfig.Default.BrowserPath = ""

		log.Println("[启动] 生成默认配置数据")
		// 将默认配置写入文件
		configData, err := toml.Marshal(defaultConfig)
		if err != nil {
			log.Printf("[启动] 生成默认配置数据失败: %v", err)
			return nil, fmt.Errorf("生成默认配置失败: %v", err)
		}

		log.Println("[启动] 写入默认配置到文件: config.toml")
		err = os.WriteFile("config.toml", configData, 0644)
		if err != nil {
			log.Printf("[启动] 写入默认配置文件失败: %v", err)
			return nil, fmt.Errorf("写入默认配置失败: %v", err)
		}

		log.Println("[启动] 默认配置创建完成")
		return defaultConfig, nil
	}

	log.Println("[启动] 成功读取配置文件: config.toml")
	log.Println("[启动] 配置文件内容:", string(data))

	config := &Config{}
	log.Println("[启动] 开始解析配置文件")
	// 直接解析 toml 到结构体
	err = toml.Unmarshal(data, config)
	if err != nil {
		log.Printf("[启动] 解析配置文件失败: %v", err)
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	log.Println("[启动] 配置文件解析完成")
	// 验证配置
	if config.Default.URL == "" {
		log.Println("[启动] URL配置为空，使用默认值 ./res/index.html")
		config.Default.URL = "./res/index.html"
	}

	log.Printf("[启动] 最终配置: %+v", config)
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
	log.Println("[桌面穿透] 开始设置壁纸")
	if wallpaperCancel != nil {
		log.Println("[桌面穿透] 取消上一个壁纸设置上下文")
		wallpaperCancel() // 先停止上一个goroutine
	}
	log.Println("[桌面穿透] 创建新的壁纸设置上下文")
	wallpaperCtx, wallpaperCancel = context.WithCancel(context.Background())

	// 使用增强的桌面设置功能
	log.Printf("[桌面穿透] 调用SetupAdvancedWallpaper(%s)", lorcaname)
	ret := SetupAdvancedWallpaper(lorcaname)
	log.Printf("[桌面穿透] SetupAdvancedWallpaper(%s) 返回: %v", lorcaname, ret)

	log.Println("[桌面穿透] 创建定时器")
	t = time.NewTicker(time.Second)
	log.Println("[桌面穿透] 启动定时任务")
	go func(ctx context.Context) {
		log.Println("[桌面穿透] 定时任务开始执行")
		failCount := 0
		for {
			select {
			case <-ctx.Done():
				log.Println("[桌面穿透] 定时任务上下文被取消")
				return
			case <-t.C:
				log.Println("[桌面穿透] 定时任务执行中")
				if hwnd := FindWindowByTitle(lorcaname); hwnd != 0 {
					log.Printf("[桌面穿透] 找到窗口句柄: %d", hwnd)
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
				} else {
					log.Println("[桌面穿透] 未找到窗口句柄")
				}
			}
		}
	}(wallpaperCtx)
	log.Println("[桌面穿透] 壁纸设置完成")
}

func runLorcaUI() {
	log.Printf("[Lorca] 开始创建UI: URL=%s, BrowserPath=%s", urlStr, BwPath)
	ui, err := lorca.New(urlStr, "", BwPath, 0, 0, "--kiosk", "--autoplay-policy=no-user-gesture-required")
	if err != nil {
		log.Printf("[Lorca] 创建UI失败: %v", err)
		// 失败后清理步骤
		if ui != nil {
			log.Println("[Lorca] 清理已创建的UI实例")
			_ = ui.Close()
		}
		windowMu.Lock()
		mainWindow = nil
		windowMu.Unlock()
		return
	}
	log.Println("[Lorca] UI创建成功")
	windowMu.Lock()
	mainWindow = ui
	windowMu.Unlock()
	lorcaname = "classpaper" + generateRandomString(6)
	log.Printf("[Lorca] 生成随机窗口标题: %s", lorcaname)
	ui.Eval("document.title='" + lorcaname + "'")
	log.Printf("[Lorca] 设置窗口标题: %s", lorcaname)

	// 绑定前端可调用的Go函数
	log.Println("[Lorca] 开始绑定前端可调用的Go函数")
	ui.Bind("getWidth", GetScreenWidth)
	ui.Bind("getHeight", GetScreenHeight)
	ui.Bind("readFile", func(path string) (string, error) {
		log.Printf("[Lorca] readFile 调用: %s", path)
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("[Lorca] readFile 失败: %v", err)
			return "", err
		}
		return string(data), err
	})
	ui.Bind("writeFile", func(path, content string) error {
		log.Printf("[Lorca] writeFile 调用: %s", path)
		return os.WriteFile(path, []byte(content), 0644)
	})
	ui.Bind("readDir", func(dir string) ([]string, error) {
		log.Printf("[Lorca] readDir 调用: %s", dir)
		entries, err := os.ReadDir(dir)
		if err != nil {
			log.Printf("[Lorca] readDir 失败: %v", err)
			return nil, err
		}
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		log.Printf("[Lorca] readDir 返回 %d 个项目", len(names))
		return names, nil
	})
	log.Println("[Lorca] 前端函数绑定完成")

	// 增加检测机制，确保窗口已创建并渲染后再设置壁纸
	log.Println("[Lorca] 等待窗口渲染完成")
	maxWait := 30 // 最多等待30次（约3秒）
	for i := 0; i < maxWait; i++ {
		// 检查窗口是否可用（可根据实际情况调整检测条件）
		res := ui.Eval("document.readyState")
		if res != nil && res.String() == "complete" {
			log.Printf("[Lorca] 窗口渲染完成，尝试次数: %d", i+1)
			break
		}
		time.Sleep(100 * time.Millisecond)
		if i == maxWait-1 {
			log.Println("[Lorca] 警告：等待窗口渲染超时，强制继续")
		}
	}

	log.Println("[Lorca] 开始设置壁纸")
	setWallpaper()
	log.Println("[Lorca] 壁纸设置完成，等待窗口关闭...")
	go func() {
		<-ui.Done()
		strictCloseMainWindow()
		log.Println("[Lorca] mainWindow 已关闭并清理")
	}()
	<-ui.Done()
	strictCloseMainWindow()
	log.Println("[Lorca] 窗口已关闭")
}

// ====== 托盘菜单相关函数区 ======
func reloadPage() {
	log.Println("[托盘] 触发网页重载")
	if mainWindow == nil {
		log.Println("[托盘] 警告: mainWindow 为 nil，无法重载网页")
		return
	}
	log.Println("[托盘] 执行网页重载")
	mainWindow.Eval("location.reload(true)")
	log.Println("[托盘] 网页重载完成")
}

func restartWebpageDisplayProgram() {
	log.Println("[托盘] 重启网页显示程序...")
	if mainWindow != nil {
		log.Println("[托盘] 关闭现有主窗口")
		mainWindow.Close()
	}
	log.Println("[托盘] 启动新的Lorca UI")
	go runLorcaUI()
	log.Println("[托盘] Lorca UI启动完成")
}

func onReady() {
	log.Println("[托盘] 开始初始化托盘菜单")
	// // 判断系统是否为夜间模式，选择不同的图标
	// var iconLight, iconDark []byte
	// iconLight = IconDataLight
	// iconDark = IconDataDark

	// // 检测夜间模式
	// isDarkMode := false
	// // Windows下可用注册表检测，macOS可用AppleScript，Linux需根据桌面环境
	// isDarkMode = GetWindowsDarkMode()

	// if isDarkMode {
	// 	systray.SetTemplateIcon(iconDark, iconDark)
	// } else {
	// 	systray.SetTemplateIcon(iconLight, iconLight)
	// }
	systray.SetTemplateIcon(IconData, IconData)

	systray.SetTitle(title)
	systray.SetTooltip(title)
	log.Println("[托盘] 创建菜单项")
	reloadMenuItem := systray.AddMenuItem("重载网页", "Reload Page")
	setPenetrationMenuItem := systray.AddMenuItem("设置程序桌面穿透", "Set Window Penetration")
	restartWebpageMenuItem := systray.AddMenuItem("重启网页显示程序", "Restart Webpage Display")
	settingsMenuItem := systray.AddMenuItem("设置", "Open Settings")
	restartMenuItem := systray.AddMenuItem("重启程序", "Restart Application")
	quitMenuItem := systray.AddMenuItem("退出程序", "Quit Application")
	log.Println("[托盘] 托盘菜单已初始化")
	log.Println("[托盘] 启动Lorca UI")
	go runLorcaUI()
	log.Println("[托盘] 启动托盘事件监听器")
	go func() {
		for {
			select {
			case <-reloadMenuItem.ClickedCh:
				log.Println("[托盘] 手动触发网页重载")
				reloadPage()
				log.Println("[托盘] 网页重载处理完成")
			case <-setPenetrationMenuItem.ClickedCh:
				log.Println("[托盘] 手动触发桌面穿透")
				setWallpaper()
				log.Println("[托盘] 桌面穿透设置完成")
			case <-restartWebpageMenuItem.ClickedCh:
				log.Println("[托盘] 手动重启网页显示程序")
				restartWebpageDisplayProgram()
				log.Println("[托盘] 网页显示程序重启完成")
			case <-settingsMenuItem.ClickedCh:
				log.Println("[托盘] 打开设置窗口")
				openSettings()
				log.Println("[托盘] 设置窗口打开完成")
			case <-restartMenuItem.ClickedCh:
				log.Println("[托盘] 重启主程序...")
				restartProgram()
				systray.Quit()
				log.Println("[托盘] 程序重启完成")
			case <-quitMenuItem.ClickedCh:
				log.Println("[托盘] 退出程序")
				systray.Quit()
				log.Println("[托盘] 程序退出完成")
				return
			}
		}
	}()
	log.Println("[托盘] 托盘初始化完成")
}

func openSettings() {
	log.Println("[设置] 开始打开设置窗口")
	windowMu.Lock()
	if settingsWindow != nil {
		windowMu.Unlock()
		log.Println("[设置] 设置窗口已经打开")
		return
	}
	windowMu.Unlock()

	log.Println("[设置] 创建新的设置窗口上下文")
	if settingsCancel != nil {
		log.Println("[设置] 取消上一个设置窗口的上下文")
		settingsCancel() // 关闭上一个context
	}
	settingsCtx, settingsCancel = context.WithCancel(context.Background())

	log.Println("[设置] 获取设置页面URL")
	settingsURL, err := getFilePathURL("res/settings.html")
	if err != nil {
		log.Printf("[设置] 获取设置页面URL失败: %v", err)
		return
	}

	// 获取屏幕尺寸
	log.Println("[设置] 获取屏幕尺寸")
	screenWidth := GetScreenWidth()
	screenHeight := GetScreenHeight()
	log.Printf("[设置] 屏幕尺寸: %dx%d", screenWidth, screenHeight) // 新增日志

	// 设置窗口尺寸（调整为屏幕宽度的50%，高度的70%）
	log.Println("[设置] 计算设置窗口尺寸")
	width := int(float64(screenWidth) * 0.5)   // 原0.6改为0.5
	height := int(float64(screenHeight) * 0.7) // 原0.8改为0.7

	// 确保窗口尺寸在合理范围内
	if width < 1000 {
		width = 1000 // 最小宽度
	}
	if height < 800 {
		height = 800 // 最小高度
	}

	// 计算窗口位置（居中）
	log.Println("[设置] 计算设置窗口位置")
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2

	// 创建窗口并设置位置
	log.Printf("[设置] 创建设置窗口: URL=%s, 宽度=%d, 高度=%d", settingsURL, width, height)
	ui, err := lorca.New(settingsURL, "", "", width, height)
	if err != nil {
		log.Printf("[设置] 创建设置窗口失败: %v", err)
		return
	}
	log.Println("[设置] 设置窗口创建成功")

	// 设置窗口位置
	log.Printf("[设置] 设置窗口位置: x=%d, y=%d", x, y)
	ui.Eval(fmt.Sprintf(`
		window.moveTo(%d, %d);
		document.title = "ClassPaper 设置";
	`, x, y))

	windowMu.Lock()
	settingsWindow = ui
	windowMu.Unlock()

	// 立即绑定函数
	log.Println("[设置] 开始绑定函数...")

	// 绑定配置相关函数
	log.Println("[设置] 绑定readConfig函数")
	err = ui.Bind("readConfig", func() (interface{}, error) {
		log.Println("[设置] readConfig被调用")
		config, err := ParseConfig()
		if err != nil {
			log.Printf("[设置] 读取配置失败: %v", err)
			return nil, err
		}
		log.Printf("[设置] 读取到配置: %+v", config)
		log.Println("[设置] readConfig执行完成")
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

	log.Println("[设置] 绑定saveConfig函数")
	err = ui.Bind("saveConfig", func(configJSON string) error {
		log.Printf("[设置] saveConfig被调用，配置JSON: %s", configJSON)
		var config Config
		err := json.Unmarshal([]byte(configJSON), &config)
		if err != nil {
			log.Printf("[设置] 解析配置JSON失败: %v", err)
			return fmt.Errorf("解析配置JSON失败: %v", err)
		}

		// 将配置写入文件
		log.Println("[设置] 序列化配置数据")
		configData, err := toml.Marshal(config)
		if err != nil {
			log.Printf("[设置] 序列化配置失败: %v", err)
			return fmt.Errorf("序列化配置失败: %v", err)
		}

		log.Println("[设置] 写入配置文件: config.toml")
		err = os.WriteFile("config.toml", configData, 0644)
		if err != nil {
			log.Printf("[设置] 写入配置文件失败: %v", err)
			return fmt.Errorf("写入配置文件失败: %v", err)
		}

		log.Printf("[设置] 配置已保存: %+v", config)

		// 更新当前运行的配置
		log.Println("[设置] 更新当前运行配置")
		urlStr = NormalizeURL(config.Default.URL)
		BwPath = config.Default.BrowserPath

		log.Println("[设置] saveConfig执行完成")
		return nil
	})
	if err != nil {
		log.Printf("[设置] 绑定saveConfig失败: %v", err)
	}

	// 绑定文件写入函数
	log.Println("[设置] 绑定writeFile函数")
	err = ui.Bind("writeFile", func(path string, content string) error {
		log.Printf("[设置] writeFile被调用: %s", path)
		return os.WriteFile(path, []byte(content), 0644)
	})
	if err != nil {
		log.Printf("[设置] 绑定writeFile失败: %v", err)
	}

	// 绑定壁纸扫描函数
	log.Println("[设置] 绑定scanWallpaperDir函数")
	err = ui.Bind("scanWallpaperDir", scanWallpaperDir)
	if err != nil {
		log.Printf("[设置] 绑定scanWallpaperDir失败: %v", err)
	}

	// 绑定主窗口刷新函数
	log.Println("[设置] 绑定reloadMainWindow函数")
	err = ui.Bind("reloadMainWindow", func() error {
		log.Println("[设置] reloadMainWindow被调用")
		if mainWindow != nil {
			log.Println("[设置] 刷新主窗口")
			mainWindow.Eval("location.reload(true)")
			log.Println("[设置] 主窗口刷新完成")
			return nil
		}
		log.Println("[设置] 主窗口未打开，无法刷新")
		return fmt.Errorf("主窗口未打开")
	})
	if err != nil {
		log.Printf("[设置] 绑定reloadMainWindow失败: %v", err)
	}

	// 新增浏览器打开绑定
	log.Println("[设置] 绑定openURLInBrowser函数")
	err = ui.Bind("openURLInBrowser", func(url string) bool {
		log.Printf("[设置] openURLInBrowser被调用: %s", url)
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", url)
		case "darwin":
			cmd = exec.Command("open", url)
		default: // linux
			cmd = exec.Command("xdg-open", url)
		}

		if err := cmd.Start(); err != nil {
			log.Printf("[调试] 打开浏览器失败: %v", err)
			return false
		}
		log.Println("[设置] 浏览器打开完成")
		return true
	})
	if err != nil {
		log.Printf("[设置] 绑定openURLInBrowser失败: %v", err)
	}

	log.Println("[设置] 函数绑定完成")

	// 监听窗口关闭
	log.Println("[设置] 启动设置窗口关闭监听器")
	go func(ctx context.Context) {
		select {
		case <-ui.Done():
			strictCloseSettingsWindow()
			log.Println("[设置] settingsWindow 已关闭并清理")
		case <-ctx.Done():
			log.Println("[设置] settingsWindow context 被取消，资源清理")
		}
	}(settingsCtx)
	log.Println("[设置] 设置窗口打开完成")
}

func restartProgram() {
	log.Println("[重启] 开始重启程序")
	execPath, err := os.Executable()
	if err != nil {
		log.Printf("[重启] 获取可执行文件路径失败: %v", err)
		return
	}
	log.Printf("[重启] 当前可执行文件路径: %s", execPath)
	execDir := filepath.Dir(execPath)
	log.Printf("[重启] 可执行文件目录: %s", execDir)
	cmd := exec.Command(execPath)
	cmd.Dir = execDir
	log.Println("[重启] 启动新进程")
	err = cmd.Start()
	if err != nil {
		log.Printf("[重启] 启动新进程失败: %v", err)
		return
	}
	log.Println("[重启] 新进程启动成功")
}

func strictCloseMainWindow() {
	log.Println("[关闭] 开始关闭主窗口")
	windowMu.Lock()
	defer windowMu.Unlock()
	if mainWindow != nil {
		log.Println("[关闭] 调用mainWindow.Close()")
		err := mainWindow.Close()
		if err != nil {
			log.Printf("[关闭] mainWindow.Close() 失败: %v", err)
		} else {
			log.Println("[关闭] mainWindow.Close() 成功")
		}
		mainWindow = nil
	} else {
		log.Println("[关闭] mainWindow 为 nil，无需关闭")
	}
	log.Println("[关闭] 主窗口关闭完成")
}

func strictCloseSettingsWindow() {
	log.Println("[关闭] 开始关闭设置窗口")
	windowMu.Lock()
	defer windowMu.Unlock()
	if settingsWindow != nil {
		log.Println("[关闭] 调用settingsWindow.Close()")
		err := settingsWindow.Close()
		if err != nil {
			log.Printf("[关闭] settingsWindow.Close() 失败: %v", err)
		} else {
			log.Println("[关闭] settingsWindow.Close() 成功")
		}
		settingsWindow = nil
	} else {
		log.Println("[关闭] settingsWindow 为 nil，无需关闭")
	}
	if settingsCancel != nil {
		log.Println("[关闭] 取消设置窗口上下文")
		settingsCancel()
		settingsCancel = nil
	} else {
		log.Println("[关闭] settingsCancel 为 nil，无需取消")
	}
	// TODO: 这里可扩展更多设置窗口相关资源的清理（如异步任务、临时文件等）
	log.Println("[关闭] 设置窗口关闭完成")
}

func endup() {
	log.Println("[退出] 执行endup，关闭窗口和资源...")
	strictCloseMainWindow()
	strictCloseSettingsWindow()
	if wallpaperCancel != nil {
		log.Println("[退出] 取消壁纸设置上下文")
		wallpaperCancel()
		wallpaperCancel = nil
	} else {
		log.Println("[退出] wallpaperCancel 为 nil，无需取消")
	}
	if t != nil {
		log.Println("[退出] 停止定时器")
		t.Stop()
		t = nil
	} else {
		log.Println("[退出] 定时器为 nil，无需停止")
	}
	log.Println("[退出] 退出系统托盘")
	systray.Quit()
	if logFile != nil {
		log.Println("[退出] 同步并关闭日志文件")
		logFile.Sync()
		logFile.Close()
		logFile = nil
	} else {
		log.Println("[退出] 日志文件为 nil，无需关闭")
	}
	log.Println("[退出] endup执行完成")
}

func onExit() {
	log.Println("[退出] 程序退出，清理资源...")
	strictCloseMainWindow()
	strictCloseSettingsWindow()
	if wallpaperCancel != nil {
		log.Println("[退出] 取消壁纸设置上下文")
		wallpaperCancel()
		wallpaperCancel = nil
	} else {
		log.Println("[退出] wallpaperCancel 为 nil，无需取消")
	}
	if t != nil {
		log.Println("[退出] 停止定时器")
		t.Stop()
		t = nil
	} else {
		log.Println("[退出] 定时器为 nil，无需停止")
	}
	if logFile != nil {
		log.Println("[退出] 同步并关闭日志文件")
		logFile.Sync()
		logFile.Close()
		logFile = nil
	} else {
		log.Println("[退出] 日志文件为 nil，无需关闭")
	}
	log.Println("[退出] 程序退出处理完成")
}

// ====== 主程序入口和初始化 ======
func main() {
	log.Println("[启动] 程序开始启动")
	defer func() {
		log.Println("[启动] 捕获到panic或正常退出，执行defer函数")
		if r := recover(); r != nil {
			log.Printf("[Panic] %v", r)
		}
		log.Println("[启动] 调用endup进行资源清理")
		endup()
		log.Println("[启动] 调用lorca.CloseAllUIs释放所有UI资源")
		lorca.CloseAllUIs() // 兜底释放所有UI资源
		log.Println("[启动] 程序退出")
	}()
	var err error
	log.Println("[启动] 创建日志文件: app.log")
	logFile, err = os.Create("app.log")
	if err != nil {
		log.Fatalf("[启动] 创建日志文件失败: %v", err)
	}
	log.Println("[启动] 设置日志输出到文件和标准输出")
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Println("[启动] 开始读取配置文件")
	config, err := ParseConfig()
	if err != nil {
		log.Printf("[启动] 读取配置失败: %v", err)
		return
	}
	log.Printf("[启动] 加载配置URL: %s", config.Default.URL)
	log.Println("[启动] 开始标准化URL")
	urlStr = NormalizeURL(config.Default.URL)
	BwPath = config.Default.BrowserPath
	log.Printf("[启动] 标准化URL: %s", urlStr)
	log.Printf("[启动] 浏览器路径: %s", BwPath)
	log.Println("[启动] 启动系统托盘")
	systray.Run(onReady, onExit)
	log.Println("[启动] systray.Run执行完成")
}

func init() {
	log.Println("[启动] 执行init函数")
	if runtime.GOOS == "windows" {
		log.Println("[启动] 设置Windows环境变量WINGUI_NO_CONSOLE=1")
		os.Setenv("WINGUI_NO_CONSOLE", "1")
		log.Println("[启动] 设置Windows控制台代码页为UTF-8")
		exec.Command("cmd", "/c", "chcp", "65001").Run()

		// 处理DPI感知结果
		log.Println("[启动] 设置Windows DPI感知")
		if !SetDPIAware() {
			log.Printf("[启动] DPI感知设置失败")
		} else {
			log.Println("[启动] DPI感知设置成功")
		}
	}
	log.Println("[启动] init函数执行完成")
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
