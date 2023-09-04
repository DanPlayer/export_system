package crontab

import (
	"export_system/internal/domain/common"
	"export_system/internal/domain/crontab/export"
	"github.com/gin-gonic/gin"
)

// Setup 注册业务模块
func Setup() common.ModuleOption {
	return common.ModuleOption{
		Name: "crontab",
		ChildList: []common.ModuleChild{
			{
				Route:   "/export/task",
				Method:  "GET",
				Handles: []gin.HandlerFunc{export.TaskExport},
			},
		},
	}
}
