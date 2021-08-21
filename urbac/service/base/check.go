package base

import (
	"strings"

	"github.com/haiyiyun/plugins/urbac/database/model"

	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) CheckRight(reqPath, httpMethod string, userID primitive.ObjectID) (allow bool) {
	checkRight := true
	if checkMethods, found := self.Config.IgnoreCheckRightPath[reqPath]; found {
		if len(checkMethods) == 0 || help.NewSlice(checkMethods).CheckPartItem(httpMethod, "") {
			checkRight = false
		}
	}

	if !checkRight {
		allow = true
	} else {
		applicationMap := map[string]model.Application{}
		if appTmp, ok := self.GetApplications(userID, true)["map"].(map[string]model.Application); ok {
			applicationMap = appTmp
		}

		reqPaths := strings.Split(reqPath, "/")
		if len(reqPaths) != 4 {
			allow = false
			return
		}

		appPath, modulePath, actionPath := "/"+reqPaths[1]+"/", reqPaths[2]+"/", reqPaths[3]
		if application, found := applicationMap[appPath]; found {
			if module, found := application.Modules[modulePath]; found {
				if action, found := module.Actions[actionPath]; found {
					if action.Enable {
						if action.Method != nil && len(action.Method) > 0 {
							if help.NewSlice(action.Method).CheckPartItem(httpMethod, "") {
								allow = true
								return
							}
						} else {
							allow = true
							return
						}
					}
				}
			}
		}
	}

	return
}
