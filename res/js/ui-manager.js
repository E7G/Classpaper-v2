// UI管理器 - 根据配置动态初始化界面组件
class UIManager {
  constructor() {
    this.config = CONFIG;
    this.initialized = false;
  }

  // 初始化所有UI组件
  init() {
    if (this.initialized) return;
    
    this.initializeApp();
    this.initializeComponents();
    this.initializeAudio();
    this.initializeLLMInterface();
    this.bindEvents();
    
    this.initialized = true;
    console.log('UI Manager initialized');
  }

  // 初始化应用基础设置
  initializeApp() {
    const app = this.config.app;
    
    // 设置标题
    document.title = app.title;
    const titleElement = document.getElementById('app-title');
    if (titleElement) titleElement.textContent = app.title;
    
    // 设置语言
    document.documentElement.lang = app.language;
    
    // 设置主题
    document.documentElement.setAttribute('data-theme', app.theme);
    
    // 设置图标
    const favicon = document.getElementById('app-favicon');
    if (favicon) favicon.href = app.favicon;
  }

  // 初始化组件
  initializeComponents() {
    this.initializeClassTable();
    this.initializeProgressBar();
    this.initializeClock();
    this.initializeCalendar();
    this.initializeHelp();
    this.initializeRemember();
  }

  // 初始化课程表
  initializeClassTable() {
    const component = this.config.ui.components.classTable;
    if (!component.enabled) return;

    const table = document.getElementById(component.id);
    if (!table) return;

    // 清空现有内容
    table.innerHTML = '';

    // 生成12个课程格子（c0到c11）
    for (let i = 0; i <= 11; i++) {
      const div = document.createElement('div');
      div.id = `c${i}`;
      
      if (i === 0) {
        // 表头
        div.innerHTML = '<a href="#" role="button" class="contrast table-header">课程</a>';
      } else {
        // 课程格子
        div.innerHTML = `<a href="#" role="button" class="contrast" id="c_b${i}" style="opacity: 0.5;">加载中...</a>`;
      }
      
      table.appendChild(div);
    }

    // 应用样式
    table.style.backdropFilter = `blur(${component.blurEffect}px)`;
  }

  // 初始化进度条
  initializeProgressBar() {
    const component = this.config.ui.components.progressBar;
    if (!component.enabled) return;

    const progressBar = document.getElementById(component.id);
    if (!progressBar) return;

    if (component.hidden) {
      progressBar.hidden = true;
    }
  }

  // 初始化时钟组件
  initializeClock() {
    const component = this.config.ui.components.clock;
    if (!component.enabled) {
      const clockElement = document.getElementById(component.id);
      if (clockElement) clockElement.style.display = 'none';
      return;
    }

    const header = document.getElementById('clock-header');
    const headerText = document.getElementById('clock-header-text');
    
    if (header) {
      header.hidden = !component.showHeader;
    }
    
    if (headerText) {
      headerText.textContent = component.headerText;
    }
  }

  // 初始化日历组件
  initializeCalendar() {
    const component = this.config.ui.components.calendar;
    if (!component.enabled) {
      const calendarElement = document.getElementById(component.id);
      if (calendarElement) calendarElement.style.display = 'none';
      return;
    }

    const header = document.getElementById('calendar-header');
    const headerText = document.getElementById('calendar-header-text');
    
    if (header) {
      header.hidden = !component.showHeader;
    }
    
    if (headerText) {
      headerText.textContent = component.headerText;
    }
  }

  // 初始化帮助组件
  initializeHelp() {
    const component = this.config.ui.components.help;
    if (!component.enabled) {
      const helpElement = document.getElementById(component.id);
      if (helpElement) helpElement.style.display = 'none';
      return;
    }

    const header = document.getElementById('help-header');
    const headerText = document.getElementById('help-header-text');
    const helpContent = document.getElementById('help');
    
    if (header) {
      header.hidden = !component.showHeader;
    }
    
    if (headerText) {
      headerText.textContent = component.headerText;
    }

    // 应用列布局样式
    if (helpContent) {
      helpContent.style.columnCount = component.columnCount;
      helpContent.style.columnGap = component.columnGap;
    }
  }


  // 初始化音频
  initializeAudio() {
    const notifications = this.config.notifications;
    if (!notifications.enabled) return;

    const regularAudio = document.getElementById('regularNotification');
    const endingAudio = document.getElementById('endingNotification');

    if (regularAudio) {
      regularAudio.src = notifications.sounds.regular;
    }

    if (endingAudio) {
      endingAudio.src = notifications.sounds.ending;
    }
  }

  // 初始化LLM接口
  initializeLLMInterface() {
    const llmConfig = this.config.llm;
    const llmInterface = document.getElementById('llm-interface');
    
    if (!llmInterface) return;

    if (llmConfig.enabled) {
      llmInterface.style.display = 'block';
      this.setupLLMButtons();
    } else {
      llmInterface.style.display = 'none';
    }
  }

  // 设置LLM按钮
  setupLLMButtons() {
    const adjustLayoutBtn = document.getElementById('llm-adjust-layout');
    const generateContentBtn = document.getElementById('llm-generate-content');
    const optimizeScheduleBtn = document.getElementById('llm-optimize-schedule');

    if (adjustLayoutBtn) {
      adjustLayoutBtn.addEventListener('click', async () => {
        const result = await window.LLMIntegration.adjustLayout();
        if (result) {
          console.log('Layout adjustment suggestion:', result);
          // 这里可以实现具体的布局调整逻辑
        }
      });
    }

    if (generateContentBtn) {
      generateContentBtn.addEventListener('click', async () => {
        const result = await window.LLMIntegration.generateContent();
        if (result) {
          // 更新告示牌内容
          const helpElement = document.getElementById('help');
          if (helpElement) {
            helpElement.innerHTML = result.replace(/\n/g, '<br>');
          }
        }
      });
    }

    if (optimizeScheduleBtn) {
      optimizeScheduleBtn.addEventListener('click', async () => {
        const result = await window.LLMIntegration.optimizeSchedule();
        if (result) {
          console.log('Schedule optimization suggestion:', result);
          // 这里可以显示优化建议
          alert(result);
        }
      });
    }
  }

  // 绑定事件
  bindEvents() {
    // 监听配置变更事件
    window.addEventListener('configChanged', (event) => {
      this.handleConfigChange(event.detail);
    });

    // 监听主题切换
    window.addEventListener('themeChanged', (event) => {
      this.handleThemeChange(event.detail);
    });
  }

  // 处理配置变更
  handleConfigChange(detail) {
    const { path, value } = detail;
    
    // 根据变更的配置路径执行相应的更新
    if (path.startsWith('ui.components')) {
      this.reinitializeComponent(path, value);
    } else if (path.startsWith('app')) {
      this.initializeApp();
    } else if (path.startsWith('llm')) {
      this.initializeLLMInterface();
    }
  }

  // 处理主题变更
  handleThemeChange(theme) {
    document.documentElement.setAttribute('data-theme', theme);
  }

  // 重新初始化特定组件
  reinitializeComponent(path, value) {
    const componentName = path.split('.')[2]; // ui.components.componentName
    
    switch (componentName) {
      case 'classTable':
        this.initializeClassTable();
        break;
      case 'clock':
        this.initializeClock();
        break;
      case 'calendar':
        this.initializeCalendar();
        break;
      case 'help':
        this.initializeHelp();
        break;
    }
  }

  // 动态更新样式
  updateStyles() {
    if (window.ThemeManager) {
      ThemeManager.applyTheme();
    }
  }
}

// 创建全局UI管理器实例
const uiManager = new UIManager();

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
  uiManager.init();
});

// 暴露到全局作用域
window.UIManager = uiManager;