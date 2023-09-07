package rtn

import "export_system/pkg/rtnerr"

const (
	OkRtn                        rtnerr.RtnCode = 0     // 成功
	ParamError                   rtnerr.RtnCode = 10001 // 参数错误
	CustomerError                rtnerr.RtnCode = 10002 // 一般错误
	AuthNeedLoginError           rtnerr.RtnCode = 10003 // 需要权限
	AuthOutDateError             rtnerr.RtnCode = 10004 // 权限过期错误
	UserWriteOffError            rtnerr.RtnCode = 10005 // 用户已注销
	UserBalanceLowError          rtnerr.RtnCode = 10007 // 账户余额不足
	UserExistError               rtnerr.RtnCode = 10008 // 用户已存在
	UserAkExistError             rtnerr.RtnCode = 10009 // 用户AccessKey已存在
	AuthTokenError               rtnerr.RtnCode = 10010 // 没有授权中心的权限
	AuthGenerateAccessTokenError rtnerr.RtnCode = 10011 // 生成AccessToken失败
)

var rtnCodeMap = map[rtnerr.RtnCode]string{
	OkRtn:                        "成功",
	ParamError:                   "参数错误",
	CustomerError:                "一般错误",
	AuthNeedLoginError:           "该功能需要登录!",
	AuthOutDateError:             "登录态已过期,请重新登录!",
	UserWriteOffError:            "用户不存在",
	UserBalanceLowError:          "您的余额不足",
	UserExistError:               "用户已存在",
	UserAkExistError:             "用户AccessKey已存在",
	AuthTokenError:               "没有授权中心的权限",
	AuthGenerateAccessTokenError: "生成AccessToken失败",
}

func init() {
	rtnerr.NewErr(rtnCodeMap)
}
