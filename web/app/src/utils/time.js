import i18n from '@/i18n'

/**
 * Generates a human-readable relative time string (e.g., "2 hours ago")
 * @param {string|Date} timestamp - The timestamp to convert
 * @returns {string} Relative time string
 */
export const generatePrettyTimeAgo = (timestamp) => {
  const t = i18n.global.t
  let differenceInMs = new Date().getTime() - new Date(timestamp).getTime();
  if (differenceInMs < 500) {
    return t('time.now');
  }
  if (differenceInMs > 3 * 86400000) { // If it was more than 3 days ago, we'll display the number of days ago
    let days = Math.round(differenceInMs / 86400000);
    return t('time.day', days);
  }
  if (differenceInMs > 3600000) { // If it was more than 1h ago, display the number of hours ago
    let hours = Math.round(differenceInMs / 3600000);
    return t('time.hour', hours);
  }
  if (differenceInMs > 60000) {
    let minutes = Math.round(differenceInMs / 60000);
    return t('time.minute', minutes);
  }
  let seconds = Math.round(differenceInMs / 1000);
  return t('time.second', seconds);
}

/**
 * Generates a pretty time difference string between two timestamps
 * @param {string|Date} start - Start timestamp
 * @param {string|Date} end - End timestamp
 * @returns {string} Time difference string
 */
export const generatePrettyTimeDifference = (start, end) => {
  const t = i18n.global.t
  const ms = new Date(start) - new Date(end)
  const seconds = Math.floor(ms / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)

  if (hours > 0) {
    const remainingMinutes = minutes % 60
    const hoursText = t('time.durationHour', hours)
    if (remainingMinutes > 0) {
      return hoursText + ' ' + t('time.durationMinute', remainingMinutes)
    }
    return hoursText
  } else if (minutes > 0) {
    const remainingSeconds = seconds % 60
    const minutesText = t('time.durationMinute', minutes)
    if (remainingSeconds > 0) {
      return minutesText + ' ' + t('time.durationSecond', remainingSeconds)
    }
    return minutesText
  } else {
    return t('time.durationSecond', seconds)
  }
}

/**
 * Formats a timestamp into YYYY-MM-DD HH:mm:ss format
 * @param {string|Date} timestamp - The timestamp to format
 * @returns {string} Formatted timestamp
 */
export const prettifyTimestamp = (timestamp) => {
  let date = new Date(timestamp);
  let YYYY = date.getFullYear();
  let MM = ((date.getMonth() + 1) < 10 ? "0" : "") + "" + (date.getMonth() + 1);
  let DD = ((date.getDate()) < 10 ? "0" : "") + "" + (date.getDate());
  let hh = ((date.getHours()) < 10 ? "0" : "") + "" + (date.getHours());
  let mm = ((date.getMinutes()) < 10 ? "0" : "") + "" + (date.getMinutes());
  let ss = ((date.getSeconds()) < 10 ? "0" : "") + "" + (date.getSeconds());
  return YYYY + "-" + MM + "-" + DD + " " + hh + ":" + mm + ":" + ss;
}