package utils

import (
	"bytes"
	"encoding/json"
	"export_system/internal/domain/rtn"
	"export_system/pkg/rtnerr"
	"github.com/gin-gonic/gin"
)

const (
	DATE_FORMAT             = "2006-01-02"
	DATE_FORMAT_SLASH       = "2006/01/02"
	DATEMONTH_FORMAT        = "2006-01"
	DATETIME_FORMAT         = "2006-01-02 15:04:05"
	DATETIMEWITHZONE_FORMAT = "2006-01-02 15:04:05 -07"
)

func OutJsonOk(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"rtn": rtn.OkRtn,
		"msg": msg,
	})
}

func OutJson(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"rtn":  rtn.OkRtn,
		"msg":  "成功",
		"data": data,
	})
}

func OutParamErrorJson(c *gin.Context) {
	c.JSON(200, gin.H{
		"rtn": rtn.ParamError,
		"msg": "参数错误",
	})
}

func OutErrorJson(c *gin.Context, err error) {
	c.JSON(200, gin.H{
		"rtn": rtn.CustomerError,
		"msg": err.Error(),
	})
}

func OutRtnErrorJson(c *gin.Context, r rtnerr.RtnError) {
	c.JSON(200, gin.H{
		"rtn": r.Rtn(),
		"msg": r.Error(),
	})
}

func OutBalanceNotEnoughJson(c *gin.Context, err error) {
	c.JSON(200, gin.H{
		"rtn": rtn.UserBalanceLowError,
		"msg": err.Error(),
	})
}

func OutUserWriteOff(c *gin.Context) {
	c.JSON(200, gin.H{
		"rtn": rtn.UserWriteOffError,
		"msg": "用户不存在",
	})
}

// OutAuthNeedError 需要登录
func OutAuthNeedError(c *gin.Context) {
	c.JSON(200, gin.H{
		"rtn": rtn.AuthNeedLoginError,
		"msg": "该功能需要登录!",
	})
}

// OutAuthOutdatedError 输出过期错误
func OutAuthOutdatedError(c *gin.Context) {
	c.JSON(200, gin.H{
		"rtn": rtn.AuthOutDateError,
		"msg": "登录态已过期,请重新登录!",
	})
}

func OutErrorMsg(c *gin.Context, err string) {
	c.JSON(200, gin.H{
		"rtn": rtn.CustomerError,
		"msg": err,
	})
}

// ValidateJson 监测json数据正确性
func ValidateJson(jsonString string) (right bool, data interface{}) {
	decoder := json.NewDecoder(bytes.NewReader([]byte(jsonString)))
	decoder.UseNumber()
	err := decoder.Decode(&data)
	if err != nil {
		return false, data
	}
	return true, data
}
