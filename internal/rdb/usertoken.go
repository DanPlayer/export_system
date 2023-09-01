package rdb

import (
	"context"
	"errors"
	"export_system/internal/db"
	"time"
)

const UserTokenKey = "UserToken:"

type UserToken struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

func (rs *UserToken) Get() (string, error) {
	if rs.UserID == "" {
		return "", errors.New("token缓存的必要属性没有设置")
	}
	return db.RedisClient.Get(context.Background(), rs.getKey())
}

func (rs *UserToken) Set() error {
	if rs.UserID == "" || len(rs.Token) <= 0 {
		return errors.New("token缓存的必要属性没有设置")
	}
	return db.RedisClient.Set(context.Background(), rs.getKey(), rs.Token, CacheTime*time.Second)
}

func (rs *UserToken) Del() error {
	if rs.UserID == "" {
		return errors.New("token缓存的必要属性没有设置")
	}
	return db.RedisClient.Del(context.Background(), rs.getKey())
}

func (rs *UserToken) getKey() string {
	return UserTokenKey + rs.UserID
}
