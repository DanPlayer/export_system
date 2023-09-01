package tests

import (
	"export_system/internal/domain/ip2region"
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

func TestIpAddress(t *testing.T) {
	region, err := ip2region.Searcher.SearchByStr("171.113.249.158")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(region)
}

func TestGenUid(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(utils.GetUid())
	}
}
