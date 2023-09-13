package task

import (
	"export_system/internal/domain/pojo"
	"export_system/internal/domain/service"
	"export_system/pkg/validate"
	"export_system/utils"
	"github.com/gin-gonic/gin"
)

type PageListRequest struct {
	Page int `json:"page" validate:"required"` // 页码
	Size int `json:"size" validate:"required"` // 页面大小
}

type PageListResponse struct {
	List  []pojo.TaskVo `json:"list"`  // 列表
	Count int64         `json:"count"` // 总数
}

// PageList
// @Summary 任务列表
// @Description 任务列表
// @Tags 任务
// @Produce json
// @param query query PageListRequest true "参数"
// @Success 200 {object} PageListResponse
// @Router /v1/task/page/list [get]
func PageList(c *gin.Context) {
	var req PageListRequest
	_ = c.ShouldBindJSON(&req)
	if err := validate.ParseStruct(req); err != nil {
		utils.OutErrorJson(c, err)
		return
	}

	list, count, err := service.TaskList(req.Page, req.Size)
	if err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	utils.OutJson(c, PageListResponse{
		List:  list,
		Count: count,
	})
}
