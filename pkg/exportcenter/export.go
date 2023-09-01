package exportcenter

import (
	"errors"
	"export_system/pkg/exportcenter/model"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/panjf2000/ants/v2"
	"gorm.io/gorm"
	"math"
	"sync"
	"sync/atomic"
)

var DbClient *gorm.DB

type ExportCenter struct {
	Db           *gorm.DB
	Queue        Queue
	queuePrefix  string
	sheetMaxRows int64
	poolMax      int
	goroutineMax int
}

// Options 配置
type Options struct {
	Db           gorm.DB // gorm实例
	QueuePrefix  string  // 队列前缀
	Queue        Queue   // 队列配置（必须配置）
	SheetMaxRows int64   // 数据表最大行数
	PoolMax      int     // 协程池最大数量
	GoroutineMax int     // 协程最大数量
}

// Queue 队列
type Queue interface {
	Create(key string) error            // 创建队列
	Pop(key string) string              // 拉取数据
	Push(key string, data string) error // 推送数据
	Del(key string) error               // 删除队列
}

func NewClient(options Options) (*ExportCenter, error) {
	if options.SheetMaxRows == 0 {
		return nil, errors.New("SheetMaxRows数据表最大行数必须配置大于0")
	}
	if options.PoolMax <= 0 {
		options.PoolMax = 1
	}
	if options.GoroutineMax <= 0 {
		options.GoroutineMax = 1
	}

	DbClient = &options.Db

	// 自动创建任务表
	err := DbClient.AutoMigrate(&model.Task{})
	if err != nil {
		return nil, err
	}

	return &ExportCenter{
		Db:    &options.Db,
		Queue: options.Queue,
	}, nil
}

// CreateTask 创建导出任务
func (ec *ExportCenter) CreateTask(key, name, description, source, destination, format string, count int64, options model.ExportOptions) error {
	marshal, err := json.Marshal(options)
	if err != nil {
		return err
	}

	// 创建导出任务
	task := model.Task{
		Name:          name,
		Description:   description,
		Status:        1,
		ProgressRate:  0,
		Source:        source,
		Destination:   destination,
		ExportFormat:  format,
		ExportOptions: string(marshal),
		QueueKey:      key,
		CountNum:      count,
	}
	err = task.Create()
	if err != nil {
		return err
	}

	return nil
}

// PushData 推送导出数据到队列
func (ec *ExportCenter) PushData(key string, data string) error {
	return ec.Queue.Push(key, data)
}

// PopData 拉取队列数据
func (ec *ExportCenter) PopData(key string) string {
	return ec.Queue.Pop(key)
}

func (ec *ExportCenter) GetTask(id int64) (info model.Task, err error) {

	return
}

// ExportToExcelCSV 导出成excel表格，格式csv
func (ec *ExportCenter) ExportToExcelCSV(id int64) (err error) {
	// 获取任务信息
	task, err := ec.GetTask(id)
	if err != nil {
		return err
	}

	// 根据数据量，创建导出任务的数据队列
	sheetCount := int(math.Ceil(float64(task.CountNum / ec.sheetMaxRows)))

	// 生成或者打开excel

	for i := 1; i <= sheetCount; i++ {
		queueKey := ""
		if ec.queuePrefix != "" {
			queueKey = fmt.Sprintf("%s_%s_sheet%d", ec.queuePrefix, task.QueueKey, i)
		} else {
			queueKey = fmt.Sprintf("%s_sheet%d", task.QueueKey, i)
		}

		err = ec.Queue.Create(queueKey)
		if err != nil {
			return err
		}

		// 创建sheet

	}

	// 判断sheet的数据量是否达到限制，达到限制则增加数据到下一张sheet，设置当前数据增加的sheet索引值
	sheetIndex := int32(0)
	count := int64(0)
	rowCount := int64(0)
	errRowCount := int64(0)

	// 创建并发工作组，在工作组中使用协程处理数据写入，单个协程会有一个小时的过期时间，一个小时内未完成单表设置的最大数量就会任务失败
	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(ec.poolMax, func(key interface{}) {
		queueKey := ""
		if ec.queuePrefix != "" {
			queueKey = fmt.Sprintf("%s_%s_sheet%d", ec.queuePrefix, key.(string), sheetIndex)
		} else {
			queueKey = fmt.Sprintf("%s_sheet%d", key.(string), sheetIndex)
		}

		// 拉取队列数据
		for {
			// 增加数据到当前sheet并记录当前数据行索引，达到限制新增sheet，并重置当前sheet索引值
			atomic.AddInt64(&count, 1) // 记录数据进度
			atomic.AddInt64(&rowCount, 1)
			if atomic.LoadInt64(&rowCount) > ec.sheetMaxRows {
				atomic.StoreInt64(&rowCount, 0)
				atomic.AddInt32(&sheetIndex, 1)
				break
			}

			data := ec.PopData(queueKey)
			if data == "" {
				// 记录错误数据数
				atomic.AddInt64(&errRowCount, 1)
				continue
			}

			// 写入excel文件

		}

		wg.Done()
	}, ants.WithExpiryDuration(3600))
	defer p.Release()
	// 提交协程任务
	for i := 0; i < ec.goroutineMax; i++ {
		wg.Add(1)
		_ = p.Invoke(task.QueueKey)
	}
	wg.Wait()

	// 任务进度完成（数据量达到总数包括错误数据），删除队列
	if count > task.CountNum {

	} else {
		// 任务失败，记录数据最后的id

	}

	// 将文件上传至云端，记录下载地址

	// 删除本地文件

	return
}
