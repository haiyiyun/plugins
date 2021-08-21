package base

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/haiyiyun/plugins/urbac/database/model"
	"github.com/haiyiyun/plugins/urbac/database/model/token"
	"github.com/haiyiyun/plugins/urbac/database/model/user"
	"github.com/haiyiyun/plugins/urbac/predefined"

	"github.com/golang-jwt/jwt"
	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) Login(username, password, ip, userAgent string) (m help.M, err error) {
	passwordMd5 := help.Strings(password).Md5()

	userModel := user.NewModel(self.M)
	filter := bson.D{
		{"name", username},
		{"password", passwordMd5},
		{"enable", true},
		{"delete", false},
	}

	u := model.User{}
	ctx := context.TODO()
	sr := userModel.FindOne(ctx, filter)
	if err = sr.Decode(&u); err == nil {
		tokenModel := token.NewModel(self.M)
		cnt, _ := tokenModel.CountDocuments(ctx, tokenModel.FilterByUserID(u.ID))
		if cnt == 0 || (self.AllowMultiLogin && (self.AllowMultiLoginNum == 0 || cnt < self.AllowMultiLoginNum)) {
			ExpiredTime := time.Now().Add(self.Config.TokenExpireDuration.Duration)
			jwtID := primitive.NewObjectID()
			claims := &predefined.JWTTokenClaims{
				StandardClaims: &jwt.StandardClaims{
					Id:        jwtID.Hex(),
					Audience:  u.ID.Hex(),
					Issuer:    u.ID.Hex(),
					Subject:   u.Name,
					ExpiresAt: ExpiredTime.Unix(),
				},
				TokenType: predefined.TokenTypeSelf,
			}

			jwtToken := jwt.NewWithClaims(predefined.JWTSigningMethod, claims)
			if tokenString, jwtErr := jwtToken.SignedString([]byte(u.Password)); jwtErr == nil {
				if _, err = tokenModel.Create(ctx, model.Token{
					ID:        jwtID,
					UserID:    u.ID,
					UserName:  u.Name,
					TokenType: predefined.TokenTypeSelf,
					Token:     tokenString,
					SignInfo: model.TokenSignInfo{
						IP:        ip,
						UserAgent: userAgent,
					},
					ExpiredTime: ExpiredTime,
				}); err == nil {
					m = map[string]interface{}{
						"userId": u.ID.Hex(),
						"token":  tokenString,
					}
				}
			} else {
				err = jwtErr
			}
		} else {
			err = errors.New(predefined.StatusCodeLoginLimitText)
		}
	}

	return
}

func (self *Service) Logout(r *http.Request) {
	if claims, _ := self.GetValidClaims(r); claims != nil {
		tokenString, _ := self.BearerAuth(r)
		if tokenID, err := primitive.ObjectIDFromHex(claims.Id); err == nil {
			tokenModel := token.NewModel(self.M)
			tokenModel.DeleteOne(context.TODO(), tokenModel.FilterByID(tokenID))
			cacheClaimsKey := "claims.valid." + tokenString
			self.Cache.Delete(cacheClaimsKey)
			cacheApplicationsInfoKey := "applications.info." + claims.Issuer
			self.Cache.Delete(cacheApplicationsInfoKey)
		}
	}
}

func (self *Service) GetClaims(r *http.Request) (claims *predefined.JWTTokenClaims) {
	if tokenString, found := self.BearerAuth(r); found {
		if token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &predefined.JWTTokenClaims{}); err == nil {
			claims = token.Claims.(*predefined.JWTTokenClaims)
		}
	}

	return
}

func (self *Service) GetValidClaims(r *http.Request) (claims *predefined.JWTTokenClaims, u model.User) {
	if tokenString, found := self.BearerAuth(r); found {
		cacheKey := "claims.valid." + tokenString
		if claimsAndUserI, found := self.Cache.Get(cacheKey); found {
			claimsAndUser := claimsAndUserI.(help.M)
			claims = claimsAndUser["claims"].(*predefined.JWTTokenClaims)
			u = claimsAndUser["user"].(model.User)
		} else {
			uTmp := model.User{}
			token, _ := jwt.ParseWithClaims(tokenString, &predefined.JWTTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
				claimsTmp := t.Claims.(*predefined.JWTTokenClaims)

				var key []byte
				var jwtID primitive.ObjectID
				var userID primitive.ObjectID
				jwtIDHex := claimsTmp.Id
				userIDHex := claimsTmp.Issuer

				jwtID, err := primitive.ObjectIDFromHex(jwtIDHex)
				if err == nil {
					tokenModel := token.NewModel(self.M)
					var cnt int64
					if cnt, err = tokenModel.CountDocumentsByIDAndToken(jwtID, tokenString); err == nil && cnt > 0 {
						if userID, err = primitive.ObjectIDFromHex(userIDHex); err == nil {
							tokenType := ""
							if self.Config.AllowMultiLogin {
								tokenType = predefined.TokenTypeSelf
							}

							cnt, err = tokenModel.CountDocumentsByUserIDAndType(userID, tokenType)
							if err == nil && cnt > 0 {
								if cnt == 1 || (self.Config.AllowMultiLogin && (self.Config.AllowMultiLoginNum == 0 || cnt <= self.Config.AllowMultiLoginNum)) {
									u, err := self.getUser(userID)
									if err == nil {
										uTmp = u
										key = []byte(u.Password)
									}
								}
							}
						}
					}
				}

				return key, err
			})

			if token.Valid {
				claims = token.Claims.(*predefined.JWTTokenClaims)
				u = uTmp
				self.Cache.Set(cacheKey, help.M{
					"claims": claims,
					"user":   u,
				}, time.Until(time.Unix(claims.ExpiresAt, 0)))
			}
		}
	}

	return
}

func (self *Service) BearerAuth(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}

	return token, token != ""
}
