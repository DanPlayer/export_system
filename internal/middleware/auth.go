package middleware

import (
	"encoding/json"
	"errors"
	"export_system/internal/rdb"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// ParseToken 验证token密钥
func ParseToken(token string) (info rdb.Token, err error) {
	tokenRdb := rdb.Token{Token: token}
	get, err := tokenRdb.Get()
	if err != nil {
		return
	}
	if get == "" {
		err = errors.New("token验证失败")
		return
	}
	err = json.Unmarshal([]byte(get), &info)
	if err != nil {
		return
	}

	return
}

// MakeToken 生成token
func MakeToken(phone string, expired time.Time) (tokenStr string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone":   phone,
		"expired": expired.Unix(),
		"issuer":  "yourshines",
	})
	tokenStr, err = token.SignedString([]byte("yourshines%2#@sa"))
	return
}
