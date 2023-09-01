package sonyflake

import (
	"github.com/sony/sonyflake"
	"log"
	"time"
)

var Client = New()
var (
	DateFormat = "2006-01-02"
	ASTM, _    = time.LoadLocation("Asia/Shanghai")
)

func New() *sonyflake.Sonyflake {
	flakeStartTime, err := Parse(DateFormat, "2022-07-01")
	if err != nil {
		log.Println(err)
		return nil
	}
	return sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: flakeStartTime,
	})
}

func Parse(layout, value string) (time.Time, error) {
	return ParseInLocation(layout, value, ASTM)
}

func ParseInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(layout, value, loc)
}
