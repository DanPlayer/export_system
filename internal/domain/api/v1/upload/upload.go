package upload

import (
	"export_system/internal/domain/qiniu"
	"export_system/internal/middleware"
	"export_system/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"time"
)

// QiniuUploadToken
// @Summary 获取上传凭证
// @Description 获取上传凭证
// @Tags qiniu
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200
// @Router /v1/upload/qiniu/token [get]
func QiniuUploadToken(c *gin.Context) {
	uid := middleware.GetLoginUserID(c)
	utils.OutJson(c, qiniu.MakeUpToken(uid))
}

// OpenExcel
// @Summary 打开Excel文档
// @Description 打开Excel文档
// @Tags upload
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200
// @Router /v1/upload/open/excel [get]
func OpenExcel(c *gin.Context) {
	fmt.Println("start time:", time.Now().Unix())
	f, err := excelize.OpenFile("D:\\mnt\\pro_cate_report\\Eugee_US_2023-07-13.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("read file done time:", time.Now().Unix())
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err := f.Rows("Template")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("read rows done time:", time.Now().Unix())
	for rows.Next() {
		row, err := rows.Columns()
		fmt.Println("read row done time:", time.Now().Unix())
		if err != nil {
			fmt.Println(err)
		}
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
	if err = rows.Close(); err != nil {
		fmt.Println(err)
	}
}
