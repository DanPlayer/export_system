package model

import (
	"export_system/pkg/exportcenter"
	"gorm.io/gorm"
	"time"
)

// Task 任务表
// 用于记录所有的到处任务以及导出状态
type Task struct {
	gorm.Model
	Name          string    `gorm:"type:varchar(255);comment:'任务名称'"`
	Description   string    `gorm:"type:text;comment:'描述'"`
	Status        int       `gorm:"type:tinyint(1);default:1;comment:'状态 1-待处理、2-处理中、3-已完成、4-失败、5-任务废弃'"`
	ProgressRate  int       `gorm:"type:tinyint(3);default:0;comment:'任务进度1-100'"`
	StartTime     time.Time `gorm:"type:timestamp;not null;comment:'任务开始时间'"`
	EndTime       time.Time `gorm:"type:timestamp;not null;comment:'任务结束时间'"`
	Source        string    `gorm:"type:varchar(255);comment:'数据源，描述导出数据的来源'"`
	Destination   string    `gorm:"type:varchar(255);comment:'数据目标，描述导出数据的存储位置'"`
	ExportFormat  string    `gorm:"type:varchar(255);comment:'导出格式，如CSV、JSON、XML等'"`
	ExportOptions string    `gorm:"type:text;comment:'导出选项，可存储导出任务的配置信息（可选）'"`
	QueueKey      string    `gorm:"type:varchar(255);comment:'队列key'"`
	CountNum      int64     `gorm:"type:int(11);default:0;comment:'数据总数'"`
	ErrNum        int64     `gorm:"type:int(11);default:0;comment:'错误数据数'"`
	ErrLogUrl     string    `gorm:"type:text;comment:'错误日志地址'"`
	DownloadUrl   string    `gorm:"type:text;comment:'文件下载地址'"`
}

// ExportOptions 导出选项
type ExportOptions struct {
	FileName string   `json:"file_name"` // 文件名称
	Header   []string `json:"header"`    // 表头配置
}

type TaskStatus int

const (
	TaskStatusWait      TaskStatus = 1
	TaskStatusConsult   TaskStatus = 2
	TaskStatusCompleted TaskStatus = 3
	TaskStatusFail      TaskStatus = 3
	TaskStatusAbandon   TaskStatus = 4
)

func (m *Task) TableName() string {
	return "task"
}

func (m *Task) Create() error {
	return exportcenter.DbClient.Model(&m).Create(&m).Error
}

func (m *Task) FindByID(id int64) (info Task, err error) {
	err = exportcenter.DbClient.Model(&m).Where("id = ?", id).First(&info).Error
	return
}

func (m *Task) UpdateStatusByID(id int64, status TaskStatus) error {
	return exportcenter.DbClient.Model(&m).Where("id = ?", id).UpdateColumn("status", status).Error
}
