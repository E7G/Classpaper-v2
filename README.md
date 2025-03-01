# Classpaper-v2
 version 2 of my work Classpaper


---

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

---

# 设置窗口程序集

---

### 简介

设置窗口程序集提供了一个简单的用户界面，用于编辑和保存课程表、显示内容、事件以及配置程序的 URL 和浏览器路径。用户可以通过此界面轻松地管理和调整程序的各项设置。

### 功能

- 编辑和保存课程表、显示内容、事件。
- 配置程序的 URL 和浏览器路径。
- 提供保存成功或失败的信息提示。

### 使用方法

1. 打开设置窗口程序集。
2. 在相应的输入框中编辑内容。
3. 点击对应的保存按钮以保存修改。
4. 如果保存成功，将会显示保存成功的提示；否则将显示保存失败的提示。

### 文件结构

- `res/config/lessons.js`: 课程表内容的存储文件。
- `res/config/sth.js`: 显示内容的存储文件。
- `res/config/events.js`: 事件内容的存储文件。
- `config.ini`: 程序的配置文件，包含程序的 URL 和浏览器路径。

### 注意事项

- 在保存内容之前，请确保输入的内容符合相应的格式要求。
- 请注意保存成功或失败的提示信息，以便及时了解保存操作的结果。

---

### 许可证

本项目使用 LGPL 许可证。有关详细信息，请参阅 [LICENSE](LICENSE) 文件。

### 联系方式

如有任何问题或反馈意见，请随时提出问题或联系 [e7g](https://github.com/e7g/)。
