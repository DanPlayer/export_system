package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	NUMBER = "0123456789"
)

var (
	s = rand.New(rand.NewSource(time.Now().Unix()))
)

// PageSize 通用获取分页数据
func PageSize(c *gin.Context) (page, size int, err error) {
	// 分页
	p := c.Query("page")
	if len(p) > 0 {
		page, err = strconv.Atoi(p)
		if err != nil {
			return
		}
		if page < 1 {
			page = 1
		}
	} else {
		page = 1
	}

	// 分页大小
	s := c.Query("size")
	if len(s) > 0 {
		size, err = strconv.Atoi(s)
		if err != nil {
			return
		}
		if size < 1 {
			size = 1
		}
	} else {
		size = 10
	}
	return
}

// TempPicFileName 生成临时图片路径 供前端上传
func TempPicFileName(id, name string) string {
	now := time.Now().In(ASTM)
	format := now.Format("20060102")
	uuidString := uuid.NewV4().String()
	return "default/" + id + fmt.Sprintf("/%v/", format) + uuidString + path.Ext(name)
}

// TrimHtml 去掉html标签
func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	reTag, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	reQuote, _ := regexp.Compile("&nbsp;")

	src = reTag.ReplaceAllString(src, "$1")
	src = reQuote.ReplaceAllString(src, " ")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

// Substr 按字符截取字符传
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

// InArray 判断元素是否存在于数组中
func InArray(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}

// ToSlice 将传入slice的每个元素拿出来interface()化
func ToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

// CompareIdsSame 判断是否全ids匹配
func CompareIdsSame(sourceIds []int, targetIds []int) (same bool) {
	same = true
	for _, v := range sourceIds {
		if ok, _ := InArray(v, targetIds); !ok {
			same = false
		}
	}
	return
}

// Unique 切片去重
func Unique(s []string) []string {
	m := make(map[string]struct{}, 0)
	newS := make([]string, 0)
	for _, i2 := range s {
		if _, ok := m[i2]; !ok {
			newS = append(newS, i2)
			m[i2] = struct{}{}
		}
	}
	return newS
}

// Split 分割字符串为数组
func Split(text, sep string) []string {
	if text == "" {
		return []string{}
	}
	return strings.Split(text, sep)
}

// ConvertSeconds 秒数转时分
func ConvertSeconds(value int64) (hour, minute int) {
	hour = int(value / 3600)
	minute = int(value % 3600 / 60)
	return
}
