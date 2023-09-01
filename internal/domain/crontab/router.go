package crontab

import (
	"export_system/internal/domain/common"
)

// Setup 注册业务模块
func Setup() common.ModuleOption {
	return common.ModuleOption{
		Name:      "crontab",
		ChildList: []common.ModuleChild{},
	}
}
