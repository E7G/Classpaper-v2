const WIDTH = window.screen.width || 1920;
const HEIGHT = window.screen.height || 1080;

// 获取body元素
const body = document.body;

// 优化壁纸显示，提升加载与切换体验
function changeWallpaper() {
  // 随机选择壁纸
  const wallpaperIndex = Math.floor(Math.random() * wallpaperlist.length);
  const wallpaper = wallpaperlist[wallpaperIndex];

  // 预加载壁纸图片，确保切换时已加载完成
  const img = new window.Image();
  img.src = wallpaper;
  img.onload = function () {
    // 创建一个新的div用于过渡
    const transitionDiv = document.createElement('div');
    transitionDiv.className = 'wallpaper-transition-optimized';
    Object.assign(transitionDiv.style, {
      position: 'fixed',
      left: 0,
      top: 0,
      width: '100vw',
      height: '100vh',
      zIndex: -1,
      backgroundImage: `url(${wallpaper})`,
      backgroundRepeat: 'no-repeat',
      backgroundPosition: 'center center',
      backgroundSize: 'cover',
      opacity: 0,
      transition: 'opacity 1s cubic-bezier(0.4,0,0.2,1)'
    });

    // 插入到body最底层
    body.insertBefore(transitionDiv, body.firstChild);

    // 触发过渡
    requestAnimationFrame(() => {
      transitionDiv.style.opacity = 1;
    });

    // 过渡完成后移除旧背景
    setTimeout(() => {
      // 移除之前的背景div（如果有）
      const oldDivs = document.querySelectorAll('.wallpaper-transition-optimized:not(:first-child)');
      oldDivs.forEach(div => div.remove());

      // 设置body的背景为当前壁纸（用于后续新div移除时仍有背景）
      Object.assign(body.style, {
        backgroundImage: `url(${wallpaper})`,
        backgroundRepeat: 'no-repeat',
        backgroundPosition: 'center center',
        backgroundSize: 'cover',
        backgroundAttachment: 'fixed',  // 固定背景防止滚动
        overflow: 'hidden',  // 隐藏滚动条
        margin: '0',  // 清除默认边距
        padding: '0'   // 清除默认内边距
      });
    }, 1100);
  };
}

let wallpaperTimer = null;
function startWallpaperTimer() {
  if (wallpaperTimer) clearInterval(wallpaperTimer);
  let interval = (CONFIG.wallpaperInterval || 30) * 1000;
  wallpaperTimer = setInterval(changeWallpaper, interval);
}

// 初始化时调用一次切换壁纸函数
changeWallpaper();
startWallpaperTimer();

// 如果有全局配置变更（如保存设置后），可暴露 window.reloadWallpaperTimer = startWallpaperTimer;
window.reloadWallpaperTimer = startWallpaperTimer;

// 定时更新时间
setInterval(setTime, 1000);

// 计算本年的周数
function getYearWeek(date) {
  // 本年的第一天
  const beginDate = new Date(date.getFullYear(), 0, 1);
  let beginWeek = beginDate.getDay();
  if (beginWeek === 0) beginWeek = 7;
  let endWeek = date.getDay();
  // 计算两个日期的天数差
  const millisDiff = date.getTime() - beginDate.getTime();
  const dayDiff = Math.floor((millisDiff + (beginWeek - endWeek) * 86400000) / 86400000);
  return Math.ceil(dayDiff / 7) + 1;
}

function setTime() {
  const date = new Date();
  // 从配置中获取学期起始和结束时间
  const beginDate = new Date(CONFIG.lessons.times.semester.begin);
  const endDate = new Date(CONFIG.lessons.times.semester.end);
  const msPassed = date - beginDate;
  const msOverall = endDate - beginDate;
  const percentageLeft = 100 - msPassed / msOverall * 100;
  const week = getYearWeek(date);
  // 使用配置中的周数偏移
  const week_term = CONFIG.weekOffset.enabled ? week - CONFIG.weekOffset.offset : week;
  document.getElementById("time-a").textContent = date.toLocaleTimeString('zh');
  document.getElementById("date-a").textContent = date.toLocaleDateString('zh');
  document.getElementById("weekday-a").textContent = date.toLocaleDateString('zh', { weekday: 'long' });
  document.getElementById("week-a").textContent = '第 ' + week_term + ' 周';
  const prog = document.getElementById('prog');
  prog.max = msOverall;
  prog.value = msPassed;
  let percentText = '';
  if ((CONFIG.progressPercentMode || 'left') === 'left') {
    percentText = (CONFIG.progressDescription || '高三剩余') + ' ' + (100 - msPassed / msOverall * 100).toFixed(4) + '%';
  } else {
    percentText = (CONFIG.progressDescription || '高三已过') + ' ' + (msPassed / msOverall * 100).toFixed(4) + '%';
  }
  document.getElementById('prog-description').textContent = percentText;
}

setTime();
console.log(`屏幕大小：${WIDTH}*${HEIGHT}`);