package ip2region

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

const xdbPath = "./ip2region.xdb"

var Searcher = New()

func New() *xdb.Searcher {
	cBuff, err := xdb.LoadContentFromFile(xdbPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load content from `%s`: %s\n", xdbPath, err))
	}

	// 2、用全局的 cBuff 创建完全基于内存的查询对象。
	searcher, err := xdb.NewWithBuffer(cBuff)
	if err != nil {
		panic(fmt.Sprintf("failed to create searcher with content: %s\n", err))
	}
	return searcher
}
