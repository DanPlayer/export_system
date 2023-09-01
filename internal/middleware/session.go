package middleware

import (
	"export_system/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	TokenName         = "token"
	CrontabSecretName = "secret"
	CrontabSecret     = "YourShines(23df20037sd)Y"
)

func CrossDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "1728000")
		c.Header("Access-Control-Allow-Headers", "Accept,Origin,X-Requested-With,Content-Type,token,sign,app_id,timestamp")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Origin", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func CrossDomainForDebug() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "1728000")
		c.Header("Access-Control-Allow-Headers", "Accept,Origin,X-Requested-With,Content-Type,token,sign,app_id,timestamp")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Origin", "*")
		if c.Request.Method == http.MethodGet {
			log.Println("get ", c.Request.URL.RequestURI())
			c.Next()
		} else if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			raw, err := utils.CopyBody(c.Request)
			if err == nil {
				log.Println("----------")
				log.Println(string(raw))
				log.Println("----------")
			}
			c.Next()
		} else if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			log.Println(c.Request.Method, c.Request.RequestURI)
			c.Next()
		}
	}
}

// GetToken 获取token
func GetToken(c *gin.Context) string {
	if len(c.Query(TokenName)) > 0 {
		return c.Query(TokenName)
	} else if len(c.PostForm(TokenName)) > 0 {
		return c.PostForm(TokenName)
	} else if len(c.GetHeader(TokenName)) > 0 {
		return c.GetHeader(TokenName)
	} else {
		if t, err := c.Cookie(TokenName); err == nil && len(t) > 0 {
			return t
		}
	}
	return ""
}

// Auth 权限检查
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetToken(c)
		if len(token) > 0 {
			info, err := ParseToken(token)
			if err != nil || info.UserID == "" {
				if err != nil {
				}
				utils.OutAuthOutdatedError(c)
				c.Abort()
				return
			}
			c.Set("userID", info.UserID)
			c.Set("timUserID", info.TIMUserID)
			c.Set("name", info.NickName)
			c.Set("avatar", info.Avatar)
			c.Set("token", token)
			c.Set("phone", info.Phone)
			c.Next()
			return
		}
		utils.OutAuthNeedError(c)
		c.Abort()
		return
	}
}

// CrontabAuth 脚本权限检查
func CrontabAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query(CrontabSecretName) == CrontabSecret {
			c.Next()
		} else {
			utils.OutAuthNeedError(c)
			c.Abort()
		}
	}
}

// DocAuth 文档权限检查
func DocAuth(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// GetLoginTimUserID 获取登录的腾讯IM用户ID
func GetLoginTimUserID(c *gin.Context) string {
	userid, _ := c.Get("timUserID")
	return userid.(string)
}

// GetLoginUserID 获取登录的用户ID
func GetLoginUserID(c *gin.Context) string {
	userid, _ := c.Get("userID")
	return userid.(string)
}

// GetLoginUserNikeName 获取登录用户的昵称
func GetLoginUserNikeName(c *gin.Context) string {
	nickName, _ := c.Get("name")
	return nickName.(string)
}

// GetLoginUserAvatar 获取登录用户的头像
func GetLoginUserAvatar(c *gin.Context) string {
	avatar, _ := c.Get("avatar")
	return avatar.(string)
}

// GetLoginUserPhone 获取登录用户的手机号
func GetLoginUserPhone(c *gin.Context) string {
	phone, _ := c.Get("phone")
	return phone.(string)
}

// GetLoginStoreID 获取登录商户的ID
func GetLoginStoreID(c *gin.Context) string {
	storeId, _ := c.Get("storeID")
	return storeId.(string)
}

// GetLoginStoreToken 获取登录商户的Token
func GetLoginStoreToken(c *gin.Context) string {
	token, _ := c.Get("storeToken")
	return token.(string)
}

// GetLoginAdminID 获取登录管理员的ID
func GetLoginAdminID(c *gin.Context) string {
	adminId, _ := c.Get("adminID")
	return adminId.(string)
}

// GetLoginUserToken 获取登陆用户的token
func GetLoginUserToken(c *gin.Context) string {
	token, _ := c.Get("token")
	return token.(string)
}
