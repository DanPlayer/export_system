package auth

import (
	"export_system/pkg/validate"
	"export_system/utils"
	"github.com/gin-gonic/gin"
)

type AccessTokenRequest struct {
	AccessKey    string `json:"access_key" validate:"required"`    // 密钥key
	AccessSecret string `json:"access_secret" validate:"required"` // 密钥
}

type AccessTokenResponse struct {
	Token      string `json:"token"`
	ExpireTime int64  `json:"expire_time"`
}

// AccessToken
// @Summary 获取AccessToken
// @Description 获取AccessToken
// @Tags 授权中心
// @Accept json
// @Produce json
// @Param body body AccessTokenRequest true "用户添加参数"
// @Success 200 {object} AccessTokenResponse
// @Router /v1/auth/token [post]
func AccessToken(c *gin.Context) {
	var req AccessTokenRequest
	_ = c.ShouldBindJSON(&req)
	if err := validate.ParseStruct(req); err != nil {
		utils.OutErrorJson(c, err)
		return
	}

}
