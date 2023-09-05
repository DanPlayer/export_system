package service

import (
	"export_system/internal/export"
	"export_system/internal/rabbitmq"
	"export_system/pkg/exportcenter"
	"export_system/pkg/rtnerr"
	"fmt"
)

// CreateExportTask 创建导出任务
func CreateExportTask(key, name, fileName, description, source, destination, format string, count int64, header []string) (id uint, keys []string, err rtnerr.RtnError) {
	id, keys, er := export.Client.CreateTask(key, name, description, source, destination, format, count, exportcenter.ExportOptions{
		FileName: fileName,
		Header:   header,
	})
	if er != nil {
		return 0, nil, rtnerr.New(er)
	}
	return id, keys, nil
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

// ExportToExcel 导出成excel表格，格式csv
func ExportToExcel(id int64, filePath string) error {
	err := export.Client.ExportToExcel(id, filePath, func(key string) error {
		// 重新开启消费者
		err := rabbitmq.Client.DeclareConsume(key)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
