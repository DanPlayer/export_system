package export

import (
	"export_system/internal/domain/service"
	"export_system/pkg/validate"
	"export_system/utils"
	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Key         string   `json:"key"`         // 队列KEY
	Name        string   `json:"name"`        // 任务名称
	FileName    string   `json:"file_name"`   // 文件名称
	Description string   `json:"description"` // 描述
	Source      string   `json:"source"`      // 来源
	Destination string   `json:"destination"` // 目的
	Format      string   `json:"format"`      // 导出文件类型，暂时只支持CSV
	Count       int64    `json:"count"`       // 导出数据的总行数
	Header      []string `json:"header"`      // 表格标题
}

type CreateTaskResponse struct {
	ID   uint     `json:"id"`   // 任务ID
	Keys []string `json:"keys"` // 所有队列key值，用于推送数据到不同队列，队列key由数据量生成
}

// CreateTask
// @Summary 创建导出任务
// @Description 创建导出任务
// @Tags 导出
// @Accept json
// @Produce json
// @Param body body CreateTaskRequest true "创建导出任务参数"
// @Success 200 {object} CreateTaskResponse
// @Router /v1/export/task/create [post]
func CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	_ = c.ShouldBindJSON(&req)
	if err := validate.ParseStruct(req); err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	id, keys, err := service.CreateExportTask(req.Key, req.Name, req.FileName, req.Description, req.Source, req.Destination, req.Format, req.Count, req.Header)
	if err != nil {
		utils.OutErrorJson(c, err)
		return
	}

	utils.OutJson(c, CreateTaskResponse{
		ID:   id,
		Keys: keys,
	})
}

type PushExportDataRequest struct {
	Key  string   `json:"key"`  // 队列KEY
	Data []string `json:"data"` // 数据，json字符串组成的数组
}

// PushExportData
// @Summary 推送导出数据
// @Description 推送导出数据
// @Tags 导出
// @Accept json
// @Produce json
// @Param body body CreateTaskRequest true "创建导出任务参数"
// @Success 200 {object} CreateTaskResponse
// @Router /v1/export/task/data/push [post]
func PushExportData(c *gin.Context) {
	var req PushExportDataRequest
	_ = c.ShouldBindJSON(&req)
	if err := validate.ParseStruct(req); err != nil {
		utils.OutErrorJson(c, err)
		return
	}
	err := service.PushExportData(req.Key, req.Data)
	if err != nil {
		utils.OutErrorJson(c, err)
		return
	}

	utils.OutJsonOk(c, "推入数据成功")
}
