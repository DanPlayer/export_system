package service

import (
	"export_system/internal/domain/model"
	"export_system/internal/domain/pojo"
)

// TaskList 任务列表
func TaskList(page, size int) (list []pojo.TaskVo, count int64, err error) {
	taskModel := model.Task{}
	pageList, count, err := taskModel.PageList(page, size)
	if err != nil {
		return
	}
	for i := range pageList {
		item := pojo.TaskVo{
			Name:          pageList[i].Name,
			Description:   pageList[i].Description,
			Status:        pageList[i].Status,
			ProgressRate:  pageList[i].ProgressRate,
			StartTime:     pageList[i].StartTime,
			EndTime:       pageList[i].EndTime,
			Source:        pageList[i].Source,
			Destination:   pageList[i].Destination,
			ExportFormat:  pageList[i].ExportFormat,
			ExportOptions: pageList[i].ExportOptions,
			QueueKey:      pageList[i].QueueKey,
			CountNum:      pageList[i].CountNum,
			WriteNum:      pageList[i].WriteNum,
			ErrNum:        pageList[i].ErrNum,
			ErrLogUrl:     pageList[i].ErrLogUrl,
			DownloadUrl:   pageList[i].DownloadUrl,
		}
		list = append(list, item)
	}
	return
}

// TaskInfo 任务信息
func TaskInfo(id int64) (info model.Task, err error) {
	taskModel := model.Task{}
	info, err = taskModel.FindByID(id)
	if err != nil {
		return
	}
	return
}
