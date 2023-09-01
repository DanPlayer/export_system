package utils

import (
	"database/sql"
	"time"
)

var (
	ASTM, _ = time.LoadLocation("Asia/Shanghai")
)

func Now() time.Time {
	return time.Now().In(ASTM)
}

func Parse(layout, value string) (time.Time, error) {
	return ParseInLocation(layout, value, ASTM)
}

func ParseInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(layout, value, loc)
}

func SqlNullTimeFormat(nullTime sql.NullTime, layout string) string {
	if nullTime.Valid {
		return nullTime.Time.Format(layout)
	}
	return ""
}

func SqlNullTimeUnix(nullTime sql.NullTime) int64 {
	if nullTime.Valid {
		return nullTime.Time.Unix()
	}
	return 0
}

// TimeStart 获取指定时间的当天开始时间，即凌晨 0点0分0秒
func TimeStart(atime time.Time) time.Time {
	return time.Date(atime.Year(), atime.Month(), atime.Day(), 0, 0, 0, 0, ASTM)
}

// TimeEnd 获取指定时间当天的最晚时间，即 23点59分59秒
func TimeEnd(atime time.Time) time.Time {
	return time.Date(atime.Year(), atime.Month(), atime.Day(), 23, 59, 59, 0, ASTM)
}

// Age 根据出生日期，计算年龄
func Age(birthDay time.Time) int {
	// 计算评价年龄
	avgAge := Now().Year() - birthDay.Year()
	if Now().Month()-birthDay.Month() < 0 { // 月份未满一年
		avgAge -= 1
	}
	return avgAge
}
