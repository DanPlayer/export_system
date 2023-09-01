package sms

import (
	"export_system/internal/domain/service"
	"export_system/pkg/validate"
	"export_system/utils"
	"github.com/gin-gonic/gin"
)

type SendLoginSmsRequest struct {
	Phone string `json:"phone" validate:"required"` // 手机号
}

// SendLoginSms
// @Summary 发送登录短信
// @Description 发送登录短信
// @Tags 短信
// @Accept json
// @Produce json
// @Param body body SendLoginSmsRequest true "生产场地订单参数"
// @Success 200
// @Router /v1/sms/send/login/code [post]
func SendLoginSms(c *gin.Context) {
	var req SendLoginSmsRequest
	_ = c.ShouldBindJSON(&req)
	if err := validate.ParseStruct(req); err != nil {
		utils.OutErrorJson(c, err)
		return
	}

	err := service.SendLoginSms(req.Phone)
	if err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	utils.OutJsonOk(c, "发送成功")
}

type SendUserRegisterForTeamSmsCodeRequest struct {
	Phone string `json:"phone" validate:"required"` // 手机号
}
