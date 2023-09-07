package user

import (
	"errors"
	"export_system/internal/domain/pojo"
	"export_system/internal/domain/service"
	"export_system/internal/middleware"
	"export_system/pkg/validate"
	"export_system/utils"
	"github.com/gin-gonic/gin"
)

type SmsLoginRequest struct {
	Phone string `json:"phone" validate:"required"` // 手机号
	Code  string `json:"code" validate:"required"`  // 验证码
}

type SmsLoginResponse struct {
	Token           string `json:"token"`           // 验证签名
	ProfileComplete bool   `json:"profileComplete"` // 基本资料是否填充完成
	UserID          string `json:"userID"`          // 用户ID
}

// SmsLogin
// @Summary 用户登录
// @Description 用户登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body SmsLoginRequest true "用户登录参数"
// @Success 200 {object} SmsLoginResponse
// @Router /v1/user/sms/login [post]
func SmsLogin(c *gin.Context) {
	var req SmsLoginRequest
	_ = c.ShouldBindJSON(&req)
	if err := validate.ParseStruct(req); err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	token, complete, userId, err := service.SmsLogin(req.Phone, req.Code)
	if err != nil {
		utils.OutErrorJson(c, err)
		return
	}

	utils.OutJson(c, SmsLoginResponse{
		Token:           token,
		ProfileComplete: complete,
		UserID:          userId,
	})
}

type UpdateUserInfoRequest struct {
	NickName string `json:"nickName"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
}

type CheckUserProfileResponse struct {
	Complete bool `json:"complete"` // 是否完成基本资料
}

// CheckUserProfile
// @Summary 检查用户是否完成基本资料
// @Description 检查用户是否完成基本资料
// @Tags 用户
// @Produce json
// @Success 200 {object} CheckUserProfileResponse
// @Router /v1/user/check/profile [get]
func CheckUserProfile(c *gin.Context) {
	userID := middleware.GetLoginUserID(c)
	pass := service.CheckUserProfile(userID)
	utils.OutJson(c, CheckUserProfileResponse{Complete: pass})
}

type CheckPhoneUserResponse struct {
	UserExist bool `json:"userExist"` // 用户是否存在
}

// CheckPhoneUser
// @Summary 检查手机用户是否存在
// @Description 检查手机用户是否存在
// @Tags 用户
// @Param phone query string true "手机号"
// @Produce json
// @Success 200 {object} CheckPhoneUserResponse
// @Router /v1/user/check/exist [get]
func CheckPhoneUser(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		utils.OutErrorJson(c, errors.New("手机号参数缺失"))
		return
	}

	exist, _ := service.CheckPhoneUser(phone)
	utils.OutJson(c, CheckPhoneUserResponse{
		UserExist: exist,
	})
}

type FindUserRequest struct {
	Phone string `form:"phone" json:"phone" validate:"required"` //用户手机号码
}

type FindUserResponse struct {
	List []pojo.UserInfo `json:"list"` // 用户信息
}

type InfoResponse struct {
	pojo.UserInfo
}

// Info
// @Summary 当前用户的信息
// @Description 当前用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200  {object} InfoResponse
// @Router /v1/user/info [get]
func Info(c *gin.Context) {
	userID := middleware.GetLoginUserID(c)
	userInfo, err := service.UserInfo(userID)
	if err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	utils.OutJson(c, InfoResponse{userInfo})
}

// WriteOffUser
// @Summary 注销用户
// @Description 注销用户
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200
// @Router /v1/user/write/off [post]
func WriteOffUser(c *gin.Context) {
	err := service.WriteOffUser(middleware.GetLoginUserID(c))
	if err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	utils.OutJsonOk(c, "注销用户成功")
}
