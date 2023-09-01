package httppool

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	DefaultPoolSize            = 3
	DefaultMaxIdleConns        = 6000
	DefaultMaxIdleConnsPerHost = 6000
	DefaultWriteBufferSize     = 4 * 1024
	DefaultReadBufferSize      = 4 * 1024
	DefaultTimeout             = 20 * time.Second
)

var clients = make([]*http.Client, 0)

func init() {
	for i := 0; i < DefaultPoolSize; i++ {
		clients = append(clients, &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        DefaultMaxIdleConns,
				MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
				WriteBufferSize:     DefaultWriteBufferSize,
				ReadBufferSize:      DefaultReadBufferSize,
			},
			Timeout: DefaultTimeout,
		})
	}
}

// GetClient 获取一个可用的 Http Client
func GetClient() *http.Client {
	return clients[rand.Intn(DefaultPoolSize)]
}
