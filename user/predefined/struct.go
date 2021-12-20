package predefined

import (
	"github.com/golang-jwt/jwt"
	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/user/database/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTTokenClaimsUserInfo struct {
	UserID      primitive.ObjectID `json:"user_id"`
	ExtensionID int                `json:"extension_id"`
	Name        string             `json:"name"`
	Guest       bool               `json:"guest"`
	Level       int                `json:"level"`
	Experience  int                `json:"experience"`
	Roles       []model.UserRole   `json:"roles"`
	Tags        []model.UserTag    `json:"tags"`
	IP          string             `json:"ip"`
	UserAgent   string             `json:"user_agent"`
	Location    geometry.Point     `json:"location"`
}

type JWTTokenClaims struct {
	*jwt.StandardClaims
	*JWTTokenClaimsUserInfo        //数据变化，需刷新token获取最新token
	TokenType               string `json:"token_type"`
}
