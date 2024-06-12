# Classpaper-v2
 version 2 of my work Classpaper

```markdown
# ClassPaper

---

### 简介

ClassPaper 是一个使用 Go 语言编写的程序，旨在提供无缝的桌面壁纸体验，同时在后台运行网页内容。它允许您将任何网页设置为桌面壁纸，提供桌面穿透、自动重载以及更多功能。

### 功能特性

- 将任何网页设置为桌面壁纸。
- 启用桌面穿透，以便与底层桌面元素进行平滑交互。
- 根据需求自动重新加载网页。
- 轻松重启程序或网页显示程序。
- 通过 `config.ini` 文件进行自定义设置。

### 使用方法

1. 将此存储库克隆到本地计算机。
2. 确保已安装 Go。
3. 根据需要编辑 `config.ini` 文件以指定 URL 和浏览器路径。
4. 使用命令 `go run main.go` 运行程序。
5. 在系统托盘图标上右键单击，以访问重新加载网页、设置桌面穿透、重新启动程序和退出应用程序等选项。

### 配置

您可以通过编辑 `config.ini` 文件来自定义 ClassPaper 的行为：

```ini
[url]
url = "http://example.com"  # 要显示的网页的 URL
browser_path = ""           # 浏览器可执行文件的路径（可选，留空使用默认浏览器）
```

### 依赖项

- `github.com/getlantern/systray`：用于创建系统托盘图标。
- `github.com/zserge/lorca`：用于在应用程序中嵌入浏览器窗口。
- Go 1.16 或更高版本。

### 许可证

本项目使用 LGPL 许可证。有关详细信息，请参阅 [LICENSE](LICENSE) 文件。

### 联系方式

如有任何问题或反馈意见，请随时提出问题或联系 [e7g](https://github.com/e7g/)。
```
