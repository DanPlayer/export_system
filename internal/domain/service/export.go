package service

import (
	"export_system/internal/export"
	"export_system/pkg/exportcenter"
	"export_system/pkg/rtnerr"
)

// CreateExportTask 创建导出任务
func CreateExportTask(key, name, fileName, description, source, destination, format string, count int64, header []string) (id uint, err rtnerr.RtnError) {
	id, er := export.Client.CreateTask(key, name, description, source, destination, format, count, exportcenter.ExportOptions{
		FileName: fileName,
		Header:   header,
	})
	if er != nil {
		return 0, rtnerr.New(er)
	}
	return id, nil
}

// PushExportData 推送导出数据到队列
func PushExportData(key string, data []string) rtnerr.RtnError {
	for _, datum := range data {
		err := export.Client.PushData(key, datum)
		if err != nil {
			return rtnerr.New(err)
		}
	}

	return nil
}

// ExportToExcelCSV 导出成excel表格，格式csv
func ExportToExcelCSV(id int64, filePath string) error {
	err := export.Client.ExportToExcelCSV(id, filePath)
	if err != nil {
		return nil
	}
	return err
}
