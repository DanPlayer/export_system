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
	count := int64(500000)
	_, keys, err := service.CreateExportTask(
		"test_listing_desc",
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
		return
	}

	// 查询出50w左右的数据用于导入
	listingDescModel := model.ListingDesc{}
	limit := 10000
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
			err = service.PushExportData(keys[0], []string{
				string(bytes),
			})
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func TestExport(t *testing.T) {
	er := service.ExportToExcel(int64(68), "./test_listing_desc.xlsx")
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
