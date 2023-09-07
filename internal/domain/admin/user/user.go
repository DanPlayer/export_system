package user

import (
	"export_system/internal/domain/service"
	"export_system/pkg/validate"
	"export_system/utils"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Phone    string `json:"phone" validate:"required"`     // 手机号
	Password string `json:"password" validate:"required"`  // 密码
	NickName string `json:"nick_name" validate:"required"` // 昵称
	Avatar   string `json:"avatar"`                        // 头像
}

type RegisterResponse struct {
}

// Register
// @Summary 用户添加
// @Description 用户添加
// @Tags 后台管理-用户
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "用户添加参数"
// @Success 200 {object} RegisterResponse
// @Router /admin/user/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	_ = c.ShouldBindJSON(&req)
	if err := validate.ParseStruct(req); err != nil {
		utils.OutErrorJson(c, err)
		return
	}

	// 添加用户
	err := service.BackendAddUser(req.Phone, req.Password, req.NickName, req.Avatar)
	if err != nil {
		utils.OutRtnErrorJson(c, err)
		return
	}

	utils.OutJsonOk(c, "用户添加成功")
}
