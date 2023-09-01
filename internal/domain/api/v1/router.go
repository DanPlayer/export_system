package v1

import (
	"export_system/internal/domain/api/v1/sms"
	"export_system/internal/domain/api/v1/upload"
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
				Route:   "/sms/send/login/code",
				Method:  "POST",
				Handles: []gin.HandlerFunc{sms.SendLoginSms},
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
			{ // 上传文件
				Route:   "/upload/qiniu/token",
				Method:  "GET",
				Handles: []gin.HandlerFunc{middleware.Auth(), upload.QiniuUploadToken},
			},
			{
				Route:   "/upload/open/excel",
				Method:  "GET",
				Handles: []gin.HandlerFunc{upload.OpenExcel},
			},
		},
	}
}
