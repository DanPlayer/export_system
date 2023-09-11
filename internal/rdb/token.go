package rdb

import (
	"context"
	"encoding/json"
	"errors"
	"export_system/internal/db"
	"time"
)

const TokenKey = "token:"

type Token struct {
	UserID   string `json:"user_id"`
	Phone    string `json:"phone"`
	Token    string `json:"token"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
}

func (rs *Token) Get() (string, error) {
	if rs.Token == "" {
		return "", errors.New("token缓存的必要属性没有设置")
	}
	return db.RedisClient.Get(context.Background(), rs.getKey())
}

func (rs *Token) Set() error {
	if rs.UserID == "" || len(rs.Token) <= 0 {
		return errors.New("token缓存的必要属性没有设置")
	}
	bytes, e := json.Marshal(rs)
	if e != nil {
		return e
	}
	return db.RedisClient.Set(context.Background(), rs.getKey(), string(bytes), CacheTime*time.Second)
}

func (rs *Token) Del() error {
	if rs.Token == "" {
		return errors.New("token缓存的必要属性没有设置")
	}
	return db.RedisClient.Del(context.Background(), rs.getKey())
}

func (rs *Token) getKey() string {
	return TokenKey + rs.Token
}
