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

// 新增判断当前是否有课的函数
// function hasCurrentClass(now, schedule) {
//     for (const period of schedule) {
//         const classBegin = parseTime(period.begin);
//         const classEnd = parseTime(period.end);
//         if (now >= classBegin && now <= classEnd) {
//             return true;
//         }
//     }
//     return false;
// }

// 主函数：刷新当前课程显示
function nowClass() {
    const date = new Date();
    const { today_vec, prev_vec, next_vec } = getDayVectors();
    const displayMode = CONFIG.lessons.displayMode || 'scroll';

    // 检查是否在学期时间范围
    const semesterBegin = new Date(CONFIG.lessons.times.semester.begin);
    const semesterEnd = new Date(CONFIG.lessons.times.semester.end);
    if (date < semesterBegin || date > semesterEnd) {
        for (let i = 0; i < 12; i++) {
            let opacity = (i === 6) ? "" : "opacity: 0.5;";
            document.getElementById('c' + i).innerHTML =
                `<a href="#" role="button" class="contrast" id="c_b${i}" style="${opacity}">无</a>`;
        }
        const c_b6 = document.getElementById('c_b6');
        c_b6.style.backgroundColor = '#93cee97f';
        c_b6.style.fontWeight = '600';
        c_b6.style.opacity = '1';
        for (let i = 0; i < 6; i++) {
            const el = document.getElementById('c_b' + i);
            el.style.backgroundColor = '#3daee940';
            el.style.fontWeight = '400';
        }
        for (let i = 7; i < 12; i++) {
            const el = document.getElementById('c_b' + i);
            el.style.backgroundColor = '';
            el.style.fontWeight = '400';
        }
        return;
    }

    const schedule = CONFIG.lessons.times.schedule;
    const now = date.getTime();
    let it = -1;
    let colorize = true;
    let inClassTime = false;
    let inRestTime = false;

    for (let i = 0; i < schedule.length; i++) {
        const period = schedule[i];
        const classBegin = parseTime(period.begin);
        const classEnd = parseTime(period.end);
        if (now >= classBegin && now <= classEnd) {
            it = i;
            inClassTime = true;
            const remainingTime = (classEnd - now) / (1000 * 60);
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

    if (displayMode === 'scroll') {
        // 滚动模式（新逻辑：智能填充课程）
        const { prevIdx, nextIdx } = findNearestClasses(now, schedule);
        const currentIndex = nextIdx !== -1 ? nextIdx : prevIdx;
        let arranged = arrangeClasses(currentIndex, today_vec, prev_vec, next_vec);

        // 新增休息状态处理
        const hasCurrent = inClassTime;
        if (!hasCurrent) {
            // 中间位置显示休息
            arranged[6] = "休息";
            
            // 向前追溯填充前一天课程
            for (let i = 5; i >= 0; i--) {
                if (!arranged[i] && prev_vec) {
                    const prevIndex = prev_vec.length - (6 - i);
                    arranged[i] = prevIndex >= 0 ? prev_vec[prevIndex] : "...";
                }
            }
            
            // 向后追溯填充后一天课程
            for (let i = 7; i < 12; i++) {
                if (!arranged[i] && next_vec) {
                    const nextIndex = i - 7;
                    arranged[i] = nextIndex < next_vec.length ? next_vec[nextIndex] : "...";
                }
            }
        }

        // 渲染课程表
        for (let i = 0; i < arranged.length; i++) {
            let content = arranged[i] || "";
            let opacity = (i === 6) ? "" : "opacity: 0.5;";
            document.getElementById('c' + i).innerHTML =
                `<a href="#" role="button" class="contrast" id="c_b${i}" style="${opacity}">${content}</a>`;
        }

        const c_b6 = document.getElementById('c_b6');
        c_b6.style.backgroundColor = '#93cee97f';
        c_b6.style.fontWeight = '600';
        c_b6.style.opacity = '1';
        for (let i = 0; i < 6; i++) {
            const el = document.getElementById('c_b' + i);
            el.style.backgroundColor = '#3daee940';
            el.style.fontWeight = '400';
        }
        for (let i = 7; i < 12; i++) {
            const el = document.getElementById('c_b' + i);
            el.style.backgroundColor = '';
            el.style.fontWeight = '400';
        }
    }
    else if (displayMode === 'day') {
        // 一天进度模式：只显示当天全部课程，前/当前/后课程有进度感
        const todayClasses = today_vec.slice(0, -1);
        let highlightIdx = -1;
        if (inClassTime) {
            highlightIdx = it;
            //schedule和todayClasses的差距修正
            highlightIdx++;
        } else {
            // 不在上课时间，定位最近上一节/下一节
            const { prevIdx, nextIdx } = findNearestClasses(now, schedule);
            highlightIdx = nextIdx; // 可选：也可不高亮任何课程
        }

        // console.log('[课程定位] 高亮课程',  todayClasses[highlightIdx]);

        for (let i = 0; i < todayClasses.length; i++) {
            let content = todayClasses[i] || "";
            let elStyle = '';
            if (highlightIdx === -1) {
                // 没有高亮，全部淡色
                elStyle = 'background-color:; font-weight:400; opacity:0.5;';
            } else if (i < highlightIdx) {
                // 已上过
                elStyle = 'background-color:#3daee940; font-weight:400; opacity:0.8;';
            } else if (i == highlightIdx && !inClassTime) { // 新增无课条件
                // 无课时
                elStyle = 'background-color:#3daee940; font-weight:400; opacity:0.8;';
            } else if (i === highlightIdx) {
                // 当前/即将上课
                elStyle = 'background-color:#93cee97f; font-weight:600; opacity:1;';
            } else {
                // 未上课
                elStyle = 'background-color:; font-weight:400; opacity:0.5;';
            }
            document.getElementById('c' + i).innerHTML =
                `<a href="#" role="button" class="contrast" id="c_b${i}" style="${elStyle}">${content}</a>`;
        }
    }
}

// 首次加载
nowClass();
// 每秒刷新
setInterval(nowClass, 1000);