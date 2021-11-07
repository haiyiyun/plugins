package base

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/haiyiyun/plugins/log/database/model"
	modelLog "github.com/haiyiyun/plugins/log/database/model/log"
	"github.com/haiyiyun/plugins/log/predefined"
	"github.com/haiyiyun/webrouter"

	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/realip"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) LogRequestLogin(r *http.Request) primitive.ObjectID {
	if self.logLoginPath(r.URL.Path) {
		return self.LogRequest(r, predefined.LogTypeLogin, self.Config.DefaultLoginDeleteDuration.Duration)
	}

	return primitive.NilObjectID
}

func (self *Service) LogResponseLogin(rw http.ResponseWriter, r *http.Request) {
	if self.logLoginPath(r.URL.Path) {
		self.LogResponse(rw, r)
	}
}

func (self *Service) LogRequestAuth(r *http.Request) primitive.ObjectID {
	if self.logAuthPath(r.URL.Path) {
		return self.LogRequest(r, predefined.LogTypeAuth, self.Config.DefaultAuthDeleteDuration.Duration)
	}

	return primitive.NilObjectID
}

func (self *Service) LogResponseAuth(rw http.ResponseWriter, r *http.Request) {
	if self.logAuthPath(r.URL.Path) {
		self.LogResponse(rw, r)
	}
}

func (self *Service) LogRequestOperate(r *http.Request) primitive.ObjectID {
	if self.excludeLogOperatePath(r.URL.Path) {
		return self.LogRequest(r, predefined.LogTypeOperate, self.Config.DefaultOperateDeleteDuration.Duration)
	}

	return primitive.NilObjectID
}

func (self *Service) LogResponseOperate(rw http.ResponseWriter, r *http.Request) {
	if self.excludeLogOperatePath(r.URL.Path) {
		self.LogResponse(rw, r)
	}
}

func (self *Service) LogRequest(r *http.Request, typ string, deleteDuration time.Duration) primitive.ObjectID {
	referer := r.Header.Get("Referer")
	ip := self.GetRequestAllIP(r)
	header := r.Header
	payload := self.GetRequestPayload(r)
	userID := primitive.NilObjectID
	userName := ""

	if claims := request.GetClaims(r); claims != nil {
		userIDHex := claims.Issuer
		if uid, err := primitive.ObjectIDFromHex(userIDHex); err == nil {
			userID = uid
			userName = claims.Subject
		}
	}

	logID := primitive.NewObjectID()
	go self.logRequest(logID, userID, userName, typ, r.Method, referer, r.URL.Path, r.URL.Query().Encode(), header, payload, ip, deleteDuration)

	return logID
}

func (self *Service) LogResponse(rw http.ResponseWriter, r *http.Request) {
	if lrw, ok := rw.(*webrouter.ResponseWriter); ok {
		if logIDI, found := lrw.GetData("log_id"); found {
			if logID, ok := logIDI.(primitive.ObjectID); ok {
				header := rw.Header()
				data := lrw.GetResData()
				lrw.SetGetResData(false)
				go self.logResponse(logID, header, string(data))
			}
		}
	}
}

func (self *Service) GetRequestIP(r *http.Request) (ip string) {
	ip = realip.RealIP(r)
	return
}

func (self *Service) GetRequestAllIP(r *http.Request) (ip string) {
	ip = r.RemoteAddr + `;` + strings.Join(r.Header["X-Forwarded-For"], "/") + `;` + strings.Join(r.Header["X-Real-IP"], "/")
	return
}

func (self *Service) GetRequestPayload(r *http.Request) (payload string) {
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		payload = "FILE: ,"
		for _, fhs := range r.MultipartForm.File {
			if len(fhs) > 0 {
				for _, fh := range fhs {
					payload += fh.Filename + "(" + strconv.FormatInt(fh.Size, 10) + "),"
				}
			}
		}
	} else {
		if bs, err := request.GetBody(r); err == nil {
			payload = string(bs)
		}
	}

	return
}

func (self *Service) logLoginPath(reqPath string) bool {
	requestLogLoginPath := self.Config.LogLoginPath
	for _, lPath := range requestLogLoginPath {
		if strings.HasPrefix(reqPath, lPath) {
			return true
		}
	}

	return false
}

func (self *Service) logAuthPath(reqPath string) bool {
	requestLogAuthPath := self.Config.LogAuthPath
	for _, aPath := range requestLogAuthPath {
		if strings.HasPrefix(reqPath, aPath) {
			return true
		}
	}

	return false
}

func (self *Service) excludeLogOperatePath(reqPath string) bool {
	requestLogOperateExcludePath := append(self.Config.LogOperateExcludePath, self.Config.LogLoginPath...)
	requestLogOperateExcludePath = append(requestLogOperateExcludePath, self.Config.LogAuthPath...)
	for _, ePath := range requestLogOperateExcludePath {
		if strings.HasPrefix(reqPath, ePath) {
			return false
		}
	}

	return true
}

func (self *Service) logRequest(logID primitive.ObjectID, userId primitive.ObjectID, username, typ, method, referer, path, query string, header map[string][]string, payload, ip string, deleteDuration time.Duration) {
	deleteTime := time.Now().Add(self.Config.DefaultDeleteDuration.Duration)
	if deleteDuration > 0 {
		deleteTime = time.Now().Add(deleteDuration)
	}

	logModel := modelLog.NewModel(self.M)
	logModel.Create(context.TODO(), &model.Log{
		ID:             logID,
		Type:           typ,
		UserID:         userId,
		User:           username,
		Method:         method,
		Referer:        referer,
		Path:           path,
		Query:          query,
		IP:             ip,
		RequestHeader:  header,
		RequestPayload: payload,
		DeleteTime:     deleteTime,
	})
}

func (self *Service) logResponse(logID primitive.ObjectID, header map[string][]string, data string) {
	logModel := modelLog.NewModel(self.M)
	filter := logModel.FilterByID(logID)
	setData := bson.D{
		{"response_header", header},
		{"response_data", data},
	}

	logModel.Set(context.TODO(), filter, setData)
}
