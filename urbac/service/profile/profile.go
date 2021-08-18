package profile

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/haiyiyun/plugins/urbac/database/model/user"

	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson"
)

func (self *Service) Route_POST_InfoUpdate(rw http.ResponseWriter, r *http.Request) (b bool) {
	r.ParseForm()
	u, _ := self.GetUserInfo(r)
	username := u.Name
	if _, ok := r.Form["ajax"]; ok {
		b = true

		m := help.M{
			"status":  "1",
			"message": "",
		}
		if !ok {
			m["status"] = "0"
			m["message"] = "修改用户时必须指定用户名"
		}

		userModel := user.NewModel(self.M)

		filter := bson.D{
			{"name", username},
			{"delete", false},
		}
		cnt, err := userModel.CountDocuments(context.Background(), filter)
		if err != nil || cnt == 0 {
			m["status"] = "0"
			m["message"] = "用户不存在!"
		} else {
			email := r.FormValue("email")
			password := r.FormValue("password")
			change := bson.D{}

			if email != "" {
				change = append(change, bson.E{"email", email})
			}
			if password != "" {
				password = help.Strings(password).Md5()
				change = append(change, bson.E{"password", password})
			}
			if email != "" || password != "" {
				_, err := userModel.Set(context.Background(), filter, change)
				if err != nil {
					m["status"] = "0"
					m["message"] = "用户资料更新失败"
				} else {
					m["message"] = "用户资料更新成功"
				}
			} else {
				m["status"] = "0"
				m["message"] = "请输入需要更新的资料"
			}
		}
		ret, _ := json.Marshal(m)
		rw.Write(ret)
	}

	return
}
