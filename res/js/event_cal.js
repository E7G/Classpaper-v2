/* 倒计日优化版 */

const cal_source = events;
const cal_source_vec = cal_source.split(',');

const offset_basic = 2;
const now = new Date();
const evc = document.getElementById("evcal");
let cont = "";

for (let index = offset_basic; index < cal_source_vec.length; index += offset_basic) {
    const event = cal_source_vec[index];
    const event_date = cal_source_vec[index + 1];

    // 跳过无效数据
    if (!event || !event_date) continue;

    const dt = new Date(event_date);
    // 若日期无效则跳过
    if (isNaN(dt.getTime())) continue;

    const fin = (dt - now) / (1000 * 60 * 60 * 24);
    // 已过期事件不显示
    if (fin < 0) continue;

    const fweek = Math.floor(fin / 7);
    const fday = Math.ceil(fin % 7);

    cont += `<li><strong>${event}</strong> <mark>${Math.ceil(fin)}</mark> (${fweek}周${fday}天)</li>`;
}

evc.innerHTML = cont;
