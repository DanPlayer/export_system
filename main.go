package main

import (
	"context"
	"errors"
	"export_system/internal/config"
	"export_system/internal/domain"
	"export_system/internal/domain/qiniu"
	"export_system/internal/middleware"
	"export_system/internal/timewheel"
	gracefulExit "github.com/NICEXAI/graceful-exit"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// @title qingwu
// @version 1.0
// @in header
// @name Authorization
// @host localhost:8088
func main() {
	var app *gin.Engine

	if config.Config.Application.Mode == "debug" {
		app = gin.Default()
	} else {
		app = gin.New()
	}
	app.Use(middleware.CrossDomain())

	// 注册业务模块
	domainList := domain.Registry()
	for i := range domainList {
		option := domainList[i]
		group := app.Group(option.Name)
		{
			for j := range option.ChildList {
				child := option.ChildList[j]
				switch child.Method {
				case "GET":
					group.GET(child.Route, child.Handles...)
				case "POST":
					group.POST(child.Route, child.Handles...)
				case "PUT":
					group.PUT(child.Route, child.Handles...)
				case "DELETE":
					group.DELETE(child.Route, child.Handles...)
				}
			}
		}
	}

	// 开启七牛云服务
	_ = qiniu.InitQiniu()

	// 开启时间轮服务
	timewheel.Client.Start()
	defer timewheel.Client.Stop()

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(config.Config.Application.Port),
		Handler:      app,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				os.Exit(1)
			}
		}
	}()

	// 服务优雅退出
	graceful := gracefulExit.NewGracefulExit()
	graceful.RegistryHandle("exit", func() {
		// 停止server服务
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("server down error: %s", err.Error())
		}

		log.Println("service stop successfully")
	})

	graceful.Capture()
}
