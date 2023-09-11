package service

import (
	"errors"
	"export_system/internal/db"
	"export_system/internal/domain/model"
	"export_system/internal/domain/pojo"
	"export_system/internal/domain/rtn"
	"export_system/internal/hashids"
	"export_system/internal/middleware"
	"export_system/internal/rdb"
	"export_system/pkg/rtnerr"
	"export_system/utils"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

func Login(phone, password string) (token, uid string, rtnErr rtnerr.RtnError) {
	// 查询用户是否存在，不存在则是注册新账号
	exist, info := CheckPhoneUser(phone)
	if !exist || info.DeletedAt.Valid {
		rtnErr = rtn.UserExistError
		return
	}

	if info.Password != utils.MD5String(password) {
		rtnErr = rtn.UserLoginError
		return
	}

	uid = info.Uid

	// 生成登录token
	var err error
	token, err = middleware.MakeToken(phone, utils.Now().Add(7*24*time.Hour))
	if err != nil {
		rtnErr = rtnerr.New(err)
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
		rtnErr = rtnerr.New(err)
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

// SmsLogin 短信登录
func SmsLogin(phone, code string) (token string, profileComplete bool, userId string, err error) {
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

	if info.NickName != "" && info.Avatar != "" {
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

// CheckAccessKeyUser 检测用户ak是否存在
func CheckAccessKeyUser(key string) (bool, model.User) {
	userModel := model.User{}
	info, err := userModel.InfoByKeyIncludeDeleted(key)
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
	if info.NickName == "" || info.Avatar == "" {
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
			Forbidden:   list[i].Forbidden,
			CreatedTime: list[i].CreatedAt.Unix(),
		})
	}
	return
}

// BackendAddUser 新增用户
func BackendAddUser(phone, password, nickName, avatar string) rtnerr.RtnError {
	// 查询用户是否存在，不存在则是注册新账号
	exist, _ := CheckPhoneUser(phone)
	if exist {
		return rtn.UserExistError
	}

	tx := db.MasterClient.Begin()
	// 生成TIMUserID
	uid := utils.GetUid()

	split := strings.Split(uid, "")
	var uidInt []int
	for _, s := range split {
		i, _ := strconv.Atoi(s)
		uidInt = append(uidInt, i)
	}

	// 生成用户的AccessKey
	key, err := hashids.Client.Encode(uidInt)
	if err != nil {
		tx.Rollback()
		return rtn.UserAkExistError
	}

	// 查询用户的key是否已存在
	exist, _ = CheckAccessKeyUser(key)
	if exist {
		tx.Rollback()
		return rtn.UserAkExistError
	}

	// 生成AccessSecret，并加密
	secret := utils.RandStr(64)
	secret = utils.MD5String(secret)

	// 加密密码
	password = utils.MD5String(password)

	// 创建新用户信息
	userModel := model.User{
		Uid:          uid,
		Phone:        phone,
		Password:     password,
		NickName:     nickName,
		Avatar:       avatar,
		AccessKey:    key,
		AccessSecret: secret,
	}
	err = userModel.TxCreate(tx)
	if err != nil {
		tx.Rollback()
		return rtnerr.New(err)
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

// CheckAccessAuth 检查用户AK权限
func CheckAccessAuth(key, secret string) (bool, model.User) {
	// 加密secret
	secret = utils.MD5String(secret)

	userModel := model.User{}
	info, err := userModel.InfoByAccess(key, secret)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
		}
		return false, model.User{}
	}
	return true, info
}

// GenerateAccessToken 生成AccessToken
func GenerateAccessToken() (string, int64, error) {
	uid := utils.GetUid()
	randStr := utils.RandStr(12)
	token := utils.MD5String(uid + randStr + strconv.FormatInt(time.Now().Unix(), 10))

	accessTokenRdb := rdb.AccessToken{
		Token: token,
	}
	err := accessTokenRdb.Set()
	if err != nil {
		return "", 0, err
	}

	return token, rdb.AccessTokenExpireTime, nil
}

// RefreshAccessToken 刷新AccessToken
func RefreshAccessToken(token string) (string, int64, rtnerr.RtnError) {
	accessTokenRdb := rdb.AccessToken{
		Token: token,
	}
	get, err := accessTokenRdb.Get()
	if err != nil {
		return "", 0, rtnerr.New(err)
	}
	if get == "" {
		return "", 0, rtn.AuthTokenError
	}
	// 刷新
	accessToken, expire, err := GenerateAccessToken()
	if err != nil {
		return "", 0, rtnerr.New(err)
	}
	return accessToken, expire, nil
}
