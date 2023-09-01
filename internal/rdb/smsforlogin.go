package rdb

import (
	"context"
	"errors"
	"export_system/internal/db"
	"time"
)

const SmsForLoginKey = "SmsForLogin:"
const SmsForLoginExpire = 60

type SmsForLogin struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func (rs *SmsForLogin) Get() (string, error) {
	if rs.Phone == "" {
		return "", errors.New("验证码缓存的必要属性没有设置")
	}
	return db.RedisClient.Get(context.Background(), rs.getKey())
}

func (rs *SmsForLogin) Set() error {
	if rs.Phone == "" || len(rs.Code) <= 0 {
		return errors.New("验证码缓存的必要属性没有设置")
	}
	return db.RedisClient.Set(context.Background(), rs.getKey(), rs.Code, SmsForLoginExpire*time.Second)
}

func (rs *SmsForLogin) Del() error {
	if rs.Phone == "" {
		return errors.New("验证码缓存的必要属性没有设置")
	}
	return db.RedisClient.Del(context.Background(), rs.getKey())
}

func (rs *SmsForLogin) getKey() string {
	return SmsForLoginKey + rs.Phone
}
