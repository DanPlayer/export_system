package model

import (
	"export_system/internal/db"
	"export_system/internal/domain/pojo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User 用户
type User struct {
	gorm.Model
	Uid          string `gorm:"type:varchar(80);unique;not null;comment:'用户唯一标识'"`             //用户唯一标识
	Phone        string `gorm:"type:varchar(255);unique;comment:'手机号'"`                        // 手机号
	Password     string `gorm:"type:varchar(255);comment:'密码'"`                                // 密码
	NickName     string `gorm:"type:varchar(255);comment:'昵称'"`                                // 昵称
	Avatar       string `gorm:"type:text;comment:'头像URL'"`                                     // 头像
	Forbidden    bool   `gorm:"type:tinyint(1);default:0;comment:'是否禁用 0-正常， 1-被禁用'"`          // 是否禁用 0-正常， 1-被禁用
	AccessKey    string `gorm:"type:varchar(255);index:idx_access_key_secret;comment:'密令Key'"` // 密令Key
	AccessSecret string `gorm:"type:text;index:idx_access_key_secret;comment:'密令'"`            // 密令
}

func (m *User) TableName() string {
	return "user"
}

func (m *User) InfoByIncludeDeleted(phone string) (info User, err error) {
	err = db.MasterClient.Model(&m).Unscoped().Clauses(clause.Locking{Strength: "UPDATE"}).Where("phone = ?", phone).First(&info).Error
	return
}

func (m *User) InfoByKeyIncludeDeleted(key string) (info User, err error) {
	err = db.MasterClient.Model(&m).Unscoped().Clauses(clause.Locking{Strength: "UPDATE"}).Where("access_key = ?", key).First(&info).Error
	return
}

func (m *User) InfoByAccess(key, secret string) (info User, err error) {
	err = db.MasterClient.Model(&m).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("access_key = ?", key).
		Where("access_secret = ?", secret).
		First(&info).Error
	return
}

func (m *User) InfoBy(phone string) (info User, err error) {
	err = db.MasterClient.Model(&m).Clauses(clause.Locking{Strength: "UPDATE"}).Where("phone = ?", phone).First(&info).Error
	return
}

func (m *User) InfoByUid(uid string) (info User, err error) {
	err = db.MasterClient.Model(&m).Clauses(clause.Locking{Strength: "UPDATE"}).Where("uid = ?", uid).First(&info).Error
	return
}

func (m *User) InfoById(id string) (info User, err error) {
	err = db.MasterClient.Model(&m).Where("uid = ? OR tim_user_id = ?", id, id).First(&info).Error
	return
}

func (m *User) TxCreate(tx *gorm.DB) error {
	return tx.Model(&m).Create(&m).Error
}

func (m *User) TxUpdateBasicInfoByTimUserId(tx *gorm.DB, timUserId string) error {
	upMap := m.setMap()
	return tx.Model(&m).Where("tim_user_id = ?", timUserId).Updates(upMap).Error
}

func (m *User) setMap() map[string]interface{} {
	upMap := map[string]interface{}{}
	if m.NickName != "" {
		upMap["nick_name"] = m.NickName
	}
	if m.Avatar != "" {
		upMap["avatar"] = m.Avatar
	}
	return upMap
}

// InfoListBy  批量查询用户信息
func (m *User) InfoListBy(ids []string) (list []pojo.UserInfo) {
	if len(ids) == 0 {
		return make([]pojo.UserInfo, 0)
	}
	db.MasterClient.Model(&m).Where("uid IN ?", ids).Find(&list)
	return
}

func (m *User) MapUserInfoBy(ids []string) (mapInfo map[string]User, err error) {
	mapInfo = make(map[string]User, 0)
	if len(ids) == 0 {
		return
	}
	var list []User
	err = db.MasterClient.Model(&m).Where("uid IN ?", ids).Find(&list).Error
	if err != nil {
		return
	}
	for i := range list {
		mapInfo[list[i].Uid] = list[i]
	}
	return
}

func (m *User) SearchUserList(name, phone, channelCode string, page, size int) (list []User, count int64, err error) {
	table := db.MasterClient.Model(&m)
	if name != "" {
		table.Where("nick_name LIKE ?", name+"%")
	}
	if phone != "" {
		table.Where("phone = ?", phone)
	}
	if channelCode != "" {
		table.Where("register_from_channel = ?", channelCode)
	}
	table.Order("created_at DESC")
	table.Count(&count)
	err = table.Limit(size).Offset((page - 1) * size).Find(&list).Error
	return
}

// UpdateForbidden 更新用户的禁用状态
func (m *User) UpdateForbidden(userID string, isForbidden bool) error {
	return db.MasterClient.Model(&m).Where("uid = ?", userID).UpdateColumn("forbidden", isForbidden).Error
}

// TxWriteOff 注销用户
func (m *User) TxWriteOff(tx *gorm.DB, userID string) error {
	return tx.Where("uid = ?", userID).Delete(&m).Error
}
