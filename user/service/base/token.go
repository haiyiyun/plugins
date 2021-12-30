package base

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/haiyiyun/plugins/user/database/model"
	"github.com/haiyiyun/plugins/user/database/model/token"
	"github.com/haiyiyun/plugins/user/database/model/user"
	"github.com/haiyiyun/plugins/user/predefined"

	"github.com/golang-jwt/jwt"
	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Service) Login(ctx context.Context, username, password, ip, userAgent string, coordinates geometry.PointCoordinates, expiredTime time.Time) (m help.M, err error) {
	userModel := user.NewModel(self.M)

	var u model.User
	if u, err = userModel.CheckNameAndPassword(username, password); err == nil {
		m, err = self.CreateToken(ctx, u, ip, userAgent, coordinates, expiredTime)
	}

	return
}

func (self *Service) LoginByUserID(ctx context.Context, userID primitive.ObjectID, ip, userAgent string, coordinates geometry.PointCoordinates, expiredTime time.Time) (m help.M, err error) {
	userModel := user.NewModel(self.M)

	var u model.User
	if u, err = userModel.GetUserByID(userID); err == nil {
		m, err = self.CreateToken(ctx, u, ip, userAgent, coordinates, expiredTime)
	}

	return
}

func (self *Service) Logout(r *http.Request) {
	if claims, _ := self.GetValidClaims(r); claims != nil {
		tokenString, _ := self.BearerAuth(r)
		if tokenID, err := primitive.ObjectIDFromHex(claims.Id); err == nil {
			tokenModel := token.NewModel(self.M)
			tokenModel.DeleteOne(r.Context(), tokenModel.FilterByID(tokenID))
			cacheClaimsKey := "claims.valid." + tokenString
			self.Cache.Delete(cacheClaimsKey)
		}
	}
}

func (self *Service) CreateToken(ctx context.Context, u model.User, ip, userAgent string, coordinates geometry.PointCoordinates, expiredTime time.Time) (m help.M, err error) {
	tokenModel := token.NewModel(self.M)

	if !help.NewSlice(self.Config.OnlySingleLoginUserIDUnlimited).CheckItem(u.ID.Hex()) {
		if self.Config.OnlySingleLogin {
			if _, err = tokenModel.DeleteMany(ctx, tokenModel.FilterByUserID(u.ID)); err != nil {
				return
			}
		}

		if len(self.Config.OnlySingleLoginUserID) > 0 {
			if help.NewSlice(self.Config.OnlySingleLoginUserID).CheckItem(u.ID.Hex()) {
				if _, err = tokenModel.DeleteMany(ctx, tokenModel.FilterByUserID(u.ID)); err != nil {
					return
				}
			}
		}
	}

	cnt, _ := tokenModel.CountDocuments(ctx, tokenModel.FilterByUserID(u.ID))

	if cnt == 0 ||
		(self.AllowMultiLogin &&
			((len(self.AllowMultiLoginUserIDUnlimited) > 0 && help.NewSlice(self.AllowMultiLoginUserIDUnlimited).CheckItem(u.ID.Hex())) ||
				(self.AllowMultiLoginNum == 0 || cnt < self.AllowMultiLoginNum))) {

		if expiredTime.IsZero() {
			expiredTime = time.Now().Add(self.Config.TokenExpireDuration.Duration)
			if dur, found := self.Config.SpecifyUserIDTokenExpireDuration[u.ID.Hex()]; found {
				expiredTime = time.Now().Add(dur.Duration)
			}
		}

		jwtID := primitive.NewObjectID()

		//过滤未开始或者已经到期的
		roles := []model.UserRole{}
		if len(u.Roles) > 0 {
			for _, role := range u.Roles {
				if role.EndTime.IsZero() {
					roles = append(roles, role)
				} else if !role.StartTime.IsZero() && !role.EndTime.IsZero() {
					if role.EndTime.After(role.StartTime) {
						roles = append(roles, role)
					}
				}
			}
		}

		//过滤未开始或者已经到期的
		tags := []model.UserTag{}
		if len(u.Roles) > 0 {
			for _, tag := range u.Tags {
				if tag.EndTime.IsZero() {
					tags = append(tags, tag)
				} else if !tag.StartTime.IsZero() && !tag.EndTime.IsZero() {
					if tag.EndTime.After(tag.StartTime) {
						tags = append(tags, tag)
					}
				}
			}
		}

		claims := &predefined.JWTTokenClaims{
			StandardClaims: &jwt.StandardClaims{
				Id:        jwtID.Hex(),
				Audience:  u.ID.Hex(),
				Issuer:    u.ID.Hex(),
				Subject:   u.Name,
				ExpiresAt: expiredTime.Unix(),
			},
			TokenType: predefined.TokenTypeSelf,
			JWTTokenClaimsUserInfo: &predefined.JWTTokenClaimsUserInfo{
				UserID:      u.ID,
				ExtensionID: u.ExtensionID,
				Name:        u.Name,
				Guest:       u.Guest,
				Level:       u.Level,
				Experience:  u.Experience,
				Roles:       roles,
				Tags:        tags,
				IP:          ip,
				UserAgent:   userAgent,
				Location:    geometry.NewPoint(coordinates),
			},
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
					Location:  geometry.NewPoint(coordinates),
				},
				ExpiredTime: expiredTime,
			}); err == nil {
				m = map[string]interface{}{
					"user_id": u.ID.Hex(),
					"token":   tokenString,
				}
			}
		} else {
			err = jwtErr
		}
	} else {
		err = errors.New(predefined.StatusCodeLoginLimitText)
	}

	return
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
			if token, err := jwt.ParseWithClaims(tokenString, &predefined.JWTTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
				var key []byte
				var err error
				if t.Claims != nil {
					if claimsTmp, ok := t.Claims.(*predefined.JWTTokenClaims); ok {
						var jwtID primitive.ObjectID
						var userID primitive.ObjectID
						jwtIDHex := claimsTmp.Id
						userIDHex := claimsTmp.Issuer

						jwtID, err = primitive.ObjectIDFromHex(jwtIDHex)
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
										if cnt == 1 ||
											(self.Config.AllowMultiLogin &&
												((len(self.AllowMultiLoginUserIDUnlimited) > 0 && help.NewSlice(self.AllowMultiLoginUserIDUnlimited).CheckItem(userID.Hex())) ||
													(self.Config.AllowMultiLoginNum == 0 || cnt <= self.Config.AllowMultiLoginNum))) {
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
					}
				}

				return key, err
			}); err == nil {
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

	if token == "" && self.Config.TokenByUrlQuery {
		queryName := "token"
		if self.TokenByUrlQueryName != "" {
			queryName = self.TokenByUrlQueryName
		}

		token = r.URL.Query().Get(queryName)
	}

	return token, token != ""
}

func (self *Service) GetTokensByUsernameAndPassword(ctx context.Context, username, password string) (ts []help.M, err error) {
	userModel := user.NewModel(self.M)

	var u model.User
	if u, err = userModel.CheckNameAndPassword(username, password); err == nil {
		tokenModel := token.NewModel(self.M)
		var cur *mongo.Cursor
		if cur, err = tokenModel.Find(ctx, tokenModel.FilterByUserID(u.ID), options.Find().SetProjection(bson.D{
			{"token", 0},
		})); err == nil {
			err = cur.All(ctx, &ts)
		}
	}

	return
}

func (self *Service) DeleteTokenByUsernameAndPassword(ctx context.Context, tokenID primitive.ObjectID, username, password string) error {
	userModel := user.NewModel(self.M)

	u, err := userModel.CheckNameAndPassword(username, password)
	if err == nil {
		tokenModel := token.NewModel(self.M)

		filter := tokenModel.FilterByUserID(u.ID)
		filter = append(filter, tokenModel.FilterByID(tokenID)...)

		var dr *mongo.DeleteResult
		if dr, err = tokenModel.DeleteOne(ctx, filter); err == nil {
			if dr.DeletedCount == 0 {
				err = mongo.ErrNoDocuments
			}
		}
	}

	return err
}

func (self *Service) GetTokensByToken(r *http.Request) (ts []help.M, err error) {
	if claims, u := self.GetValidClaims(r); claims != nil {
		tokenID, _ := primitive.ObjectIDFromHex(claims.Id)
		tokenModel := token.NewModel(self.M)

		filter := tokenModel.FilterByUserID(u.ID)
		filter = append(filter, bson.D{
			{"$neq", bson.D{
				{"_id", tokenID},
			}},
		}...)

		var cur *mongo.Cursor
		if cur, err = tokenModel.Find(r.Context(), filter, options.Find().SetProjection(bson.D{
			{"token", 0},
		})); err == nil {
			err = cur.All(r.Context(), &ts)
		}
	} else {
		err = errors.New("Invalid Claims")
	}

	return
}

func (self *Service) DeleteTokenByToken(tokenID primitive.ObjectID, r *http.Request) (err error) {
	if claims, u := self.GetValidClaims(r); claims != nil {
		tokenID, _ := primitive.ObjectIDFromHex(claims.Id)
		tokenModel := token.NewModel(self.M)

		filter := tokenModel.FilterByUserID(u.ID)
		filter = append(filter, tokenModel.FilterByID(tokenID)...)

		var dr *mongo.DeleteResult
		if dr, err = tokenModel.DeleteOne(r.Context(), filter); err == nil {
			if dr.DeletedCount == 0 {
				err = mongo.ErrNoDocuments
			}
		}
	} else {
		err = errors.New("Invalid Claims")
	}

	return
}
