package application

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/urbac/database/model"
	"github.com/haiyiyun/plugins/urbac/database/model/application"
	"github.com/haiyiyun/plugins/urbac/predefined"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/validator"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Service) GetRoute(appsSorted []help.M, ignoreEnable bool, excludePathisURL, childrenAction, childrenMethod bool) []help.M {
	routeList := []help.M{}
	if len(self.Config.DefaultRoute) > 0 {
		routeList = append(routeList, self.Config.DefaultRoute)
	}
	//app
	for _, application := range appsSorted {
		if excludePathisURL && strings.Index(application["path"].(string), "http") != -1 {
			continue
		}

		appRoute := help.M{
			"_id":                   application["_id"],
			"type":                  application["type"],
			"path":                  application["path"],
			"name":                  application["name"],
			"parent_path":           "",
			"absolute_path":         self.GetRoutePath(application["path"].(string), "", "", ""),
			"level":                 predefined.ApplicationLevelApp,
			"meta":                  application["meta"],
			"default_module_action": application["default_module_action"],
			"order":                 application["order"],
			"enable":                application["enable"],
			"create_time":           application["create_time"],
			"update_time":           application["update_time"],
		}

		//module
		if appModulesI, found := application["modules"]; found {
			if appModules, ok := appModulesI.([]help.M); ok {
				modulesRoute := []help.M{}
				for _, module := range appModules {
					if excludePathisURL && strings.Index(module["path"].(string), "http") != -1 {
						continue
					}

					if enable, ok := module["enable"].(bool); ok && ignoreEnable || enable {
						moduleRoute := help.M{
							"type":          module["type"],
							"path":          module["path"],
							"name":          module["name"],
							"parent_path":   appRoute["absolute_path"],
							"absolute_path": self.GetRoutePath(application["path"].(string), module["path"].(string), "", ""),
							"level":         predefined.ApplicationLevelModule,
							"meta":          module["meta"],
							"order":         module["order"],
							"enable":        module["enable"],
						}

						//action
						if childrenAction {
							if appActionsI, found := module["actions"]; found {
								if appActions, ok := appActionsI.([]help.M); ok {
									actionsRoute := []help.M{}
									for _, action := range appActions {
										if excludePathisURL && strings.Index(action["path"].(string), "http") != -1 {
											continue
										}

										if enable, ok := action["enable"].(bool); ok && ignoreEnable || enable {
											actionRoute := help.M{
												"type":          action["type"],
												"path":          action["path"],
												"name":          action["name"],
												"parent_path":   moduleRoute["absolute_path"],
												"absolute_path": self.GetRoutePath(application["path"].(string), module["path"].(string), action["path"].(string), ""),
												"level":         predefined.ApplicationLevelAction,
												"meta":          action["meta"],
												"order":         action["order"],
												"enable":        action["enable"],
												"method":        action["method"],
											}

											if childrenMethod {
												httpMethod, _ := action["method"].([]string)
												if httpMethod != nil && len(httpMethod) > 0 {
													methodsRoute := []help.M{}
													for _, method := range httpMethod {
														methodRoute := help.M{
															"type":          predefined.ApplicationTypeCode,
															"name":          method,
															"path":          method,
															"parent_path":   actionRoute["absolute_path"],
															"absolute_path": self.GetRoutePath(application["path"].(string), module["path"].(string), action["path"].(string), method),
															"level":         predefined.ApplicationLevelMethod,
														}

														methodsRoute = append(methodsRoute, methodRoute)
													}

													actionRoute["children"] = methodsRoute
												}
											}

											actionsRoute = append(actionsRoute, actionRoute)

										}
									}

									if len(actionsRoute) > 0 {
										moduleRoute["children"] = actionsRoute
									}
								}
							}
						}

						modulesRoute = append(modulesRoute, moduleRoute)
					}

				}

				if len(modulesRoute) > 0 {
					appRoute["children"] = modulesRoute
				}
			}
		}

		routeList = append(routeList, appRoute)
	}

	return routeList
}

func (self *Service) Route_GET_RouteList(rw http.ResponseWriter, r *http.Request) {
	appModel := application.NewModel(self.M)
	filter := bson.D{}
	opt := options.Find().SetSort(bson.D{
		{"order", 1},
		{"create_time", -1},
	})

	apps := map[string]model.Application{}
	ctx := r.Context()
	if cur, err := appModel.Find(ctx, filter, opt); err != nil {
		log.Error(err)
	} else {
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var application model.Application
			if err = cur.Decode(&application); err != nil {
				log.Error(err)
				break
			} else {
				apps[application.Path] = application
			}
		}
	}

	sortApp := self.SortApps(apps)
	routeApp := self.GetRoute(sortApp, false, false, true, true)

	response.JSON(rw, 0, routeApp, "")
}

func (self *Service) Route_GET_Index(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	childrenAction := true
	childrenActionStr := r.FormValue("children_action")
	if ca, err := strconv.ParseBool(childrenActionStr); err == nil {
		childrenAction = ca
	}

	excludePathisURL := false
	excludePathisURLStr := r.FormValue("exclude_path_is_url")
	if epu, err := strconv.ParseBool(excludePathisURLStr); err == nil {
		excludePathisURL = epu
	}

	appModel := application.NewModel(self.M)
	filter := bson.D{}

	Types := r.Form["type[]"]
	if len(Types) > 0 {
		filter = append(filter, bson.E{
			"type", bson.D{
				{"$in", Types},
			},
		})
	}

	if name := r.FormValue("name"); name != "" {
		filter = append(filter, bson.E{"name", name})
	}

	if path := r.FormValue("path"); path != "" {
		filter = append(filter, bson.E{"path", path})
	}

	cnt, _ := appModel.CountDocuments(context.Background(), filter)
	pg := pagination.Parse(r, cnt)

	opt := options.Find().SetSort(bson.D{
		{"order", 1},
		{"create_time", -1},
	}).SetProjection(bson.D{}).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	apps := map[string]model.Application{}
	ctx := r.Context()
	if cur, err := appModel.Find(ctx, filter, opt); err != nil {
		log.Error(err)
	} else {
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var application model.Application
			if err = cur.Decode(&application); err != nil {
				log.Error(err)
				break
			} else {
				apps[application.Path] = application
			}
		}
	}

	sortApp := self.SortApps(apps)
	routeApp := self.GetRoute(sortApp, true, excludePathisURL, childrenAction, true)

	rpr := response.ResponsePaginationResult{
		Total: cnt,
		Items: routeApp,
	}

	response.JSON(rw, 0, rpr, "")
}

func (self *Service) Route_POST_Update(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	parentPath := r.FormValue("parent_path")
	typ := r.FormValue("type")
	name := r.FormValue("name")
	pathTmp := r.FormValue("path")
	level := r.FormValue("level")
	orderStr := r.FormValue("order")
	enableStr := r.FormValue("enable")
	metaHideMenuStr := r.FormValue("meta_hide_menu")
	metaTitle := r.FormValue("meta_title")
	metaIcon := r.FormValue("meta_icon")
	metaFrameSrc := r.FormValue("meta_frame_src")

	valid := &validator.Validation{}
	valid.Required(typ).Key("type").Message("type字段为空")
	valid.Required(name).Key("name").Message("name字段为空")
	valid.Required(pathTmp).Key("path").Message("path字段为空")
	valid.Required(level).Key("level").Message("level字段为空")
	valid.Required(enableStr).Key("enable").Message("enable字段为空")
	valid.Required(metaTitle).Key("meta_title").Message("meta_title字段为空")
	valid.Required(orderStr).Key("order").Message("order字段为空")
	valid.Required(metaHideMenuStr).Key("meta_hide_menu").Message("meta_hide_menu字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	aPath := []string{}
	isURL := false

	if pos := strings.Index(pathTmp, "http"); pos != -1 {
		isURL = true
		aPath = strings.Split(parentPath, "/")
		aPath = append(aPath, pathTmp)
	} else {
		aPath = strings.Split(parentPath+pathTmp, "/")
	}

	if len(aPath) < 4 {
		if level != predefined.ApplicationLevelApp {
			response.JSON(rw, http.StatusBadRequest, nil, "")
			return
		}
	}

	appPath := "/" + aPath[1] + "/"
	if isURL && level == predefined.ApplicationLevelApp {
		appPath = aPath[1]
	}

	appModel := application.NewModel(self.M)
	filter := bson.D{
		{"path", appPath},
	}

	enable, _ := strconv.ParseBool(enableStr)
	order, _ := strconv.Atoi(orderStr)
	metaHideMenu, _ := strconv.ParseBool(metaHideMenuStr)
	if typ == predefined.ApplicationTypeCode {
		metaFrameSrc = ""
	}

	switch level {
	case predefined.ApplicationLevelApp:
		change := bson.D{
			{"name", name},
			{"enable", enable},
			{"order", order},
			{"meta.hide_menu", metaHideMenu},
			{"meta.title", metaTitle},
			{"meta.icon", metaIcon},
			{"meta.frame_src", metaFrameSrc},
		}
		if _, err := appModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		} else {
			log.Debug(err)
		}
	case predefined.ApplicationLevelModule:
		modulePath := aPath[2] + "/"
		if isURL {
			modulePath = aPath[3]
		}

		moduleField := "modules." + modulePath
		change := bson.D{
			{moduleField + ".name", name},
			{moduleField + ".enable", enable},
			{moduleField + ".order", order},
			{moduleField + ".meta.hide_menu", metaHideMenu},
			{moduleField + ".meta.title", metaTitle},
			{moduleField + ".meta.icon", metaIcon},
			{moduleField + ".meta.frame_src", metaFrameSrc},
		}

		if _, err := appModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		} else {
			log.Debug(err)
		}
	case predefined.ApplicationLevelAction:
		modulePath := aPath[2] + "/"
		actionPath := aPath[3]
		if isURL {
			actionPath = aPath[4]
		}

		actionField := "modules." + modulePath + ".actions." + actionPath
		change := bson.D{
			{actionField + ".name", name},
			{actionField + ".enable", enable},
			{actionField + ".order", order},
			{actionField + ".meta.hide_menu", metaHideMenu},
			{actionField + ".meta.title", metaTitle},
			{actionField + ".meta.icon", metaIcon},
			{actionField + ".meta.frame_src", metaFrameSrc},
		}

		if _, err := appModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		} else {
			log.Debug(err)
		}
	default:
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
	return
}

func (self *Service) Route_PUT_CreateVirtualApplication(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	parentPath := r.FormValue("parent_path")
	typ := r.FormValue("type")
	name := r.FormValue("name")
	pathTmp := r.FormValue("path")
	level := r.FormValue("level")
	orderStr := r.FormValue("order")
	enableStr := r.FormValue("enable")
	metaHideMenuStr := r.FormValue("meta_hide_menu")
	metaTitle := r.FormValue("meta_title")
	metaIcon := r.FormValue("meta_icon")
	metaFrameSrc := r.FormValue("meta_frame_src")

	valid := &validator.Validation{}
	valid.Required(typ).Key("type").Message("type字段为空")
	valid.Required(name).Key("name").Message("name字段为空")
	valid.Required(pathTmp).Key("path").Message("path字段为空")
	valid.Required(level).Key("level").Message("level字段为空")
	valid.Required(enableStr).Key("enable").Message("enable字段为空")
	valid.Required(metaTitle).Key("meta_title").Message("meta_title字段为空")
	valid.Required(orderStr).Key("order").Message("order字段为空")
	valid.Required(metaHideMenuStr).Key("meta_hide_menu").Message("meta_hide_menu字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	if pos := strings.Index(parentPath, "http"); pos != -1 {
		response.JSON(rw, http.StatusBadRequest, nil, "不能在外链下建立虚拟应用")
		return
	}

	aPath := []string{}
	isURL := false

	if pos := strings.Index(pathTmp, "http"); pos == 0 {
		isURL = true
		aPath = strings.Split(parentPath, "/")
		aPath = append(aPath, pathTmp)
	} else {
		aPath = strings.Split(parentPath+pathTmp, "/")
	}

	if len(aPath) < 4 {
		if level != predefined.ApplicationLevelApp {
			response.JSON(rw, http.StatusBadRequest, nil, "")
			return
		}
	}

	appPath := "/" + aPath[1] + "/"
	if isURL && level == predefined.ApplicationLevelApp {
		appPath = aPath[1]
	}

	appModel := application.NewModel(self.M)
	filter := bson.D{
		{"path", appPath},
	}

	enable, _ := strconv.ParseBool(enableStr)
	order, _ := strconv.Atoi(orderStr)
	metaHideMenu, _ := strconv.ParseBool(metaHideMenuStr)
	if typ == predefined.ApplicationTypeCode {
		metaFrameSrc = ""
	}

	switch level {
	case predefined.ApplicationLevelApp:
		change := bson.D{
			{"type", predefined.ApplicationTypeVirtual},
			{"name", name},
			{"path", appPath},
			{"enable", enable},
			{"order", order},
			{"meta.hide_menu", metaHideMenu},
			{"meta.title", metaTitle},
			{"meta.icon", metaIcon},
			{"meta.frame_src", metaFrameSrc},
		}
		if _, err := appModel.SetAndSetOnInsert(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		} else {
			log.Debug(err)
		}
	case predefined.ApplicationLevelModule:
		modulePath := aPath[2] + "/"
		if isURL {
			modulePath = aPath[3]
		}

		moduleField := "modules." + modulePath
		change := bson.D{
			{moduleField + ".type", predefined.ApplicationTypeVirtual},
			{moduleField + ".name", name},
			{moduleField + ".path", modulePath},
			{moduleField + ".enable", enable},
			{moduleField + ".order", order},
			{moduleField + ".meta.hide_menu", metaHideMenu},
			{moduleField + ".meta.title", metaTitle},
			{moduleField + ".meta.icon", metaIcon},
			{moduleField + ".meta.frame_src", metaFrameSrc},
		}

		if _, err := appModel.SetAndSetOnInsert(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		} else {
			log.Debug(err)
		}
	case predefined.ApplicationLevelAction:
		modulePath := aPath[2] + "/"
		actionPath := aPath[3]
		if isURL {
			actionPath = aPath[4]
		}

		actionField := "modules." + modulePath + ".actions." + actionPath
		change := bson.D{
			{actionField + ".type", predefined.ApplicationTypeVirtual},
			{actionField + ".name", name},
			{actionField + ".path", actionPath},
			{actionField + ".enable", enable},
			{actionField + ".order", order},
			{actionField + ".meta.hide_menu", metaHideMenu},
			{actionField + ".meta.title", metaTitle},
			{actionField + ".meta.icon", metaIcon},
			{actionField + ".meta.frame_src", metaFrameSrc},
		}

		if _, err := appModel.SetAndSetOnInsert(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		} else {
			log.Debug(err)
		}
	default:
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
	return
}

func (self *Service) Route_PUT_CreateCodeApplication(rw http.ResponseWriter, r *http.Request) {
	responseMsg := "创建应用失败"

	applicationName := r.FormValue("name")
	applicationPath := r.FormValue("path")

	valid := &validator.Validation{}
	valid.Required(applicationName).Key("name").Message("name字段为空")
	valid.Required(applicationPath).Key("path").Message("path字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	if applicationPath[0] == '/' && applicationPath[len(applicationPath)-1] == '/' {
		appModel := application.NewModel(self.M)

		filter := bson.D{
			{"$or", []bson.D{
				{{"path", applicationPath}},
				{{"name", applicationName}},
			}},
		}

		cnt, err := appModel.CountDocuments(context.Background(), filter)
		if err != nil {
			responseMsg = "获取应用信息失败"
		} else if cnt != 0 {
			responseMsg = "存在相同的应用名"
		} else {
			regroutes := webrouter.Registers()
			var foundApp bool
			app := model.Application{
				Type: predefined.ApplicationTypeCode,
				Name: applicationName,
				Path: applicationPath,
				Meta: model.ApplicationMeta{
					Title: applicationName,
				},
				Enable:  self.Config.DefaultEnableApp,
				Modules: map[string]model.ApplicationModule{},
			}

			for rPath, rI := range regroutes {
				if !help.NewSlice(self.Config.IgnoreAppModuleInfo).CheckItem(rPath) {
					if strings.HasPrefix(rPath, applicationPath) {
						if !foundApp {
							foundApp = true
						}
						modPath := rPath[len(applicationPath):]

						app.Modules[modPath] = model.ApplicationModule{
							Type:   predefined.ApplicationTypeCode,
							Name:   modPath,
							Path:   modPath,
							Enable: self.Config.DefaultEnableModule,
							Meta: model.ApplicationMeta{
								Title: modPath,
							},
							Actions: map[string]model.ApplicationModuleAction{},
						}

						rt := rI.Type
						for j := 0; j < rt.NumMethod(); j++ {
							if filterMN, httpMN := webrouter.GetFilterMethodNameAndHTTPMethodName(rt.Method(j).Name); filterMN != "" {
								mName := webrouter.MakePattern(filterMN)

								if ama, found := app.Modules[modPath].Actions[mName]; !found {
									var actOrder int64
									if mName == "index" {
										actOrder = -1
									}

									app.Modules[modPath].Actions[mName] = model.ApplicationModuleAction{
										Type:   predefined.ApplicationTypeCode,
										Name:   mName,
										Path:   mName,
										Order:  actOrder,
										Enable: self.Config.DefaultEnableAction,
										Meta: model.ApplicationMeta{
											Title: mName,
										},
										Method: []string{httpMN},
									}
								} else {
									ama.Method = append(ama.Method, httpMN)
									app.Modules[modPath].Actions[mName] = ama
								}
							}
						}
					}
				}
			}

			if foundApp {
				if _, err := appModel.Create(context.Background(), &app); err != nil {
					log.Error(err)
				} else {
					response.JSON(rw, 0, nil, "")
					return
				}
			}
		}
	}

	response.JSON(rw, http.StatusBadRequest, nil, responseMsg)
}

func (self *Service) Route_POST_ReadCodeApplication(rw http.ResponseWriter, r *http.Request) {
	responseMsg := "获取应用信息失败"

	applicationName := r.FormValue("name")
	applicationPath := r.FormValue("path")

	valid := &validator.Validation{}
	valid.Required(applicationName).Key("name").Message("name字段为空")
	valid.Required(applicationPath).Key("path").Message("path字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	appModel := application.NewModel(self.M)
	filter := bson.D{
		{"path", applicationPath},
		{"name", applicationName},
	}

	sr := appModel.FindOne(context.Background(), filter)
	if sr.Err() == nil {
		app := model.Application{}
		if err := sr.Decode(&app); err == nil {
			//取反射
			regroutes := webrouter.Registers()
			//初始化
			aModules := app.Modules
			modules := map[string]model.ApplicationModule{}
			//循环添加虚拟模块
			for aMPath, aModule := range aModules {
				if aModule.Type == predefined.ApplicationTypeVirtual {
					modules[aMPath] = aModule
				}
			}
			var foundApp bool
			//循环反射应用
			for rPath, rI := range regroutes {
				//过滤
				if !help.NewSlice(self.Config.IgnoreAppModuleInfo).CheckItem(rPath) {
					//判断是否含有当前应用
					if strings.HasPrefix(rPath, applicationPath) {
						if foundApp == false {
							foundApp = true
						}

						mPath := rPath[len(applicationPath):]
						//判断反射出的应用模块是否包含当前应用已经存在的模块
						if mod, found := aModules[mPath]; found {
							modules[mPath] = model.ApplicationModule{
								Type:    mod.Type,
								Name:    mod.Name,
								Enable:  mod.Enable,
								Path:    mod.Path,
								Order:   mod.Order,
								Meta:    mod.Meta,
								Actions: map[string]model.ApplicationModuleAction{},
							}
							//循环添加虚拟方法
							for aAPath, aAction := range mod.Actions {
								if aAction.Type == predefined.ApplicationTypeVirtual {
									modules[mPath].Actions[aAPath] = aAction
								}
							}

							rt := rI.Type
							//循环反射已存在模块下的方法
							for j := 0; j < rt.NumMethod(); j++ {
								//判断是否已经存在的方法
								if filterMN, httpMN := webrouter.GetFilterMethodNameAndHTTPMethodName(rt.Method(j).Name); filterMN != "" {
									aPath := webrouter.MakePattern(filterMN)
									aAction, foundModAct := mod.Actions[aPath]
									ama, foundAct := modules[mPath].Actions[aPath]

									if foundModAct && !foundAct {
										modules[mPath].Actions[aPath] = model.ApplicationModuleAction{
											Type:   aAction.Type,
											Name:   aAction.Name,
											Enable: aAction.Enable,
											Path:   aAction.Path,
											Order:  aAction.Order,
											Meta:   aAction.Meta,
											Method: []string{httpMN},
										}
									} else {
										if !foundAct {
											var actOrder int64
											if aPath == "index" {
												actOrder = -1
											}

											modules[mPath].Actions[aPath] = model.ApplicationModuleAction{
												Type:   predefined.ApplicationTypeCode,
												Name:   aPath,
												Enable: self.Config.DefaultEnableAction,
												Path:   aPath,
												Order:  actOrder,
												Meta: model.ApplicationMeta{
													Title: aPath,
												},
												Method: []string{httpMN},
											}
										} else {
											ama.Method = append(ama.Method, httpMN)
											modules[mPath].Actions[aPath] = ama
										}
									}
								}
							}
						} else {
							//不包含已经存在的模块，等于发现新模块
							modules[mPath] = model.ApplicationModule{
								Type:   predefined.ApplicationTypeCode,
								Name:   mPath,
								Path:   mPath,
								Enable: self.Config.DefaultEnableModule,
								Meta: model.ApplicationMeta{
									Title: mPath,
								},
								Actions: map[string]model.ApplicationModuleAction{},
							}
							rt := rI.Type
							for j := 0; j < rt.NumMethod(); j++ {
								if filterMN, httpMN := webrouter.GetFilterMethodNameAndHTTPMethodName(rt.Method(j).Name); filterMN != "" {
									aPath := webrouter.MakePattern(filterMN)

									if ama, found := modules[mPath].Actions[aPath]; !found {
										var actOrder int64
										if aPath == "index" {
											actOrder = -1
										}

										modules[mPath].Actions[aPath] = model.ApplicationModuleAction{
											Type:   predefined.ApplicationTypeCode,
											Name:   aPath,
											Path:   aPath,
											Order:  actOrder,
											Enable: self.Config.DefaultEnableAction,
											Meta: model.ApplicationMeta{
												Title: aPath,
											},

											Method: []string{httpMN},
										}
									} else {
										ama.Method = append(ama.Method, httpMN)
										modules[mPath].Actions[aPath] = ama
									}
								}
							}
						}
					}
				}
			}

			if foundApp {
				if _, err := appModel.Set(context.Background(), filter, bson.D{
					{"modules", modules},
				}); err != nil {
					responseMsg = "读取应用信息失败"
				} else {
					response.JSON(rw, 0, nil, "")
					return
				}
			} else {
				if dr, err := appModel.DeleteOne(context.Background(), filter); dr.DeletedCount > 0 && err == nil {
					response.JSON(rw, 0, nil, "")
					return
				}
			}
		}

	}

	response.JSON(rw, http.StatusBadRequest, nil, responseMsg)

	return

}

func (self *Service) Route_DELETE_Delete(rw http.ResponseWriter, r *http.Request) {
	values, _ := request.ParseDeleteForm(r)

	pathTmp := values.Get("path")
	level := values.Get("level")

	valid := &validator.Validation{}
	valid.Required(pathTmp).Key("path").Message("path字段为空")
	valid.Required(level).Key("level").Message("level字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	aPath := []string{}
	isURL := false

	if pos := strings.Index(pathTmp, "http"); pos != -1 {
		isURL = true
		link := pathTmp[pos:]
		aPath = strings.Split(pathTmp[0:pos], "/")
		aPath = append(aPath, link)
	} else {
		aPath = strings.Split(pathTmp, "/")
	}

	if len(aPath) < 4 {
		if level != predefined.ApplicationLevelApp {
			response.JSON(rw, http.StatusBadRequest, nil, "")
			return
		}
	}

	appPath := "/" + aPath[1] + "/"
	if isURL && level == predefined.ApplicationLevelApp {
		appPath = aPath[1]
	}

	appModel := application.NewModel(self.M)
	filter := bson.D{
		{"path", appPath},
	}

	switch level {
	case predefined.ApplicationLevelApp:
		if dr, err := appModel.DeleteOne(context.Background(), filter); dr.DeletedCount > 0 && err == nil {
			response.JSON(rw, 0, nil, "")
			return
		}
	case predefined.ApplicationLevelModule:
		modulePath := aPath[2] + "/"
		if isURL {
			modulePath = aPath[3]
		}

		moduleField := "modules." + modulePath

		unset := bson.D{
			{moduleField, ""},
		}

		if _, err := appModel.UnSet(r.Context(), filter, unset); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		}
	case predefined.ApplicationLevelAction:
		modulePath := aPath[2] + "/"
		actionPath := aPath[3]
		if isURL {
			actionPath = aPath[4]
		}

		actionField := "modules." + modulePath + ".actions." + actionPath

		unset := bson.D{
			{actionField, ""},
		}

		if _, err := appModel.UnSet(r.Context(), filter, unset); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		}
	case predefined.ApplicationLevelMethod:
		modulePath := aPath[2] + "/"
		actionPath := aPath[3]
		methodField := "modules." + modulePath + ".actions." + actionPath + ".method"

		if pos := strings.Index(aPath[0], "_"); pos != -1 {
			httpMethod := aPath[0][0:pos]
			pull := bson.D{
				{methodField, httpMethod},
			}
			if _, err := appModel.Pull(r.Context(), filter, pull); err == nil {
				response.JSON(rw, 0, nil, "")
				return
			}
		}
	default:
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
	return
}

func (self *Service) Route_POST_Enable(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	pathTmp := r.FormValue("path")
	level := r.FormValue("level")
	enableStr := r.FormValue("enable")

	valid := &validator.Validation{}
	valid.Required(pathTmp).Key("path").Message("path字段为空")
	valid.Required(level).Key("level").Message("level字段为空")
	valid.Required(enableStr).Key("enable").Message("enable字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	aPath := []string{}
	isURL := false

	if pos := strings.Index(pathTmp, "http"); pos != -1 {
		isURL = true
		link := pathTmp[pos:]
		aPath = strings.Split(pathTmp[0:pos], "/")
		aPath = append(aPath, link)
	} else {
		aPath = strings.Split(pathTmp, "/")
	}

	if len(aPath) < 4 {
		if level != predefined.ApplicationLevelApp {
			response.JSON(rw, http.StatusBadRequest, nil, "")
			return
		}
	}

	appPath := "/" + aPath[1] + "/"
	if isURL && level == predefined.ApplicationLevelApp {
		appPath = aPath[1]
	}

	appModel := application.NewModel(self.M)
	filter := bson.D{
		{"path", appPath},
	}

	enable, _ := strconv.ParseBool(enableStr)

	switch level {
	case predefined.ApplicationLevelApp:
		change := bson.D{
			{"enable", enable},
		}
		if _, err := appModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		}
	case predefined.ApplicationLevelModule:
		modulePath := aPath[2] + "/"
		if isURL {
			modulePath = aPath[3]
		}

		moduleField := "modules." + modulePath

		change := bson.D{
			{moduleField + ".enable", enable},
		}

		if _, err := appModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		}
	case predefined.ApplicationLevelAction:
		modulePath := aPath[2] + "/"
		actionPath := aPath[3]
		if isURL {
			actionPath = aPath[4]
		}

		actionField := "modules." + modulePath + ".actions." + actionPath

		change := bson.D{
			{actionField + ".enable", enable},
		}
		if _, err := appModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		}
	default:
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
	return
}

func (self *Service) Route_POST_Hide(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	pathTmp := r.FormValue("path")
	level := r.FormValue("level")
	hideStr := r.FormValue("hide")

	valid := &validator.Validation{}
	valid.Required(pathTmp).Key("path").Message("path字段为空")
	valid.Required(level).Key("level").Message("level字段为空")
	valid.Required(hideStr).Key("hide").Message("hide字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	aPath := []string{}
	isURL := false

	if pos := strings.Index(pathTmp, "http"); pos != -1 {
		isURL = true
		link := pathTmp[pos:]
		aPath = strings.Split(pathTmp[0:pos], "/")
		aPath = append(aPath, link)
	} else {
		aPath = strings.Split(pathTmp, "/")
	}

	if len(aPath) < 4 {
		if level != predefined.ApplicationLevelApp {
			response.JSON(rw, http.StatusBadRequest, nil, "")
			return
		}
	}

	appPath := "/" + aPath[1] + "/"
	if isURL && level == predefined.ApplicationLevelApp {
		appPath = aPath[1]
	}

	appModel := application.NewModel(self.M)
	filter := bson.D{
		{"path", appPath},
	}

	hide, _ := strconv.ParseBool(hideStr)

	switch level {
	case predefined.ApplicationLevelApp:
		change := bson.D{
			{"meta.hide_menu", hide},
		}
		if _, err := appModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		}
	case predefined.ApplicationLevelModule:
		modulePath := aPath[2] + "/"
		if isURL {
			modulePath = aPath[3]
		}

		moduleField := "modules." + modulePath

		change := bson.D{
			{moduleField + ".meta.hide_menu", hide},
		}

		if _, err := appModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		}
	case predefined.ApplicationLevelAction:
		modulePath := aPath[2] + "/"
		actionPath := aPath[3]
		if isURL {
			actionPath = aPath[4]
		}

		actionField := "modules." + modulePath + ".actions." + actionPath

		change := bson.D{
			{actionField + ".meta.hide_menu", hide},
		}
		if _, err := appModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
			return
		}
	default:
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
	return
}
