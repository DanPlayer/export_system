package pojo

import (
	"database/sql"
)

// TaskVo 任务视图
type TaskVo struct {
	Name          string       `json:"name,omitempty"`           // 任务名称
	Description   string       `json:"description,omitempty"`    // 任务描述
	Status        int          `json:"status,omitempty"`         // 任务状态
	ProgressRate  int          `json:"progress_rate,omitempty"`  // 任务进度
	StartTime     sql.NullTime `json:"start_time"`               // 任务开始时间
	EndTime       sql.NullTime `json:"end_time"`                 // 任务结束时间
	Source        string       `json:"source,omitempty"`         // 任务来源
	Destination   string       `json:"destination,omitempty"`    // 任务目标
	ExportFormat  string       `json:"export_format,omitempty"`  // 导出文件类型
	ExportOptions string       `json:"export_options,omitempty"` // 导出配置
	QueueKey      string       `json:"queue_key,omitempty"`      // 队列key
	CountNum      int64        `json:"count_num,omitempty"`      // 导出总数
	WriteNum      int64        `json:"write_num,omitempty"`      // 写入数量
	ErrNum        int64        `json:"err_num,omitempty"`        // 错误数量
	ErrLogUrl     string       `json:"err_log_url,omitempty"`    // 错误日志地址
	DownloadUrl   string       `json:"download_url,omitempty"`   // 下载地址
}
