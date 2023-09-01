package utils

import (
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"os"
	"time"
)

var hostname = "unknow"
var logger *producer.Producer = nil

func init() {
	logger := makeProducer()
	logger.Start()

	name, err := os.Hostname()
	if err == nil {
		hostname = name
	}
}

func makeProducer() *producer.Producer {
	cfg := producer.GetDefaultProducerConfig()
	cfg.Endpoint = "cn-hangzhou.log.aliyuncs.com"
	cfg.AccessKeyID = "LTAInAUi1WTgkUhb"
	cfg.AccessKeySecret = "Ky9UEeT3Aia8Kw2IGjwCZ1YCjO7cWh"
	return producer.InitProducer(cfg)
}

func Info(item map[string]string) error {
	logs := producer.GenerateLog(uint32(time.Now().Unix()), item)
	return logger.SendLog("puahome-web", "app_server", "icitysecret_mg", hostname, logs)
}

func LogRequest(item map[string]string) error {
	logs := producer.GenerateLog(uint32(time.Now().Unix()), item)
	return logger.SendLog("puahome-web", "app_server", "icitysecret_mg", hostname, logs)
}
