package predefined

import (
	"github.com/golang-jwt/jwt"
)

type JWTTokenClaimsUserInfo struct {
	UserID      string   `json:"user_id"`
	ExtensionID int      `json:"extension_id"`
	Name        string   `json:"name"`
	Guest       bool     `json:"guest"`
	Level       int      `json:"level"`
	Role        []string `json:"role"`
}

type JWTTokenClaims struct {
	*jwt.StandardClaims
	*JWTTokenClaimsUserInfo        //数据变化，需刷新token获取最新token
	TokenType               string `json:"token_type"`
}
