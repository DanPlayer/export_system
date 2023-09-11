package v1

import (
	"export_system/internal/domain/api/v1/auth"
	"export_system/internal/domain/api/v1/export"
	"export_system/internal/domain/api/v1/user"
	"export_system/internal/domain/common"
	"export_system/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Setup 注册业务模块
func Setup() common.ModuleOption {
	return common.ModuleOption{
		Name: "v1",
		ChildList: []common.ModuleChild{
			{
				Route:   "/user/login",
				Method:  "POST",
				Handles: []gin.HandlerFunc{user.Login},
			},
			{
				Route:   "/user/sms/login",
				Method:  "POST",
				Handles: []gin.HandlerFunc{user.SmsLogin},
			},
			{
				Route:   "/user/check/exist",
				Method:  "GET",
				Handles: []gin.HandlerFunc{user.CheckPhoneUser},
			},
			{
				Route:   "/user/info",
				Method:  "GET",
				Handles: []gin.HandlerFunc{middleware.Auth(), user.Info},
			},
			{
				Route:   "/user/check/profile",
				Method:  "GET",
				Handles: []gin.HandlerFunc{middleware.Auth(), user.CheckUserProfile},
			},
			{
				Route:   "/user/write/off",
				Method:  "POST",
				Handles: []gin.HandlerFunc{middleware.Auth(), user.WriteOffUser},
			},
			// 授权中心
			{
				Route:   "/auth/token",
				Method:  "POST",
				Handles: []gin.HandlerFunc{auth.AccessToken},
			},
			{
				Route:   "/auth/token/refresh",
				Method:  "POST",
				Handles: []gin.HandlerFunc{auth.RefreshAccessToken},
			},
			// 导出系统
			{
				Route:   "/export/task/create",
				Method:  "POST",
				Handles: []gin.HandlerFunc{middleware.AccessAuth(), export.CreateTask},
			},
			{
				Route:   "/export/task/data/push",
				Method:  "POST",
				Handles: []gin.HandlerFunc{middleware.AccessAuth(), export.PushExportData},
			},
		},
	}
}
