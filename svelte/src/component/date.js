/**
 * @param {Date} d
 * @returns {string}
 */
export function dateToYYYYMM(d) {
  const yyyy = `${d.getFullYear()}`.padStart(4, "0");
  const mm = `${d.getMonth()}`.padStart(2, "0");
  return `${yyyy}${mm}`;
}

/**
 * @param {Date} d
 * @returns {string}
 */
export function dateToYYYYMMDDhhmm(d) {
  const YYYY = `${d.getFullYear()}`.padStart(4, "0");
  const MM = `${d.getMonth() + 1}`.padStart(2, "0");
  const DD = `${d.getDate()}`.padStart(2, "0");
  const hh = `${d.getHours()}`.padStart(2, "0");
  const mm = `${d.getMinutes()}`.padStart(2, "0");
  return `${YYYY}${MM}${DD}${hh}${mm}`;
}

/**
 * https://developer.mozilla.org/ja/docs/Web/JavaScript/Reference/Global_Objects/Date/Date#%E6%A7%8B%E6%96%87
 * @param yyyymm {string}
 * @returns {Date}
 */
export function yyyymmToDate(yyyymm) {
  const year = parseInt(yyyymm.slice(0, 4), 10);
  const monthIndex = parseInt(yyyymm.slice(4, 6), 10) - 1;
  return new Date(year, monthIndex);
}
