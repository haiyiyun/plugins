package predefined

import (
	"github.com/golang-jwt/jwt"
)

type JWTTokenClaims struct {
	*jwt.StandardClaims
	TokenType string `json:"token_type"`
}
