package mysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Config 数据配置信息
type Config struct {
	// 数据库，支持：MySQL、PostgreSQL, SQLServer等
	// 数据库DNS,，例如：user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	DNS string
	// 最大连接数
	MaxOpenConns int `database:"default:60"`
	// 最大空闲数
	MaxIdleConns int `database:"default:10"`
}

// New 初始化MySQL数据库连接
func New(c Config) *gorm.DB {
	mysqlDial := mysql.Open(c.DNS)
	conn, err := initConnection(mysqlDial, c)
	// 开启debug模式
	if err != nil {
		panic(fmt.Sprintf("MySQL初始化失败：%v", err.Error()))
	}
	return conn
}

func Open(dns string) gorm.Dialector {
	return mysql.Open(dns)
}

// 初始化数据库连接
func initConnection(dial gorm.Dialector, config Config) (db *gorm.DB, err error) {
	var (
		originDB    *gorm.DB
		originSqlDB *sql.DB
	)
	if originDB, err = gorm.Open(dial, &gorm.Config{
		// 设置数据库表名称为单数(User,复数Users末尾自动添加s)
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}); err != nil {
		return nil, err
	}

	if originSqlDB, err = originDB.DB(); err != nil {
		return nil, err
	}

	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = getConfigTagDefaultValue("MaxOpenConns", "database")
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = getConfigTagDefaultValue("MaxIdleConns", "database")
	}

	originSqlDB.SetMaxIdleConns(config.MaxIdleConns)
	// 避免并发太高导致连接mysql出现too many connections的错误
	originSqlDB.SetMaxOpenConns(config.MaxOpenConns)
	// 设置数据库闲链接超时时间
	originSqlDB.SetConnMaxLifetime(time.Second * 30)
	return originDB, nil
}

// 获取配置文件指定字段默认属性
func getConfigTagDefaultValue(name string, tag string) (value int) {
	openField, _ := reflect.TypeOf(Config{}).FieldByName(name)
	openReg := regexp.MustCompile(`default:(\d*)`)
	vList := openReg.FindStringSubmatch(openField.Tag.Get(tag))
	if len(vList) == 2 {
		value, _ = strconv.Atoi(vList[1])
	}
	return value
}
