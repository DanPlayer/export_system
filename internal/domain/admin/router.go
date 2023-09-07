package admin

import (
	"export_system/internal/domain/admin/user"
	"export_system/internal/domain/common"
	"github.com/gin-gonic/gin"
)

// Setup 注册业务模块
func Setup() common.ModuleOption {
	return common.ModuleOption{
		Name: "admin",
		ChildList: []common.ModuleChild{
			{
				Route:   "/user/register",
				Method:  "POST",
				Handles: []gin.HandlerFunc{user.Register},
			},
		},
	}
}
