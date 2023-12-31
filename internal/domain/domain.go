package domain

import (
	"export_system/internal/domain/admin"
	v1 "export_system/internal/domain/api/v1"
	"export_system/internal/domain/common"
	"export_system/internal/domain/crontab"
	"export_system/internal/domain/docs"
	"export_system/internal/domain/tool"
)

// Registry 挂载业务模块
func Registry() common.ModuleOptionList {
	return common.ModuleOptionList{
		docs.Setup(),
		tool.Setup(),
		v1.Setup(),
		admin.Setup(),
		crontab.Setup(),
	}
}
