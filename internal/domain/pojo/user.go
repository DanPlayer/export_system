package pojo

// UserInfo 用户基本信息，基本展示用
type UserInfo struct {
	Uid      string `json:"userID"`   // 用户ID
	NickName string `json:"nickName"` // 用户昵称
	Avatar   string `json:"avatar"`   // 用户头像
	Phone    string `json:"phone"`    // 手机号码
}

type BackendUserInfo struct {
	Uid         string `json:"userID"`      // 用户ID
	Phone       string `json:"phone"`       // 手机号
	NickName    string `json:"nickName"`    // 用户昵称
	Avatar      string `json:"avatar"`      // 用户头像
	Forbidden   bool   `json:"forbidden"`   // 是否被禁用
	CreatedTime int64  `json:"createdTime"` // 注册时间
}

type BackendUserVerify struct {
	UserID       string `json:"userId"`       // 用户ID
	RealName     string `json:"realName"`     // 真实姓名
	CardType     int    `json:"cardType"`     // 证件类型
	IdCardNumber string `json:"idCardNumber"` // 身份证号码
	IdCardFace   string `json:"idCardFace"`   // 身份证人脸面
	IdCardBack   string `json:"idCardBack"`   // 身份证国徽面
	Status       int    `json:"status"`       // 状态：1 - 待审核 2 - 审核通过 3 - 审核不通过
	CheckAdminID string `json:"checkAdminID"` // 审核管理员ID
}
