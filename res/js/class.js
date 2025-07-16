/*课表 */

var source = lessons;

let source_vec = source.split(',');
//console.log(source_vec)

var classtable = document.getElementById('classtable');

// 音频元素
const regularNotification = document.getElementById('regularNotification');
const endingNotification = document.getElementById('endingNotification');

// 上次提示的时间
let lastNotificationTime = 0;
// 是否已经发出结束提示
let endingNotified = false;

// 播放提示音的函数
function playNotification(type) {
    // 检查提示音是否启用
    if (!CONFIG.notifications.enabled) {
        return;
    }

    if (type === 'regular') {
        regularNotification.play();
    } else if (type === 'ending') {
        endingNotification.play();
    }
}

var date = new Date();

var week = new Date().getDay();

var offset = week * 12
if (week == 0) {
    offset = 7 * 12;
}

var today_vec = source_vec.slice(offset, offset + 13);
var prev_vec = offset >= 12 ? source_vec.slice(offset - 12, offset + 1) : null;
var next_vec = offset + 25 < source_vec.length ? source_vec.slice(offset + 13, offset + 25) : null;

// 根据当前课程重新排列课程顺序
function arrangeClasses(currentIndex) {
    // 复制今天的课程数组
    let classes = today_vec.slice(0, -1);  // 去掉最后一个空元素
    
    // 计算需要在当前课程前显示多少节课
    const showBefore = 6;  // 当前课程要显示在第7个位置
    
    // 构建新的显示顺序
    let arranged = [];
    
    // 添加当前课程之前的课程
    for (let i = showBefore - 1; i >= 0; i--) {
        let index = currentIndex - (showBefore - i);
        if (index >= 0) {
            arranged[i] = classes[index];
        } else {
            // 使用前一天对应时间的课程
            if (prev_vec) {
                let prevIndex = prev_vec.length - 1 + index; // -1 因为要去掉最后一个空元素
                arranged[i] = prevIndex >= 0 ? prev_vec[prevIndex] : "";
            } else {
                arranged[i] = "";
            }
        }
    }
    
    // 添加当前课程
    arranged[showBefore] = classes[currentIndex];
    
    // 添加当前课程之后的课程
    for (let i = showBefore + 1; i < 12; i++) {
        let index = currentIndex + (i - showBefore);
        if (index < classes.length) {
            arranged[i] = classes[index];
        } else {
            // 使用下一天对应时间的课程
            if (next_vec) {
                let nextIndex = index - classes.length;
                arranged[i] = nextIndex < next_vec.length - 1 ? next_vec[nextIndex] : "";
            } else {
                arranged[i] = "";
            }
        }
    }
    
    return arranged;
}

// 解析时间字符串为Date对象
function parseTime(timeStr) {
    const [hours, minutes] = timeStr.split(':').map(Number);
    const now = new Date();
    now.setHours(hours, minutes, 0);
    return now;
}

function nowClass() {
    var date = new Date();
    
    // 检查是否在学期时间范围内
    const semesterBegin = new Date(CONFIG.lessons.times.semester.begin);
    const semesterEnd = new Date(CONFIG.lessons.times.semester.end);
    
    if (date < semesterBegin || date > semesterEnd) {
        // 不在学期时间范围内，显示空课表
        return;
    }

    var colorize = true;
    var it = 0;

    // 获取当前时间对应的课程
    const schedule = CONFIG.lessons.times.schedule;
    const now = date.getTime();

    // 遍历课程时间表，找到当前时间对应的课程
    for (let i = schedule.length - 1; i >= 0; i--) {
        const period = schedule[i];
        const classBegin = parseTime(period.begin);
        const classEnd = parseTime(period.end);
        
        if (now >= classBegin && now <= classEnd) {
            it = i;
            // 课程倒计时提示
            const remainingTime = (classEnd - date) / (1000 * 60); // 剩余分钟数
            
            // 根据设置的间隔时间提示
            const regularInterval = CONFIG.notifications.regularInterval;
            if (remainingTime > CONFIG.notifications.endingTime && 
                (now - lastNotificationTime) >= regularInterval * 60 * 1000) {
                playNotification('regular');
                lastNotificationTime = now;
            }
            
            // 根据设置的结束时间提示
            if (remainingTime <= CONFIG.notifications.endingTime && !endingNotified) {
                playNotification('ending');
                endingNotified = true;
            }
            break;
        } else if (period.rest && i > 0) {
            const restBegin = schedule[i - 1].end;
            const restTime = parseTime(restBegin);
            if (now >= classEnd && now <= restTime) {
                it = i;
                colorize = false;
                // 重置提示状态
                endingNotified = false;
                break;
            }
        }
    }

    // 重新排列课程
    let arranged = arrangeClasses(it);
    
    // 更新显示
    for (let i = 0; i < arranged.length; i++) {
        let content = arranged[i] || "";
        let opacity = i < 6 ? "opacity: 0.5;" : (i > 6 ? "opacity: 0.5;" : "");
        document.getElementById('c' + i).innerHTML = 
            `<a href="#" role="button" class="contrast" id="c_b${i}" style="${opacity}">${content}</a>`;
    }

    // 设置当前课程样式
    document.getElementById('c_b6').style.backgroundColor = colorize ? '#93cee97f' : '#3daee940';
    document.getElementById('c_b6').style.fontWeight = colorize ? '600' : '400';
    document.getElementById('c_b6').style.opacity = '1';

    // 设置之前课程的样式
    for (let i = 0; i < 6; i++) {
        document.getElementById('c_b' + i).style.backgroundColor = '#3daee940';
        document.getElementById('c_b' + i).style.fontWeight = '400';
    }

    // 清除之后课程的样式
    for (let i = 7; i < 12; i++) {
        document.getElementById('c_b' + i).style.backgroundColor = '';
        document.getElementById('c_b' + i).style.fontWeight = '400';
    }
}

nowClass();

setInterval(nowClass, 1000);