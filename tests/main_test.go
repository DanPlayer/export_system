package tests

import (
	"encoding/json"
	"export_system/internal/domain/model"
	"export_system/internal/domain/service"
	"export_system/utils"
	"fmt"
	"github.com/sony/sonyflake"
	"log"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
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

func BenchmarkExportExcel(b *testing.B) {
	count := int64(269276)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Printf("write data start time:%s\n", time.Now())
		var ids []uint
		for i := 0; i < 5; i++ {
			id := CreateTestTask(fmt.Sprintf("test_listing_desc_%d", i), count)
			ids = append(ids, id)
		}
		fmt.Printf("write data end time:%s\n", time.Now())

		fmt.Printf("cusume data start time:%s\n", time.Now())
		wg := sync.WaitGroup{}
		for _, id := range ids {
			wg.Add(1)
			go func(id uint) {
				Export(int64(id))
				wg.Done()
			}(id)
		}
		wg.Wait()
		fmt.Printf("cusume data end time:%s\n", time.Now())
	}
}

func CreateTestTask(taskKey string, count int64) uint {
	id, keys, err := service.CreateExportTask(
		taskKey,
		"test_name",
		"test_file",
		"测试性能使用",
		"本地处理的数据",
		"当作测试性能用例",
		"xlsx",
		count,
		[]string{
			"id",
			"title",
			"keywords",
			"description",
			"bullet_point1",
			"bullet_point2",
			"bullet_point3",
			"bullet_point4",
			"bullet_point5",
		},
	)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	wg := sync.WaitGroup{}

	// 查询出50w左右的数据用于导入
	listingDescModel := model.ListingDesc{}
	for _, key := range keys {
		wg.Add(1)
		go func(key string) {
			limit := 1000
			did := 0
			for i := 0; i < 50; i++ {
				descs, err := listingDescModel.ListRangeByID(did, limit)
				if err != nil {
					break
				}
				if len(descs) == 0 {
					break
				}
				did = int(descs[len(descs)-1].ID)
				for _, desc := range descs {
					s := StructToSlice(desc)
					bytes, err := json.Marshal(s)
					if err != nil {
						continue
					}
					err = service.PushExportData(key, []string{
						string(bytes),
					})
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			}
			wg.Done()
		}(key)
	}
	wg.Wait()

	return id
}

func TestExport(t *testing.T) {
	id := int64(96)
	Export(id)
}

func Export(id int64) {
	service.StartTask(id)

	er := service.ExportToExcel(id, fmt.Sprintf("./test_listing_desc_%d.xlsx", id))
	if er != nil {
		fmt.Println(er)
		return
	}
}

func StructToSlice(f model.ListingDesc) []string {
	v := reflect.ValueOf(f)
	ss := make([]string, v.NumField())
	for i := range ss {
		ss[i] = fmt.Sprintf("%v", v.Field(i))
	}
	return ss
}

func TestSimpleTaskExport(t *testing.T) {
	id, keys, err := service.CreateExportTask(
		"test_mq",
		"test_name",
		"test_file",
		"测试使用",
		"本地处理的数据",
		"本地处理的数据",
		"xlsx",
		3,
		[]string{
			"header1",
			"header2",
			"header3",
		},
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, key := range keys {
		data := []string{
			"[\"get1\",\"get1\",\"get1\"]",
			"[\"get1\",\"get1\",\"get1\"]",
			"[\"get1\",\"get1\",\"get1\"]",
		}
		err := service.PushExportData(key, data)
		if err != nil {
			return
		}
	}

	service.StartTask(int64(id))

	er := service.ExportToExcel(int64(id), "./test.xlsx")
	if er != nil {
		return
	}
}
