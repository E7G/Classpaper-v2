const WIDTH = 1920;
const HEIGHT = 1920;

// 获取body元素
const body = document.querySelector('body');

// 优化壁纸显示，提升加载与切换体验
function changeWallpaper() {
  // 随机选择壁纸
  const wallpaperIndex = Math.floor(Math.random() * wallpaperlist.length);
  const wallpaper = wallpaperlist[wallpaperIndex];

  // 预加载壁纸图片，确保切换时已加载完成
  const img = new Image();
  img.src = wallpaper;
  img.onload = function() {
    // 创建一个新的div用于过渡
    let transitionDiv = document.createElement('div');
    transitionDiv.className = 'wallpaper-transition-optimized';
    transitionDiv.style.position = 'fixed';
    transitionDiv.style.left = 0;
    transitionDiv.style.top = 0;
    transitionDiv.style.width = '100vw';
    transitionDiv.style.height = '100vh';
    transitionDiv.style.zIndex = -1;
    transitionDiv.style.backgroundImage = `url(${wallpaper})`;
    transitionDiv.style.backgroundRepeat = 'no-repeat';
    transitionDiv.style.backgroundPosition = 'center center';
    transitionDiv.style.backgroundSize = 'cover';
    transitionDiv.style.opacity = 0;
    transitionDiv.style.transition = 'opacity 1s cubic-bezier(0.4,0,0.2,1)';

    // 插入到body最底层
    body.insertBefore(transitionDiv, body.firstChild);

    // 触发过渡
    setTimeout(() => {
      transitionDiv.style.opacity = 1;
    }, 10);

    // 过渡完成后移除旧背景
    setTimeout(() => {
      // 移除之前的背景div（如果有）
      let oldDivs = document.querySelectorAll('.wallpaper-transition-optimized:not(:first-child)');
      oldDivs.forEach(div => div.remove());

      // 设置body的背景为当前壁纸（用于后续新div移除时仍有背景）
      body.style.backgroundImage = `url(${wallpaper})`;
      body.style.backgroundRepeat = 'no-repeat';
      body.style.backgroundPosition = 'center center';
      body.style.backgroundSize = 'cover';
    }, 1100);
  };
}

// 初始化时调用一次切换壁纸函数
changeWallpaper();

// 每30秒调用一次切换壁纸函数
setInterval(changeWallpaper, 30*1000);

// 监听窗口大小变化事件，当窗口大小发生变化时重新调用切换壁纸函数
//window.addEventListener('resize', changeWallpaper);

setInterval(setTime, 1000);
//计算本年的周数
function getYearWeek(endDate) {
    //本年的第一天

    var beginDate = new Date(endDate.getFullYear(), 0, 1);
    //星期从0-6,0代表星期天，6代表星期六

    var endWeek = endDate.getDay();

    //if (endWeek == 0) endWeek = 7;

    var beginWeek = beginDate.getDay();

    if (beginWeek == 0) beginWeek = 7;
    //计算两个日期的天数差

    var millisDiff = endDate.getTime() - beginDate.getTime();
    var dayDiff = Math.floor((millisDiff + (beginWeek - endWeek) * (24 * 60 * 60 * 1000)) / 86400000);
    return Math.ceil(dayDiff / 7) + 1;
}
function setTime() {
    var date = new Date();
    // 从配置中获取学期起始和结束时间
    var beginDate = new Date(CONFIG.lessons.times.semester.begin);
    var endDate = new Date(CONFIG.lessons.times.semester.end);
    var msPassed = date - beginDate;
    var msOverall = endDate - beginDate;
    var percentageLeft = 100 - msPassed / msOverall * 100;
    var week = getYearWeek(date);
    // 使用配置中的周数偏移
    var week_term = CONFIG.weekOffset.enabled ? week - CONFIG.weekOffset.offset : week;
    document.getElementById("time-a").innerHTML = date.toLocaleTimeString('zh');
    document.getElementById("date-a").innerHTML = date.toLocaleDateString('zh');
    document.getElementById("weekday-a").innerHTML = date.toLocaleDateString('zh', { weekday: 'long' });
    document.getElementById("week-a").innerHTML = '第 ' + week_term + ' 周';
    document.getElementById('prog').setAttribute('max', msOverall);
    document.getElementById('prog').setAttribute('value', msPassed);
    document.getElementById('prog-description').innerHTML = "高三剩余 " + percentageLeft.toFixed(4) + "%";
    //document.getElementById("nav-time").insertAdjacentHTML('beforeend', '<a href="#" role="button"><h1>' +date.toLocaleTimeString('zh', { hour12: true })  + '</h1></div>');
}

setTime();
console.log("屏幕大小：" + WIDTH + '*' + HEIGHT);