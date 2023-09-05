package export

import (
	"export_system/internal/domain/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

// TaskExport 任务导出
func TaskExport(c *gin.Context) {
	err := service.ExportToExcel(1, "/")
	if err != nil {
		fmt.Println(err)
		return
	}
}
