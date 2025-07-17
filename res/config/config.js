const CONFIG = {
  "lessons": {
    "headers": [
      "星期",
      "1",
      "2",
      "3",
      "4",
      "5",
      "6",
      "7",
      "8",
      "9",
      "10",
      "11"
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
  "weekOffset": {
    "enabled": true,
    "offset": 7
  },
  "notifications": {
    "enabled": true,
    "regularInterval": 5,
    "endingTime": 5
  },
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
  "wallpapers": [
    "wallpaper/bg1.jpg",
    "wallpaper/bg2.jpg",
    "wallpaper/kdedark.png",
    "wallpaper/kdelight.png"
  ],
  "wallpaperInterval": 30,
  "progressDescription": "高三剩余",
  "progressPercentMode": "left",
  "sth": "一鸣从此始，相望青云端"
};

// 为了保持向后兼容，导出原有的变量名
const lessons = CONFIG.lessons.headers.join(",") + "\n" + 
  CONFIG.lessons.schedule.map(day => " ," + day.day + "," + day.classes.join(",")).join("\n") + "\n";

const events = "事件,日期,\n" + 
  CONFIG.events.map(event => `${event.name},${event.date},`).join("\n");

const wallpaperlist = CONFIG.wallpapers;

const sth = CONFIG.sth;