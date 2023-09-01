package rdb

import (
	"context"
	"export_system/internal/db"
	"time"
)

const (
	LockKey = "lock"
	LockVal = "locked"
)

type Lock struct {
	Key    string
	Expire time.Duration
}

func (l *Lock) Lock() (bool, error) {
	return db.RedisClient.SetNx(context.Background(), l.getKey(), LockVal, l.Expire)
}

func (l *Lock) UnLock() error {
	return db.RedisClient.Del(context.Background(), l.getKey())
}

func (l *Lock) Locked() bool {
	get, _ := db.RedisClient.Get(context.Background(), l.getKey())
	if get != "" {
		return true
	}
	return false
}

func (l *Lock) getKey() string {
	return LockKey + "::" + l.Key
}
