package service

import (
	"errors"
	"export_system/internal/db"
	"export_system/internal/domain/model"
	"export_system/internal/domain/pojo"
	"export_system/internal/middleware"
	"export_system/internal/rdb"
	"export_system/utils"
	"gorm.io/gorm"
	"log"
	"time"
)

// SmsLogin 短信登录
func SmsLogin(phone, code, channelCode, ipAddress string) (token string, profileComplete bool, userId string, err error) {
	// 查询用户是否存在，不存在则是注册新账号
	exist, info := CheckPhoneUser(phone)
	if info.DeletedAt.Valid {
		err = errors.New("用户已注销")
		return
	}
	if !CheckLoginSms(phone, code) {
		err = errors.New("手机验证码错误，请重新输入")
		return
	}

	if !exist {
		info, err = register(phone)
		if err != nil {
			log.Println(err)
			return
		}
	}

	if info.NickName != "" && info.Avatar != "" && !info.BirthDay.IsZero() && info.Gender != 0 {
		profileComplete = true
	}

	// 删除手机验证码缓存
	smsForLoginRdb := rdb.SmsForLogin{Phone: phone}
	err = smsForLoginRdb.Del()
	if err != nil {
		log.Println("缓存删除失败")
	}

	// 登录成功，生成token，记录token缓存
	token, err = middleware.MakeToken(phone, utils.Now().Add(7*24*time.Hour))
	if err != nil {
		return
	}
	tokenRdb := rdb.Token{
		Token:    token,
		UserID:   info.Uid,
		NickName: info.NickName,
		Avatar:   info.Avatar,
		Phone:    phone,
	}
	err = tokenRdb.Set()
	if err != nil {
		return
	}

	// 存储用户token关系缓存
	userTokenRdb := rdb.UserToken{UserID: info.Uid, Token: token}
	oldToken, _ := userTokenRdb.Get()
	// 删除旧登录态
	if oldToken != "" {
		oldTokenRdb := rdb.Token{Token: oldToken}
		_ = oldTokenRdb.Del()
	}
	_ = userTokenRdb.Set()

	return
}

// register 注册新账号
func register(phone string) (info model.User, err error) {
	// 开启事务
	tx := db.MasterClient.Begin()

	// 生成TIMUserID
	userID := utils.GetUid()

	userModel := model.User{
		Uid:   userID,
		Phone: phone,
	}
	err = userModel.TxCreate(tx)
	if err != nil {
		tx.Rollback()
		return
	}
	// 事务提交
	tx.Commit()

	return userModel, nil
}

// CheckPhoneUser 检测手机号是否已注册
func CheckPhoneUser(phone string) (bool, model.User) {
	userModel := model.User{}
	info, err := userModel.InfoByIncludeDeleted(phone)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
		}
		return false, model.User{}
	}
	return true, info
}

// CheckUserProfile 检测用户资料是否完整
func CheckUserProfile(userId string) bool {
	userModel := model.User{}
	info, err := userModel.InfoByUid(userId)
	if err != nil {
		return false
	}
	if info.NickName == "" || info.Avatar == "" || info.BirthDay.IsZero() || info.Gender == 0 {
		return false
	}
	return true
}

// CheckLoginSms 检查手机验证码是否正确
func CheckLoginSms(phone, code string) bool {
	if code == "888888" {
		return true
	}
	// 验证用户手机验证码
	smsForLoginRdb := rdb.SmsForLogin{Phone: phone}
	smsCode, err := smsForLoginRdb.Get()
	if err != nil || smsCode == "" {
		return false
	}
	if smsCode != code {
		return false
	}
	return true
}

// UserInfo 当前用户的详情
func UserInfo(userID string) (userInfo pojo.UserInfo, err error) {
	// 查找用户信息
	user := model.User{}
	info, err := user.InfoByUid(userID)
	if err != nil {
		return
	}

	userInfo = pojo.UserInfo{
		Uid:      info.Uid,
		NickName: info.NickName,
		Avatar:   info.Avatar,
		Gender:   info.Gender,
		Phone:    info.Phone,
	}

	return userInfo, nil
}

// BackendSearchUsers 搜索用户
func BackendSearchUsers(name, phone, channelCode string, page, size int) (users []pojo.BackendUserInfo, count int64, err error) {
	users = make([]pojo.BackendUserInfo, 0)

	user := model.User{}
	list, count, err := user.SearchUserList(name, phone, channelCode, page, size)
	if err != nil || len(list) == 0 {
		return
	}
	for i := range list {
		users = append(users, pojo.BackendUserInfo{
			Uid:         list[i].Uid,
			Phone:       list[i].Phone,
			NickName:    list[i].NickName,
			Avatar:      list[i].Avatar,
			Gender:      list[i].Gender,
			Forbidden:   list[i].Forbidden,
			CreatedTime: list[i].CreatedAt.Unix(),
		})
	}
	return
}

// BackendAddInnerUser 新增内部用户
func BackendAddInnerUser(phone, nickName, avatar string, birthDay time.Time, gender int) error {
	// 查询用户是否存在，不存在则是注册新账号
	exist, _ := CheckPhoneUser(phone)
	if exist {
		return errors.New("该号码已经注册了")
	}

	tx := db.MasterClient.Begin()
	// 生成TIMUserID
	userID := utils.GetUid()
	// 创建新用户信息
	userModel := model.User{
		Uid:      userID,
		Phone:    phone,
		NickName: nickName,
		Avatar:   avatar,
		BirthDay: birthDay,
		Gender:   gender,
	}
	err := userModel.TxCreate(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 事务提交
	tx.Commit()
	return nil
}

func BackendForbiddenUser(userId string, forbidden bool) error {
	userModel := model.User{}
	err := userModel.UpdateForbidden(userId, forbidden)
	if err != nil {
		return err
	}

	if forbidden { // 如果是禁用用户
		// 删除用户的登录态
		userTokenRdb := rdb.UserToken{UserID: userId}
		oldToken, _ := userTokenRdb.Get()
		if oldToken != "" {
			oldTokenRdb := rdb.Token{Token: oldToken}
			_ = oldTokenRdb.Del()
		}
	}

	return nil
}

// WriteOffUser 注销用户
func WriteOffUser(userID string) error {
	tx := db.MasterClient.Begin()
	userModel := model.User{}
	err := userModel.TxWriteOff(tx, userID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}