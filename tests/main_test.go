package tests

import (
	"export_system/internal/domain/service"
	"export_system/utils"
	"fmt"
	"github.com/sony/sonyflake"
	"log"
	"strconv"
	"testing"
)

func main(m *testing.M) {
	fmt.Println("main test")
}

func TestGenNumber(t *testing.T) {
	flakeStartTime, err := utils.Parse(utils.DATE_FORMAT, "2022-07-01")
	if err != nil {
		log.Println(err)
		return
	}
	newSonyflake := sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: flakeStartTime,
	})

	outOrderSnUint, err := newSonyflake.NextID()
	if err != nil {
		return
	}
	number := strconv.FormatUint(outOrderSnUint, 32)
	fmt.Println(number)
	number = strconv.FormatUint(outOrderSnUint, 10)
	fmt.Println(number)
}

func TestGenUid(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(utils.GetUid())
	}
}

func TestTaskExport(t *testing.T) {
	id, keys, err := service.CreateExportTask(
		"test52000",
		"test_name",
		"test_file",
		"测试使用",
		"本地处理的数据",
		"当作测试用例",
		"xlsx",
		2,
		[]string{
			"header1",
			"header2",
			"header3",
		},
	)

	for _, key := range keys {
		err = service.PushExportData(key, []string{
			"[\"get1\",\"get1\",\"get1\"]",
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	er := service.ExportToExcel(int64(id), "./test.xlsx")
	if er != nil {
		fmt.Println(er)
		return
	}
}
