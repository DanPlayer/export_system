package auth

import (
	"export_system/internal/domain/rtn"
	"export_system/internal/domain/service"
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
// @Param body body AccessTokenRequest true "参数"
// @Success 200 {object} AccessTokenResponse
// @Router /v1/auth/token [post]
func AccessToken(c *gin.Context) {
	var req AccessTokenRequest
	_ = c.ShouldBindJSON(&req)
	if err := validate.ParseStruct(req); err != nil {
		utils.OutErrorJson(c, err)
		return
	}

	pass, _ := service.CheckAccessAuth(req.AccessKey, req.AccessSecret)
	if pass {
		// 生成token
		token, expire, err := service.GenerateAccessToken()
		if err != nil {
			utils.OutRtnErrorJson(c, rtn.AuthGenerateAccessTokenError)
			return
		}
		utils.OutJson(c, AccessTokenResponse{
			Token:      token,
			ExpireTime: expire,
		})
		return
	} else {
		utils.OutRtnErrorJson(c, rtn.AuthTokenError)
		return
	}
}

type RefreshAccessTokenRequest struct {
	Token string `json:"token" validate:"required"` // 已获取的token
}

type RefreshAccessTokenResponse struct {
	Token      string `json:"token"`
	ExpireTime int64  `json:"expire_time"`
}

// RefreshAccessToken
// @Summary 刷新AccessToken
// @Description 刷新AccessToken
// @Tags 授权中心
// @Accept json
// @Produce json
// @Param body body RefreshAccessTokenRequest true "参数"
// @Success 200 {object} RefreshAccessTokenResponse
// @Router /v1/auth/token/refresh [post]
func RefreshAccessToken(c *gin.Context) {
	var req RefreshAccessTokenRequest
	_ = c.ShouldBindJSON(&req)
	if err := validate.ParseStruct(req); err != nil {
		utils.OutErrorJson(c, err)
		return
	}

	// 刷新token
	token, expire, err := service.RefreshAccessToken(req.Token)
	if err != nil {
		utils.OutRtnErrorJson(c, err)
		return
	}
	utils.OutJson(c, AccessTokenResponse{
		Token:      token,
		ExpireTime: expire,
	})
	return
}
