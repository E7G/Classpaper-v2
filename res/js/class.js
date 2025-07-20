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
function playNotification(type, className) {
    if (!CONFIG.notifications.enabled) return;
    if (!className || className === '无') return; // 课程为无时不响铃
    if (type === 'regular') {
        regularNotification && regularNotification.play();
    } else if (type === 'ending') {
        endingNotification && endingNotification.play();
    }
}

// 获取今天、前一天、后一天的课程数组，并带有详细注释
function getDayVectors() {
    // 获取当前星期几，getDay()返回0-6，0表示周日
    let week = new Date().getDay();
    // 如果是周日（0），则将其转换为7，方便后续计算（1=周一，7=周日）
    week = week === 0 ? 7 : week;

    const dayCount = 7;           // 一周有7天
    const lessonsPerDay = 12;     // 每天有12节课
    const total = (dayCount+1) * lessonsPerDay; // 一周总共的课程数

    // 计算今天在课程数组中的起始下标
    // 例如：周一为12，周二为24，依此类推
    let offset =  week  * lessonsPerDay;
    // console.log(`[getDayVectors] week: ${week}, offset: ${offset}`);

    // 取出今天的课程数组
    // 这里slice的长度为13，是因为每天12节课+1个空元素（可能用于占位或防止越界）
    let today_vec = source_vec.slice(offset, offset + lessonsPerDay).map(item => item.replace(/\n/g, ""));
    today_vec.push("");
    // console.log(`[getDayVectors] today_vec:`, today_vec);

    // 计算前一天的起始下标
    let prev_offset = offset - lessonsPerDay;
    // 如果前一天小于0，说明已经到周一的前一天了，需要循环到周日
    if (prev_offset < 1*lessonsPerDay) prev_offset = total-lessonsPerDay;
    // console.log(`[getDayVectors] prev_offset: ${prev_offset}`);

    // 计算后一天的起始下标
    let next_offset = offset + lessonsPerDay;
    // 如果后一天超出总课程数，说明已经到周日的后一天了，需要循环到周一
    if (next_offset >= total) next_offset = 1*lessonsPerDay;
    // console.log(`[getDayVectors] next_offset: ${next_offset}`);

    // 取出前一天的课程数组，并在末尾加一个空元素
    let prev_vec = source_vec.slice(prev_offset, prev_offset + lessonsPerDay).map(item => item.replace(/\n/g, ""));
    prev_vec.push("");
    // 取出后一天的课程数组，并在末尾加一个空元素
    let next_vec = source_vec.slice(next_offset, next_offset + lessonsPerDay).map(item => item.replace(/\n/g, ""));
    next_vec.push("");

    // console.log(`[getDayVectors] prev_vec:`, prev_vec);
    // console.log(`[getDayVectors] next_vec:`, next_vec);

    // 返回今天、前一天、后一天的课程数组
    return {
        today_vec, // 今天的课程
        prev_vec,  // 前一天的课程
        next_vec   // 后一天的课程
    };
}

// 重新排列课程顺序，当前课程在第7个位置
function arrangeClasses(currentIndex, today_vec, prev_vec, next_vec) {
    // showBefore 表示当前课程前面要显示的课程数，这里为6
    const showBefore = 6;
    // classes 是今天的课程数组，去掉最后一个元素（通常为占位或空元素）
    const classes = today_vec.slice(0, -1);
    // arranged 用于存放最终排列好的12节课
    const arranged = new Array(12);

    // console.log(`[arrangeClasses] currentIndex: ${currentIndex}`);
    // console.log(`[arrangeClasses] today_vec:`, today_vec);
    // console.log(`[arrangeClasses] prev_vec:`, prev_vec);
    // console.log(`[arrangeClasses] next_vec:`, next_vec);

    // 填充前6节课
    for (let i = 0; i < showBefore; i++) {
        // idx 计算当前要填充的课程在 today_vec 中的下标
        let idx = currentIndex - (showBefore - i);
        if (idx >= 0) {
            // 如果 idx 合法，直接取今天的课程
            arranged[i] = classes[idx];
            // console.log(`[arrangeClasses] arranged[${i}] = today_vec[${idx}]`, classes[idx]);
        } else if (prev_vec) {
            // 如果 idx 不合法，尝试从前一天的课程补齐
            let prevIndex = prev_vec.length - 1 + idx;
            if (prevIndex >= 0) {
                arranged[i] = prev_vec[prevIndex];
                // console.log(`[arrangeClasses] arranged[${i}] = prev_vec[${prevIndex}]`, prev_vec[prevIndex]);
            } else {
                arranged[i] = "";
                // console.log(`[arrangeClasses] arranged[${i}] = "" (prevIndex < 0)`);
            }
        } else {
            // 没有前一天的课程，填空
            arranged[i] = "";
            // console.log(`[arrangeClasses] arranged[${i}] = "" (no prev_vec)`);
        }
    }

    // 填充当前课程
    arranged[showBefore] = classes[currentIndex];
    // console.log(`[arrangeClasses] arranged[${showBefore}] = today_vec[${currentIndex}]`, classes[currentIndex]);

    // 填充后5节课
    for (let i = showBefore + 1; i < 12; i++) {
        // idx 计算当前要填充的课程在 today_vec 中的下标
        let idx = currentIndex + (i - showBefore);
        if (idx < classes.length) {
            // 如果 idx 合法，直接取今天的课程
            arranged[i] = classes[idx];
            // console.log(`[arrangeClasses] arranged[${i}] = today_vec[${idx}]`, classes[idx]);
        } else if (next_vec) {
            // 如果 idx 超出 today_vec，尝试从后一天的课程补齐
            let nextIndex = idx - classes.length;
            if (nextIndex < next_vec.length - 1) {
                arranged[i] = next_vec[nextIndex];
                // console.log(`[arrangeClasses] arranged[${i}] = next_vec[${nextIndex}]`, next_vec[nextIndex]);
            } else {
                arranged[i] = "";
                // console.log(`[arrangeClasses] arranged[${i}] = "" (nextIndex 超界)`);
            }
        } else {
            // 没有后一天的课程，填空
            arranged[i] = "";
            // console.log(`[arrangeClasses] arranged[${i}] = "" (no next_vec)`);
        }
    }

    // console.log(`[arrangeClasses] arranged:`, arranged);
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
            // 获取当前课程名
            let className = '';
            if (displayMode === 'day') {
                const { today_vec } = getDayVectors();
                className = today_vec[i] || '';
            } else {
                const { today_vec } = getDayVectors();
                className = today_vec[i] || '';
            }
            if (remainingTime > CONFIG.notifications.endingTime &&
                (now - lastNotificationTime) >= regularInterval * 60 * 1000) {
                playNotification('regular', className);
                lastNotificationTime = now;
            }
            if (remainingTime <= CONFIG.notifications.endingTime && !endingNotified) {
                playNotification('ending', className);
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