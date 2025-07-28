# 前端代码重构总结

## 重构目标
1. 将所有可配置的元素提取到 `config.js` 中
2. 创建模块化的配置管理系统
3. 预留 LLM 集成接口用于智能调整界面
4. 提高代码的可维护性和可扩展性

## 主要变更

### 1. 配置文件重构 (`res/config/config.js`)

#### 新增配置结构：
- **app**: 应用基础配置（标题、语言、主题、图标）
- **ui**: 界面布局和组件配置
- **styles**: 样式配置（字体、颜色、效果、间距）
- **lessons**: 课程表配置
- **time**: 时间相关配置
- **notifications**: 通知配置
- **events**: 事件日历配置
- **wallpapers**: 壁纸配置
- **announcement**: 告示牌内容
- **llm**: LLM集成配置

#### 新增工具类：
- **LLMIntegration**: LLM集成工具类
- **ConfigManager**: 配置管理工具类
- **ThemeManager**: 主题管理工具类

### 2. UI管理器 (`res/js/ui-manager.js`)

创建了统一的UI管理器，负责：
- 根据配置动态初始化所有组件
- 处理配置变更事件
- 管理组件的显示/隐藏状态
- 设置LLM集成按钮和事件

### 3. 样式系统重构 (`res/css/custom.css`)

- 引入CSS自定义属性（CSS Variables）
- 所有颜色、字体、间距等样式值都通过配置驱动
- 支持动态主题切换

### 4. HTML结构优化 (`res/index_new.html`)

- 添加了更多的ID和类名用于配置驱动
- 组件结构更加模块化
- 预留了LLM集成界面

### 5. JavaScript代码更新

#### `res/js/main.js`:
- 更新壁纸管理逻辑使用新配置结构
- 时间显示格式可配置
- 支持主题动态应用

## LLM集成接口

### 功能特性
1. **智能布局调整**: 根据当前时间和课程安排调整界面布局
2. **内容生成**: 生成适合当前时间段的激励性文案
3. **课程优化**: 分析课程安排并提供学习建议
4. **语音控制**: 预留语音控制接口

### API接口
```javascript
// 调用LLM进行布局调整
await window.LLMIntegration.adjustLayout();

// 生成内容
await window.LLMIntegration.generateContent();

// 优化课程安排
await window.LLMIntegration.optimizeSchedule();
```

### 配置示例
```javascript
"llm": {
  "enabled": true,
  "apiEndpoint": "https://api.openai.com/v1/chat/completions",
  "apiKey": "your-api-key",
  "model": "gpt-3.5-turbo",
  "features": {
    "autoAdjustLayout": true,
    "smartScheduling": true,
    "contentGeneration": true,
    "voiceControl": false
  }
}
```

## 配置管理

### 动态配置更新
```javascript
// 更新配置
ConfigManager.updateConfig('app.theme', 'light');

// 获取配置
const theme = ConfigManager.getConfig('app.theme');

// 导出配置
const configJson = ConfigManager.exportConfig();

// 导入配置
ConfigManager.importConfig(configJson);
```

### 事件监听
```javascript
// 监听配置变更
window.addEventListener('configChanged', (event) => {
  console.log('配置已更新:', event.detail);
});
```

## 主题系统

### 动态主题切换
```javascript
// 切换主题
ThemeManager.switchTheme('light');

// 应用主题
ThemeManager.applyTheme();
```

### CSS变量支持
所有样式都通过CSS变量控制，支持实时更新：
- `--primary-color`
- `--secondary-color`
- `--accent-color`
- `--backdrop-blur`
- `--border-radius`
- 等等...

## 向后兼容性

保持了原有的变量导出，确保现有代码正常工作：
- `lessons`
- `events`
- `wallpaperlist`
- `sth`

## 使用方法

1. **基础配置**: 修改 `CONFIG` 对象中的相应字段
2. **LLM集成**: 设置 `llm.enabled: true` 并配置API信息
3. **主题定制**: 修改 `styles` 配置或使用 `ThemeManager`
4. **组件控制**: 通过 `ui.components` 配置启用/禁用组件

## 扩展性

- 易于添加新的配置项
- 支持插件式组件开发
- LLM集成为未来AI功能预留了接口
- 模块化设计便于维护和扩展

这次重构大大提高了代码的可配置性和可维护性，同时为未来的AI集成功能奠定了基础。