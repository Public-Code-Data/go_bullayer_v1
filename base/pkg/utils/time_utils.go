package utils

import "time"

// Now 获取当前时间
func Now() time.Time {
	return time.Now()
}

// FormatTime 格式化时间
// t: 时间对象
// layout: 格式化模板，如 "2006-01-02 15:04:05"
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// FormatDateTime 格式化为日期时间字符串
// 格式: 2006-01-02 15:04:05
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDate 格式化为日期字符串
// 格式: 2006-01-02
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// ParseTime 解析时间字符串
// timeStr: 时间字符串
// layout: 时间格式模板
func ParseTime(timeStr, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}
