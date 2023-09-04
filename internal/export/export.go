package export

import (
	"export_system/internal/db"
	"export_system/internal/rabbitmq"
	"export_system/pkg/exportcenter"
)

var Client = NewClient()

func NewClient() *exportcenter.ExportCenter {
	// 开启导出中心
	center, err := exportcenter.NewClient(exportcenter.Options{
		Db:           db.MasterClient,
		QueuePrefix:  "yb_",
		Queue:        rabbitmq.Client,
		SheetMaxRows: 500000,
		PoolMax:      2,
		GoroutineMax: 30,
	})
	if err != nil {
		return nil
	}

	return center
}
