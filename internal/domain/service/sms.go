package service

import (
	"export_system/internal/rdb"
	"export_system/internal/sms"
	"export_system/utils"
	"fmt"
)

const (
	AliSmsSign = ""

	AliSmsCode = ""
)

// SendLoginSms 发送用户登录验证短信
func SendLoginSms(phone string) error {
	// 生成登录验证码，存储至redis
	code := utils.GenValidateCode(6)
	smsForLoginRdb := rdb.SmsForLogin{Phone: phone, Code: code}
	err := smsForLoginRdb.Set()
	if err != nil {
		return err
	}
	return sms.SendSmsCustomer(phone, AliSmsSign, AliSmsCode, fmt.Sprintf(`{"code":"%s"}`, code))
}
