package db

import (
	"context"
	"export_system/internal/config"
	"export_system/pkg/mysql"
	"export_system/pkg/redis"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// RedisClient 初始化Redis客户端
var RedisClient = NewRedisClient()

// MasterClient 初始化MySQL客户端
var MasterClient = NewMasterClient()

// AmazonClient 易佰中台亚马逊数据库
var AmazonClient = NewAmazonClient()

// AutoMigrateApiServerTable 自动迁移数据表
func AutoMigrateApiServerTable(modes ...interface{}) {
	_ = MasterClient.AutoMigrate(modes...)
}

func NewMasterClient() *gorm.DB {
	client := mysql.New(mysql.Config{
		DNS:         config.Config.MySQL.Master,
		TablePrefix: config.Config.MySQL.MasterTablePrefix,
	})
	// 从库自动加载
	_ = client.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{mysql.Open(config.Config.MySQL.Slave)},
	}))
	return client
}

func NewAmazonClient() *gorm.DB {
	return mysql.New(mysql.Config{
		DNS: config.Config.MySQL.Amazon,
	})
}

func NewRedisClient() *redis.Redis {
	return redis.New(context.Background(), redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
}
