package profile

import (
	"context"
	"net/http"

	"github.com/haiyiyun/plugins/urbac/database/model/user"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/response"
	"go.mongodb.org/mongo-driver/bson"
)

func (self *Service) Route_POST_InfoUpdate(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if u, found := self.GetUserInfo(r); found {
		realName := r.FormValue("real_name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		avatar := r.FormValue("avatar")
		description := r.FormValue("description")

		userModel := user.NewModel(self.M)
		filter := userModel.FilterByID(u.ID)
		change := bson.D{}

		if realName != "" {
			change = append(change, bson.E{"real_name", realName})
		}

		if email != "" {
			change = append(change, bson.E{"email", email})
		}

		if password != "" {
			password = help.Strings(password).Md5()
			change = append(change, bson.E{"password", password})
		}

		if avatar != "" {
			change = append(change, bson.E{"avatar", avatar})
		}

		if description != "" {
			change = append(change, bson.E{"description", description})
		}

		if ur, err := userModel.Set(context.Background(), filter, change); err == nil && ur.ModifiedCount > 0 {
			response.JSON(rw, 0, nil, "")
			return
		}
	}

	response.JSON(rw, http.StatusUnauthorized, nil, "")

}
