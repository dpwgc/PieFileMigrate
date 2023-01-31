package util

import "time"

var cstSh, _ = time.LoadLocation("Asia/Shanghai")

// GetLocalDate 年月日期获取
func GetLocalDate() string {
	return time.Now().In(cstSh).Local().Format("2006-01-02")
}

// GetLocalDateTime 秒级日期获取
func GetLocalDateTime() string {
	return time.Now().In(cstSh).Local().Format("2006-01-02 15:04:05")
}
