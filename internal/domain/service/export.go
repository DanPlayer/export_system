package service

import (
	"export_system/internal/domain/model"
	"export_system/pkg/rtnerr"
)

const sheetMax = 500000 // sheet最大行数限制

// CreateExportTask 创建导出任务
func CreateExportTask(name, description, source, destination, format string, count int64, options model.ExportOptions) (id int64, err rtnerr.RtnError) {
	// 创建导出任务

	// 根据数据量，创建导出任务的数据队列

	return
}

// PushExportData 推送导出数据到队列
func PushExportData(id int64, data []string) rtnerr.RtnError {

	return nil
}

// ExportToExcelCSV 导出成excel表格，格式csv
func ExportToExcelCSV(id int64) {
	// 获取任务信息

	// 创建并发工作组，在工作组中使用协程处理数据写入

	// 拉取队列数据

	// 生成或者打开excel

	// 读取所有sheet

	// 判断sheet的数据量是否达到限制，达到限制则增加数据到下一张sheet，设置当前数据增加的sheet索引值

	// 增加数据到当前sheet并记录当前数据行索引，达到限制新增sheet，并重置当前sheet索引值

	// 记录数据进度，记录错误数据数

	// 任务进度完成（数据量达到总数包括错误数据），删除队列

	// 将文件上传至云端，记录下载地址

	// 删除本地文件
}
