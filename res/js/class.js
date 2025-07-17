/* 课表优化版 */

const source = lessons;
const source_vec = source.split(',');

const classtable = document.getElementById('classtable');

// 音频元素
const regularNotification = document.getElementById('regularNotification');
const endingNotification = document.getElementById('endingNotification');

// 上次提示的时间戳
let lastNotificationTime = 0;
// 是否已经发出结束提示
let endingNotified = false;

// 播放提示音
function playNotification(type) {
    if (!CONFIG.notifications.enabled) return;
    if (type === 'regular') {
        regularNotification && regularNotification.play();
    } else if (type === 'ending') {
        endingNotification && endingNotification.play();
    }
}

// 获取今天、前一天、后一天的课程数组
function getDayVectors() {
    const week = new Date().getDay();
    let offset = week === 0 ? 7 * 12 : week * 12;
    // 13是因为每天12节课+1个空元素
    return {
        today_vec: source_vec.slice(offset, offset + 13),
        prev_vec: offset >= 12 ? source_vec.slice(offset - 12, offset + 1) : null,
        next_vec: offset + 25 < source_vec.length ? source_vec.slice(offset + 13, offset + 25) : null
    };
}

// 重新排列课程顺序，当前课程在第7个位置
function arrangeClasses(currentIndex, today_vec, prev_vec, next_vec) {
    const showBefore = 6;
    const classes = today_vec.slice(0, -1);
    const arranged = new Array(12);

    // 前6节
    for (let i = 0; i < showBefore; i++) {
        let idx = currentIndex - (showBefore - i);
        if (idx >= 0) {
            arranged[i] = classes[idx];
        } else if (prev_vec) {
            let prevIndex = prev_vec.length - 1 + idx;
            arranged[i] = prevIndex >= 0 ? prev_vec[prevIndex] : "";
        } else {
            arranged[i] = "";
        }
    }
    // 当前
    arranged[showBefore] = classes[currentIndex];
    // 后5节
    for (let i = showBefore + 1; i < 12; i++) {
        let idx = currentIndex + (i - showBefore);
        if (idx < classes.length) {
            arranged[i] = classes[idx];
        } else if (next_vec) {
            let nextIndex = idx - classes.length;
            arranged[i] = nextIndex < next_vec.length - 1 ? next_vec[nextIndex] : "";
        } else {
            arranged[i] = "";
        }
    }
    return arranged;
}

// 解析时间字符串为当天的Date对象
function parseTime(timeStr) {
    const [hours, minutes] = timeStr.split(':').map(Number);
    const now = new Date();
    now.setHours(hours, minutes, 0, 0);
    return now.getTime();
}

// 查找最近的上一节课和下一节课
function findNearestClasses(now, schedule) {
    let prevIdx = -1, nextIdx = -1;
    let prevTime = -Infinity, nextTime = Infinity;
    for (let i = 0; i < schedule.length; i++) {
        const classBegin = parseTime(schedule[i].begin);
        const classEnd = parseTime(schedule[i].end);
        if (classEnd <= now && classEnd > prevTime) {
            prevTime = classEnd;
            prevIdx = i;
        }
        if (classBegin > now && classBegin < nextTime) {
            nextTime = classBegin;
            nextIdx = i;
        }
    }
    return { prevIdx, nextIdx };
}

// 主函数：刷新当前课程显示
function nowClass() {
    const date = new Date();
    const { today_vec, prev_vec, next_vec } = getDayVectors();

    // 检查是否在学期时间范围
    const semesterBegin = new Date(CONFIG.lessons.times.semester.begin);
    const semesterEnd = new Date(CONFIG.lessons.times.semester.end);
    if (date < semesterBegin || date > semesterEnd) {
        // 非学期内，全部显示“无”
        for (let i = 0; i < 12; i++) {
            let opacity = (i === 6) ? "" : "opacity: 0.5;";
            document.getElementById('c' + i).innerHTML =
                `<a href="#" role="button" class="contrast" id="c_b${i}" style="${opacity}">无</a>`;
        }
        // 当前课程高亮
        const c_b6 = document.getElementById('c_b6');
        c_b6.style.backgroundColor = '#93cee97f';
        c_b6.style.fontWeight = '600';
        c_b6.style.opacity = '1';
        // 之前课程样式
        for (let i = 0; i < 6; i++) {
            const el = document.getElementById('c_b' + i);
            el.style.backgroundColor = '#3daee940';
            el.style.fontWeight = '400';
        }
        // 之后课程样式
        for (let i = 7; i < 12; i++) {
            const el = document.getElementById('c_b' + i);
            el.style.backgroundColor = '';
            el.style.fontWeight = '400';
        }
        return;
    }

    // 获取当前时间对应的课程
    const schedule = CONFIG.lessons.times.schedule;
    const now = date.getTime();
    let it = -1;
    let colorize = true;
    let inClassTime = false;
    let inRestTime = false;

    // 判断当前是否在课程时间或课间休息
    for (let i = 0; i < schedule.length; i++) {
        const period = schedule[i];
        const classBegin = parseTime(period.begin);
        const classEnd = parseTime(period.end);

        if (now >= classBegin && now <= classEnd) {
            it = i;
            inClassTime = true;
            // 课程倒计时提示
            const remainingTime = (classEnd - now) / (1000 * 60); // 剩余分钟数
            const regularInterval = CONFIG.notifications.regularInterval;
            if (remainingTime > CONFIG.notifications.endingTime &&
                (now - lastNotificationTime) >= regularInterval * 60 * 1000) {
                playNotification('regular');
                lastNotificationTime = now;
            }
            if (remainingTime <= CONFIG.notifications.endingTime && !endingNotified) {
                playNotification('ending');
                endingNotified = true;
            }
            break;
        } else if (period.rest && i > 0) {
            // 课间休息
            const restBegin = schedule[i - 1].end;
            const restTime = parseTime(restBegin);
            if (now >= classEnd && now <= restTime) {
                it = i;
                colorize = false;
                endingNotified = false;
                inRestTime = true;
                break;
            }
        }
    }

    if (inClassTime || inRestTime) {
        // 课中或课间，正常显示
        const arranged = arrangeClasses(it, today_vec, prev_vec, next_vec);

        for (let i = 0; i < arranged.length; i++) {
            let content = arranged[i] || "";
            let opacity = (i === 6) ? "" : "opacity: 0.5;";
            document.getElementById('c' + i).innerHTML =
                `<a href="#" role="button" class="contrast" id="c_b${i}" style="${opacity}">${content}</a>`;
        }

        // 当前课程高亮
        const c_b6 = document.getElementById('c_b6');
        c_b6.style.backgroundColor = colorize ? '#93cee97f' : '#3daee940';
        c_b6.style.fontWeight = colorize ? '600' : '400';
        c_b6.style.opacity = '1';

        // 之前课程样式
        for (let i = 0; i < 6; i++) {
            const el = document.getElementById('c_b' + i);
            el.style.backgroundColor = '#3daee940';
            el.style.fontWeight = '400';
        }
        // 之后课程样式
        for (let i = 7; i < 12; i++) {
            const el = document.getElementById('c_b' + i);
            el.style.backgroundColor = '';
            el.style.fontWeight = '400';
        }
    } else {
        // 不在课程时间，显示“暂时无课/休息”，前后显示最近的上一节课和下一节课
        const { prevIdx, nextIdx } = findNearestClasses(now, schedule);
        const classes = today_vec.slice(0, -1);

        let prevClass = prevIdx >= 0 ? classes[prevIdx] : "";
        let nextClass = nextIdx >= 0 ? classes[nextIdx] : "";

        // 构造显示数组
        let arranged = new Array(12).fill("");

        // 填充上一节、下一节
        arranged[5] = prevClass || "无";
        arranged[6] = "休息";
        arranged[7] = nextClass || "无";

        // 向前填充上上一节、上上上一节...
        let p = prevIdx - 1;
        for (let i = 4; i >= 0; i--) {
            if (p >= 0) {
                arranged[i] = classes[p];
                p--;
            } else {
                arranged[i] = "";
            }
        }
        // 向后填充下下一节、下下下一节...
        let n = nextIdx + 1;
        for (let i = 8; i < 12; i++) {
            if (n < classes.length) {
                arranged[i] = classes[n];
                n++;
            } else {
                arranged[i] = "";
            }
        }

        for (let i = 0; i < arranged.length; i++) {
            let content = arranged[i] || "";
            let opacity = (i === 6) ? "" : "opacity: 0.5;";
            document.getElementById('c' + i).innerHTML =
                `<a href="#" role="button" class="contrast" id="c_b${i}" style="${opacity}">${content}</a>`;
        }

        // 当前课程高亮
        const c_b6 = document.getElementById('c_b6');
        c_b6.style.backgroundColor = '#3daee940';
        c_b6.style.fontWeight = '600';
        c_b6.style.opacity = '1';

        // 之前课程样式
        for (let i = 0; i < 6; i++) {
            const el = document.getElementById('c_b' + i);
            el.style.backgroundColor = '#3daee940';
            el.style.fontWeight = '400';
        }
        // 之后课程样式
        for (let i = 7; i < 12; i++) {
            const el = document.getElementById('c_b' + i);
            el.style.backgroundColor = '';
            el.style.fontWeight = '400';
        }
    }
}

// 首次加载
nowClass();
// 每秒刷新
setInterval(nowClass, 1000);