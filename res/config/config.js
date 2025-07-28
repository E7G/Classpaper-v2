const CONFIG = {
  // 应用基础配置
  "app": {
    "title": "背景课表",
    "language": "zh",
    "theme": "dark",
    "favicon": "favicon.ico"
  },

  // 界面布局配置
  "ui": {
    "layout": {
      "containerFluid": true,
      "glassEffect": true,
      "overflow": "hidden"
    },
    "grid": {
      "enabled": true,
      "columns": "auto"
    },
    "components": {
      "classTable": {
        "enabled": true,
        "id": "classtable",
        "blurEffect": 12
      },
      "progressBar": {
        "enabled": true,
        "id": "nav-pro",
        "hidden": true
      },
      "clock": {
        "enabled": true,
        "id": "clockart",
        "showHeader": false,
        "headerText": "⏰ 时钟",
        "timeFormat": "zh",
        "dateFormat": "zh",
        "weekdayFormat": "long"
      },
      "calendar": {
        "enabled": true,
        "id": "calart",
        "showHeader": false,
        "headerText": "📅 事件日历"
      },
      "help": {
        "enabled": true,
        "id": "helpart",
        "showHeader": false,
        "headerText": "📚 告示牌",
        "columnCount": 2,
        "columnGap": "40px"
      },
    }
  },

  // 样式配置
  "styles": {
    "fonts": {
      "primary": "-apple-system, \"HarmonyOS Sans SC\", \"Noto Sans CJK SC\", MiSans, \"MiSans L3\", sans-serif, \"Apple Color Emoji\", \"Segoe UI Emoji\", \"Segoe UI Symbol\", \"Noto Color Emoji\"",
      "serif": "\"Noto Serif CJK SC\", serif"
    },
    "colors": {
      "primary": "#fcfcfc",
      "secondary": "#31363d30",
      "accent": "#3daee9",
      "background": "#7f8c8d",
      "text": "#fcfcfc",
      "mark": "black"
    },
    "effects": {
      "backdropBlur": 12,
      "borderRadius": 25,
      "glassOpacity": 0.3
    },
    "spacing": {
      "base": "5px",
      "padding": "15px 30px",
      "margin": "20px"
    },
    "fontSize": {
      "base": "36px",
      "large": "52px",
      "small": "24px"
    }
  },

  // 课程表配置
  "lessons": {
    "headers": [
      "星期",
      "早读",
      "第1节课",
      "第2节课", 
      "第3节课",
      "第4节课",
      "第5节课",
      "第6节课",
      "第7节课",
      "第8节课",
      "晚修1",
      "晚修2"
    ],
    "displayMode": "scroll",
    "schedule": [
      {
        "day": "周一",
        "classes": [
          "升旗",
          "班会",
          "物理",
          "英语",
          "数学",
          "数学",
          "生物",
          "化学",
          "语文",
          "晚修",
          "晚修"
        ]
      },
      {
        "day": "周二",
        "classes": [
          "语文",
          "语文",
          "数学",
          "英语",
          "体育",
          "物理",
          "生物",
          "化学",
          "自习",
          "晚修",
          "晚修"
        ]
      },
      {
        "day": "周三",
        "classes": [
          "英语",
          "物理",
          "英语",
          "数学",
          "生物",
          "语文",
          "化学",
          "英测",
          "英测",
          "晚修",
          "晚修"
        ]
      },
      {
        "day": "周四",
        "classes": [
          "语文",
          "物理",
          "语文",
          "生物",
          "体育",
          "英语",
          "化学",
          "数测",
          "数测",
          "晚修",
          "晚修"
        ]
      },
      {
        "day": "周五",
        "classes": [
          "英语",
          "化学",
          "生物",
          "数学",
          "语文",
          "语文",
          "自习",
          "物理",
          "英语",
          "无",
          "无"
        ]
      },
      {
        "day": "周六",
        "classes": [
          "无",
          "无",
          "无",
          "无",
          "无",
          "无",
          "无",
          "无",
          "无",
          "无",
          "无"
        ]
      },
      {
        "day": "周日",
        "classes": [
          "无",
          "无",
          "无",
          "无",
          "无",
          "无",
          "无",
          "无",
          "无",
          "晚修",
          "晚修"
        ]
      }
    ],
    "times": {
      "semester": {
        "begin": "2023-07-31",
        "end": "2026-06-07"
      },
      "schedule": [
        {
          "period": 1,
          "begin": "07:20",
          "end": "07:55",
          "rest": "07:55-08:00"
        },
        {
          "period": 2,
          "begin": "08:00",
          "end": "08:40",
          "rest": "08:40-08:50"
        },
        {
          "period": 3,
          "begin": "08:50",
          "end": "09:30",
          "rest": "09:30-09:40"
        },
        {
          "period": 4,
          "begin": "09:40",
          "end": "10:20",
          "rest": "10:20-10:45"
        },
        {
          "period": 5,
          "begin": "10:45",
          "end": "11:25",
          "rest": "11:25-11:35"
        },
        {
          "period": 6,
          "begin": "11:35",
          "end": "12:15",
          "rest": "12:15-14:20"
        },
        {
          "period": 7,
          "begin": "14:20",
          "end": "15:00",
          "rest": "15:00-15:15"
        },
        {
          "period": 8,
          "begin": "15:15",
          "end": "15:55",
          "rest": "15:55-16:05"
        },
        {
          "period": 9,
          "begin": "16:05",
          "end": "16:45",
          "rest": "16:45-19:00"
        },
        {
          "period": 10,
          "begin": "19:00",
          "end": "20:25",
          "rest": "20:25-20:40"
        },
        {
          "period": 11,
          "begin": "20:40",
          "end": "22:00",
          "rest": null
        }
      ]
    }
  },

  // 时间配置
  "time": {
    "weekOffset": {
      "enabled": true,
      "offset": 7
    },
    "updateInterval": 1000,
    "progressDescription": "高三剩余",
    "progressPercentMode": "left"
  },

  // 通知配置
  "notifications": {
    "enabled": true,
    "regularInterval": 5,
    "endingTime": 5,
    "sounds": {
      "regular": "audio/regular_notification.mp3",
      "ending": "audio/ending_notification.mp3"
    }
  },

  // 事件日历配置
  "events": [
    {
      "name": "高考",
      "date": "2026-06-07T00:00:00"
    },
    {
      "name": "明天", 
      "date": "2025-07-18T06:50:00"
    }
  ],

  // 壁纸配置
  "wallpapers": {
    "list": [
      "wallpaper/bg1.jpg",
      "wallpaper/bg2.jpg", 
      "wallpaper/kdedark.png",
      "wallpaper/kdelight.png"
    ],
    "interval": 30,
    "transition": {
      "duration": 1000,
      "easing": "cubic-bezier(0.4,0,0.2,1)"
    }
  },

  // 告示牌内容
  "announcement": "一鸣从此始，相望青云端",

  // LLM集成配置
  "llm": {
    "enabled": false,
    "apiEndpoint": "",
    "apiKey": "",
    "model": "gpt-3.5-turbo",
    "features": {
      "autoAdjustLayout": false,
      "smartScheduling": false,
      "contentGeneration": false,
      "voiceControl": false
    },
    "prompts": {
      "layoutAdjustment": "请根据当前时间和课程安排，调整界面布局以突出重要信息。",
      "contentGeneration": "请生成适合当前时间段的激励性文案。",
      "scheduleOptimization": "请分析当前课程安排，提供学习建议。"
    }
  }
};

// LLM集成工具类
class LLMIntegration {
  constructor(config) {
    this.config = config.llm;
    this.isEnabled = this.config.enabled;
  }

  async callLLM(prompt, context = {}) {
    if (!this.isEnabled || !this.config.apiEndpoint) {
      console.warn('LLM integration is disabled or not configured');
      return null;
    }

    try {
      const response = await fetch(this.config.apiEndpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.config.apiKey}`
        },
        body: JSON.stringify({
          model: this.config.model,
          messages: [
            {
              role: 'system',
              content: `你是一个智能课表助手。当前配置信息：${JSON.stringify(context)}`
            },
            {
              role: 'user', 
              content: prompt
            }
          ]
        })
      });

      const data = await response.json();
      return data.choices?.[0]?.message?.content || null;
    } catch (error) {
      console.error('LLM API call failed:', error);
      return null;
    }
  }

  async adjustLayout() {
    if (!this.config.features.autoAdjustLayout) return null;
    
    const context = {
      currentTime: new Date().toISOString(),
      currentConfig: CONFIG
    };
    
    return await this.callLLM(this.config.prompts.layoutAdjustment, context);
  }

  async generateContent() {
    if (!this.config.features.contentGeneration) return null;
    
    const context = {
      currentTime: new Date().toISOString(),
      announcement: CONFIG.announcement
    };
    
    return await this.callLLM(this.config.prompts.contentGeneration, context);
  }

  async optimizeSchedule() {
    if (!this.config.features.smartScheduling) return null;
    
    const context = {
      schedule: CONFIG.lessons.schedule,
      currentTime: new Date().toISOString()
    };
    
    return await this.callLLM(this.config.prompts.scheduleOptimization, context);
  }
}

// 配置管理工具类
class ConfigManager {
  static updateConfig(path, value) {
    const keys = path.split('.');
    let current = CONFIG;
    
    for (let i = 0; i < keys.length - 1; i++) {
      if (!current[keys[i]]) current[keys[i]] = {};
      current = current[keys[i]];
    }
    
    current[keys[keys.length - 1]] = value;
    this.notifyConfigChange(path, value);
  }

  static getConfig(path) {
    const keys = path.split('.');
    let current = CONFIG;
    
    for (const key of keys) {
      if (current[key] === undefined) return undefined;
      current = current[key];
    }
    
    return current;
  }

  static notifyConfigChange(path, value) {
    // 触发配置变更事件
    window.dispatchEvent(new CustomEvent('configChanged', {
      detail: { path, value, config: CONFIG }
    }));
  }

  static exportConfig() {
    return JSON.stringify(CONFIG, null, 2);
  }

  static importConfig(configString) {
    try {
      const newConfig = JSON.parse(configString);
      Object.assign(CONFIG, newConfig);
      this.notifyConfigChange('*', CONFIG);
      return true;
    } catch (error) {
      console.error('Failed to import config:', error);
      return false;
    }
  }
}

// 主题管理工具类
class ThemeManager {
  static applyTheme() {
    const theme = CONFIG.app.theme;
    document.documentElement.setAttribute('data-theme', theme);
    
    // 应用自定义样式
    const styles = CONFIG.styles;
    const root = document.documentElement;
    
    root.style.setProperty('--primary-color', styles.colors.primary);
    root.style.setProperty('--secondary-color', styles.colors.secondary);
    root.style.setProperty('--accent-color', styles.colors.accent);
    root.style.setProperty('--background-color', styles.colors.background);
    root.style.setProperty('--text-color', styles.colors.text);
    root.style.setProperty('--mark-color', styles.colors.mark);
    
    root.style.setProperty('--backdrop-blur', `${styles.effects.backdropBlur}px`);
    root.style.setProperty('--border-radius', `${styles.effects.borderRadius}px`);
    root.style.setProperty('--glass-opacity', styles.effects.glassOpacity);
    
    root.style.setProperty('--base-spacing', styles.spacing.base);
    root.style.setProperty('--base-padding', styles.spacing.padding);
    root.style.setProperty('--base-margin', styles.spacing.margin);
    
    root.style.setProperty('--base-font-size', styles.fontSize.base);
    root.style.setProperty('--large-font-size', styles.fontSize.large);
    root.style.setProperty('--small-font-size', styles.fontSize.small);
    
    root.style.setProperty('--primary-font', styles.fonts.primary);
    root.style.setProperty('--serif-font', styles.fonts.serif);
  }

  static switchTheme(newTheme) {
    ConfigManager.updateConfig('app.theme', newTheme);
    this.applyTheme();
  }
}

// 全局实例
const llmIntegration = new LLMIntegration(CONFIG);

// 暴露到全局作用域
window.CONFIG = CONFIG;
window.ConfigManager = ConfigManager;
window.ThemeManager = ThemeManager;
window.LLMIntegration = llmIntegration;

// 为了保持向后兼容，导出原有的变量名
const lessons = CONFIG.lessons.headers.join(",") + "\n" + 
  CONFIG.lessons.schedule.map(day => " ," + day.day + "," + day.classes.join(",")).join("\n") + "\n";

const events = "事件,日期,\n" + 
  CONFIG.events.map(event => `${event.name},${event.date},`).join("\n");

const wallpaperlist = CONFIG.wallpapers.list;

const sth = CONFIG.announcement;