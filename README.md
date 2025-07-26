# ClassPaper-v2

ClassPaper-v2 是一个用 Go 语言开发的高级桌面壁纸程序，可以将任意网页作为动态壁纸显示在桌面。程序采用先进的 Windows API 技术，提供专业级的桌面集成体验，支持 Windows 10/11 最新版本（包括 24H2）。

---

## 核心功能特性

### 🖥️ 高级桌面集成
- **智能桌面穿透**：采用 DWM (Desktop Window Manager) 技术实现真正的桌面嵌入
- **多版本兼容**：支持 Windows 10/11 各版本，包括最新的 24H2 更新
- **Z序监控系统**：自动维护窗口层次，防止被其他程序遮挡
- **透明效果**：使用 DWM 扩展框架实现专业级透明效果

### 🌐 网页壁纸功能
- 将任意网页作为动态桌面壁纸显示
- 支持本地 HTML 文件和在线网页
- 自动全屏适配，完美融入桌面环境
- 支持交互式网页内容

### 🎛️ 系统托盘控制
- **网页重载**：实时刷新壁纸内容
- **桌面穿透设置**：手动触发高级桌面集成
- **程序管理**：重启网页显示、重启程序、安全退出
- **设置界面**：图形化配置管理

### 🔧 高级技术特性
- **进程安全检测**：识别合法的 Explorer 工作窗口
- **窗口样式优化**：自动配置最佳的窗口属性
- **资源管理**：严格的窗口生命周期管理，防止内存泄漏
- **错误恢复**：自动检测和修复常见的显示问题

---

## 前端设置与功能优化（2025年7月重点更新）

- **设置界面美化**：整体风格更简洁，tab切换流畅，课程表tab与壁纸tab均有美化，去除斑马纹，保留圆角和悬停高亮。
- **进度条描述与百分比显示**：进度条下方描述支持自定义，百分比可切换“剩余/已过”两种方式。
- **重置按钮弹窗**：重置按钮增加自定义对话框确认弹窗，防止误操作。
- **智能推算课程时间**：支持一键自动补全后续节次时间，兼容input type=time格式（如08:50），无需手动填写全部节次。
- **课表高亮优化**：动态高亮当前课程，避免死板高亮第6格。
- **状态提示优化**：每次调用showStatus都能弹出提示，消息显示更及时。

---

## 壁纸管理与美化

- **壁纸tab结构优化**：壁纸列表支持缩略图、卡片分组，按钮美化，切换间隔可自定义。
- **壁纸切换体验提升**：切换时平滑过渡，支持定时自动切换。

---

## 🔧 技术架构与实现

### Windows API 集成
- **DWM 集成**：使用 `DwmExtendFrameIntoClientArea` 和 `DwmEnableBlurBehindWindow` 实现专业透明效果
- **窗口管理**：通过 `SetWindowLongPtr` 和 `SetWindowPos` 精确控制窗口属性和层次
- **桌面检测**：智能识别 `Progman` 和 `SHELLDLL_DefView` 窗口，适配不同 Windows 版本
- **进程验证**：使用 PSAPI 验证 Explorer 进程，确保系统安全

### 兼容性设计
- **双模式架构**：
  - **高级模式**（Windows 10/11）：DWM 透明效果 + Z序监控 + 多版本适配
  - **传统模式**（Windows 7/8/8.1）：经典 Progman/WorkerW 方案
- **自动版本检测**：运行时检测 Windows 版本并自动选择最佳方案
- **智能回退**：高级模式失败时自动回退到传统模式
- **错误处理**：全面的错误检测和恢复机制，确保程序稳定运行
- **高DPI支持**：自动适配高分辨率显示器，图标和界面完美显示

### 性能优化
- **异步监控**：Z序监控使用独立 goroutine，不影响主程序性能
- **资源管理**：严格的窗口生命周期管理，防止内存泄漏
- **智能检测**：仅在必要时进行窗口调整，减少系统资源消耗

---

## 🐛 故障排除

### 桌面穿透问题
- **症状**：壁纸显示在桌面图标上方
- **解决**：点击托盘菜单中的"设置程序桌面穿透"重新设置
- **原因**：系统更新或其他程序可能影响窗口层次

### 程序无法启动
- **系统检查**：支持 Windows 7 及以上版本，推荐 Windows 10/11
- **权限**：以管理员身份运行可能解决某些权限问题
- **日志**：查看 `app.log` 文件获取详细错误信息
- **兼容性**：老系统会自动使用传统模式，功能略有限制

### 网页显示异常
- **调试**：使用设置界面的"在浏览器打开调试"功能
- **兼容性**：确保网页内容兼容 Chromium 内核
- **路径**：检查配置文件中的 URL 路径是否正确

---

## 🚀 快速开始

### 系统要求
- **操作系统**：
  - **最佳体验**：Windows 10 1903+ 或 Windows 11（支持所有高级功能）
  - **基础支持**：Windows 7/8/8.1（传统桌面穿透模式）
- **开发环境**：Go 1.19+ （如需从源码编译）
- **硬件要求**：
  - **Windows 10/11**：支持 DWM 的显卡，4GB+ 内存
  - **Windows 7/8/8.1**：2GB+ 内存即可

### 安装使用

#### 方式一：直接运行（推荐）
1. 下载最新的 `classpaper.exe` 可执行文件
2. 创建或编辑 `config.toml` 配置文件
3. 双击运行 `classpaper.exe`
4. 程序将在系统托盘显示图标

#### 方式二：从源码编译
1. 克隆本仓库：
   ```bash
   git clone https://github.com/your-repo/classpaper.git
   cd classpaper
   ```
2. 安装依赖：
   ```bash
   go mod tidy
   ```
3. 编译运行：
   ```bash
   go build .
   ./classpaper.exe
   ```

### 首次配置
1. 程序启动后，右键点击托盘图标
2. 选择"设置"打开配置界面
3. 配置网页 URL 和其他选项
4. 点击"设置程序桌面穿透"激活壁纸模式

---

## 配置文件说明

配置文件为 `config.toml`，支持两种写法：

**推荐写法：**

```toml
url = "http://example.com"         # 要显示的网页URL
browser_path = ""                  # 浏览器可执行文件路径（可选，留空用默认浏览器）
```

**或带区块写法：**

```toml
[default]
url = "http://example.com"
browser_path = ""
```

---

## 🎛️ 托盘菜单功能

| 菜单项 | 功能说明 | 使用场景 |
|--------|----------|----------|
| **重载网页** | 强制刷新壁纸网页内容 | 网页内容更新后需要刷新显示 |
| **设置程序桌面穿透** | 重新激活高级桌面集成 | 系统更新后壁纸被遮挡时使用 |
| **重启网页显示程序** | 重新创建壁纸窗口 | 网页显示异常或需要重新加载时 |
| **设置** | 打开图形化配置界面 | 修改 URL、浏览器路径等配置 |
| **重启程序** | 完全重启 ClassPaper | 配置更改后需要重启生效时 |
| **退出程序** | 安全退出并清理所有资源 | 停止使用程序时 |

---

## 📁 项目结构

```
classpaper/
├── main.go                 # 主程序逻辑和托盘管理
├── winapi.go              # Windows API 集成和桌面穿透
├── resources.go           # 嵌入式资源文件
├── config.toml            # 主配置文件
├── app.log               # 运行日志文件
├── res/                  # 静态资源目录
│   ├── index.html        # 默认壁纸页面
│   ├── settings.html     # 设置界面
│   ├── css/             # 样式文件
│   ├── js/              # JavaScript 文件
│   └── wallpaper/       # 壁纸图片目录
├── icon/                # 程序图标
│   ├── light/           # 亮色主题图标
│   └── dark/            # 暗色主题图标
└── lib/                 # 第三方库
    └── lorca/           # 修改版 Lorca 库
```

---

## 🔧 技术依赖

### 核心依赖
- **Go Runtime**: Go 1.19+ （编译时需要）
- **Windows API**: User32, DWMApi, Kernel32, PSAPI
- **Chromium**: 通过 Lorca 提供的嵌入式浏览器

### Go 模块依赖
```go
require (
    github.com/getlantern/systray v1.2.2  // 系统托盘
    github.com/zserge/lorca v0.1.10       // 嵌入式浏览器
    github.com/pelletier/go-toml v1.9.5   // TOML 配置解析
)
```

### 前端技术栈
- **HTML5/CSS3/ES6**: 现代 Web 标准
- **Pico.css**: 轻量级 CSS 框架
- **原生 JavaScript**: 无额外框架依赖

---

## ⚠️ 重要说明

### 系统兼容性
- **完全支持**: Windows 10 1903+ / Windows 11（高级模式 + 所有功能）
- **兼容支持**: Windows 7/8/8.1（传统模式，基础桌面穿透）
- **推荐配置**: Windows 11 22H2+ 获得最佳体验
- **自动检测**: 程序会自动检测系统版本并选择最佳兼容方案

### 安全考虑
- 程序需要访问桌面窗口层次，某些安全软件可能报警
- 建议将程序添加到安全软件白名单
- 程序不会修改系统文件或注册表

### 性能影响
- 内存占用: 通常 50-100MB（取决于网页内容）
- CPU 占用: 空闲时 < 1%，网页动画时可能增加
- 对系统桌面性能影响极小

---

## 🎯 开发计划

### 已完成功能 ✅
- [x] 高级桌面穿透技术（基于 DWM API）
- [x] Windows 10/11 多版本兼容
- [x] Z序监控和自动修复
- [x] 透明效果和窗口样式优化
- [x] 系统托盘完整功能
- [x] 图形化设置界面
- [x] 资源管理和内存泄漏防护

### 计划中功能 🚧
- [ ] 多显示器支持
- [ ] 壁纸切换动画效果
- [ ] 性能监控面板
- [ ] 插件系统支持
- [ ] 远程配置管理
- [ ] 主题切换功能

### 长期目标 🎯
- [ ] Linux 桌面环境支持（X11/Wayland）
- [ ] macOS 桌面集成
- [ ] 云端壁纸库
- [ ] 社区分享平台

---

## 🤝 贡献指南

### 报告问题
1. 查看 [Issues](https://github.com/your-repo/classpaper/issues) 确认问题未被报告
2. 提供详细的系统信息和错误日志
3. 描述重现步骤和期望行为

### 代码贡献
1. Fork 本仓库
2. 创建功能分支: `git checkout -b feature/your-feature`
3. 提交更改: `git commit -am 'Add some feature'`
4. 推送分支: `git push origin feature/your-feature`
5. 提交 Pull Request

### 开发环境设置
```bash
# 克隆仓库
git clone https://github.com/your-repo/classpaper.git
cd classpaper

# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 编译调试版本
go build -tags debug .
```

---

## 📄 许可证

本项目使用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

```
MIT License

Copyright (c) 2024 ClassPaper Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

---

## 📞 联系方式

- **项目维护者**: [e7g](https://github.com/e7g/)
- **问题反馈**: [GitHub Issues](https://github.com/your-repo/classpaper/issues)
- **功能建议**: [GitHub Discussions](https://github.com/your-repo/classpaper/discussions)

---

## 🙏 致谢

感谢以下开源项目和贡献者：

- [Lorca](https://github.com/zserge/lorca) - 提供嵌入式浏览器支持
- [Systray](https://github.com/getlantern/systray) - 跨平台系统托盘
- [Go-TOML](https://github.com/pelletier/go-toml) - TOML 配置解析
- [Pico.css](https://picocss.com/) - 轻量级 CSS 框架
- [实现桌面动态壁纸（一）\_动态壁纸原理-CSDN博客](https://blog.csdn.net/qq_59075481/article/details/125361650) - 提供桌面壁纸实现原理

特别感谢所有测试用户和贡献者的反馈和建议！

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给我们一个 Star！**

[🏠 返回顶部](#classpaper)

</div>
