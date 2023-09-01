package docs

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "export_system/docs"
	"export_system/internal/domain/common"
)

// Setup 注册文档模块
func Setup() common.ModuleOption {
	return common.ModuleOption{
		Name: "docs",
		ChildList: []common.ModuleChild{
			{
				Route:   "/*any",
				Method:  "GET",
				Handles: []gin.HandlerFunc{ginSwagger.WrapHandler(swaggerFiles.Handler)},
			},
		},
	}
}
