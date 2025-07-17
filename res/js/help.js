/** 告示牌优化版 */

const sourceStr = typeof sth === 'string' ? sth : '';

const helpElem = document.getElementById("help");
if (helpElem) {
    // 支持多行换行显示
    helpElem.innerHTML = sourceStr.replace(/\n/g, '<br>');
}