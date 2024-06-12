
// 获取body元素
const body = document.querySelector('body');

// 定义切换壁纸函数
function changeWallpaper() {
  // 随机选择壁纸
  const wallpaperIndex = Math.floor(Math.random() * wallpaperlist.length);
  const wallpaper = wallpaperlist[wallpaperIndex];

  // 设置背景图片
  body.style.backgroundImage = `url(${wallpaper})`;
  body.style.backgroundRepeat = 'no-repeat';
  body.style.backgroundPosition = 'center';
  body.style.backgroundSize = 'cover';
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
    
if (endWeek == 0) endWeek = 7;
    
var beginWeek = beginDate.getDay();
    
if (beginWeek == 0) beginWeek = 7;
    //计算两个日期的天数差
    
    var millisDiff = endDate.getTime() - beginDate.getTime();
    var dayDiff = Math.floor(( millisDiff + (beginWeek - endWeek) * (24 * 60 * 60 * 1000)) / 86400000);
    return Math.ceil(dayDiff / 7) + 1;
} 
function setTime() {
    var date = new Date();
    var week = getYearWeek(date);
    var week_term = week - 31;
    document.getElementById("time-a").innerHTML = date.toLocaleTimeString('zh', { hour12: true });
    
document.getElementById("date-a").innerHTML = date.toLocaleDateString('zh');
    document.getElementById("week-a").innerHTML = '第' + week_term + '周';
    //document.getElementById("nav-time").insertAdjacentHTML('beforeend', '<a href="#" role="button"><h1>' +date.toLocaleTimeString('zh', { hour12: true })  + '</h1></div>');
}

setTime();