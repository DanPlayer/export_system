package tool

import (
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping
// @Summary 服务健康检查
// @Tags tool
// @Accept json
// @Produce json
// @Success 200
// @Router /tool/ping [get]
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

type RealIPSchema struct {
	Client string `json:"client"` // 客户端 IP
	Server string `json:"server"` // 服务端出口IP
}

// GetRealIP
// @Summary 获取当前服务真实出口IP
// @Tags tool
// @Accept json
// @Produce json
// @Success 200 {object} RealIPSchema
// @Router /tool/real-ip [get]
func GetRealIP(c *gin.Context) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
	localAddr, _ := conn.LocalAddr().(*net.UDPAddr)
	c.JSON(http.StatusOK, RealIPSchema{
		Client: c.ClientIP(),
		Server: localAddr.IP.String(),
	})
}
