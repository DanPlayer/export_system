package utils

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
)

// UnicodeEmojiCode 表情转换
func UnicodeEmojiCode(s string) string {
	encodes := utf16.Encode([]rune(s))
	ret := ""
	for _, enc := range encodes {
		//ascii可见字符范围0-9a-zA-Z
		if (enc >= 0x30 && enc <= 0x39) || (enc >= 0x41 && enc <= 0x5A) || (enc >= 0x61 && enc <= 0x7A) {
			ret = ret + string(rune(enc))
		} else { //其他都使用utf16编码
			if encStr := strconv.FormatUint(uint64(enc), 16); encStr != "" {
				//补齐4位长度
				diff := 4 - len(encStr)
				for i := 0; i < diff; i++ {
					encStr = "0" + encStr
				}
				ret = ret + `u` + encStr
			}
		}
	}
	return ret
}

// UnicodeEmojiDecode 表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}
