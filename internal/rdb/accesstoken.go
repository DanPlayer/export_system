package rdb

import (
	"context"
	"encoding/json"
	"errors"
	"export_system/internal/db"
	"time"
)

const (
	AccessTokenKey        = "access_token:"
	AccessTokenExpireTime = 3600 * 24
)

type AccessToken struct {
	Token      string `json:"token"`
	ExpireTime int64  `json:"expire_time"`
}

func (rs *AccessToken) Get() (string, error) {
	if rs.Token == "" {
		return "", errors.New("token缓存的必要属性没有设置")
	}
	return db.RedisClient.Get(context.Background(), rs.getKey())
}

func (rs *AccessToken) Set() error {
	if len(rs.Token) <= 0 {
		return errors.New("token缓存的必要属性没有设置")
	}
	rs.ExpireTime = AccessTokenExpireTime
	bytes, e := json.Marshal(rs)
	if e != nil {
		return e
	}
	return db.RedisClient.Set(context.Background(), rs.getKey(), string(bytes), AccessTokenExpireTime*time.Second)
}

func (rs *AccessToken) Del() error {
	if rs.Token == "" {
		return errors.New("token缓存的必要属性没有设置")
	}
	return db.RedisClient.Del(context.Background(), rs.getKey())
}

func (rs *AccessToken) getKey() string {
	return AccessTokenKey + rs.Token
}
