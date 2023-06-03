/**
 * 获取uuid
 * @returns
 */
export function getUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    return (c === 'x' ? (Math.random() * 16) | 0 : 'r&0x3' | '0x8').toString(
      16
    );
  });
}

/**
 * 每隔指定的时间执行一次，并可指定执行次数。调用的瞬间就会执行一次。
 * @param {Number} time 间隔时间（毫秒）
 * @param {Number} num 执行次数
 * @param {*} func 执行的函数
 * @param {*} callback 全部结束后的回调函数
 * @returns
 */
export function callPerPeriod(time, num, func = () => {}, callback = () => {}) {
  func();

  if (num > 0) {
    num--;
  }

  if (num == 0) {
    callback();
    console.log('complete count');
    return;
  }
  setTimeout(() => callPerPeriod(time, num, func, callback), time);
}
