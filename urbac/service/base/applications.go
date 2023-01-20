package base

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/haiyiyun/plugins/urbac/database/model"
	"github.com/haiyiyun/plugins/urbac/database/model/application"
	"github.com/haiyiyun/plugins/urbac/predefined"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Service) GetApplications(userID primitive.ObjectID, checkRole bool) help.M {
	cacheKey := "applications.info." + userID.Hex()
	if appI, found := self.Cache.Get(cacheKey); found {
		appMapIs := appI.(help.M)
		return appMapIs
	}

	apps := map[string]model.Application{}
	if roles, err := self.getRole(userID); err != nil {
		log.Error("getRole error:", err)
	} else {
		filter := bson.D{
			{"enable", true},
		}

		opt := options.Find().SetSort(bson.D{
			{"order", 1},
			{"create_time", -1},
		})

		appModel := application.NewModel(self.M)
		ctx := context.Background()
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
					if !checkRole {
						apps[application.Path] = application
					} else {
						for _, role := range roles {
							switch role.Right.Scope {
							case predefined.RoleRightTypeAction:
								if rapp, found := role.Right.Applications[application.Path]; found {
									if aApp, found := apps[application.Path]; !found {
										appTmp := model.Application{}
										help.DeepCopy(application, &appTmp)
										appTmp.Modules = map[string]model.ApplicationModule{}
										for rmPath, rmValue := range rapp.Modules {
											moduleTmp := model.ApplicationModule{}
											help.DeepCopy(application.Modules[rmPath], &moduleTmp)
											moduleTmp.Actions = map[string]model.ApplicationModuleAction{}
											for raPath, raValue := range rmValue.Actions {
												actionTmp := model.ApplicationModuleAction{}
												help.DeepCopy(application.Modules[rmPath].Actions[raPath], &actionTmp)
												actionTmp.Path = raValue.Path
												actionTmp.Method = raValue.Method

												if actionTmp.Enable {
													moduleTmp.Actions[raPath] = actionTmp
												}
											}
											appTmp.Modules[rmPath] = moduleTmp
										}
										apps[application.Path] = appTmp

									} else {
										for rmPath, rmValue := range rapp.Modules {
											if amValue, found := aApp.Modules[rmPath]; !found {
												moduleTmp := model.ApplicationModule{}
												help.DeepCopy(application.Modules[rmPath], &moduleTmp)
												moduleTmp.Actions = map[string]model.ApplicationModuleAction{}
												for raPath, raValue := range rmValue.Actions {
													actionTmp := model.ApplicationModuleAction{}
													help.DeepCopy(application.Modules[rmPath].Actions[raPath], &actionTmp)
													actionTmp.Path = raValue.Path
													actionTmp.Method = raValue.Method

													if actionTmp.Enable {
														moduleTmp.Actions[raPath] = actionTmp
													}
												}
												aApp.Modules[rmPath] = moduleTmp
											} else {
												for raPath, raValue := range rmValue.Actions {
													actionTmp := model.ApplicationModuleAction{}
													help.DeepCopy(application.Modules[rmPath].Actions[raPath], &actionTmp)
													actionTmp.Path = raValue.Path
													actionTmp.Method = raValue.Method

													if actionTmp.Enable {
														amValue.Actions[raPath] = actionTmp
													}
												}
											}
										}
									}
								}
							case predefined.RoleRightTypeModule:
								if rapp, found := role.Right.Applications[application.Path]; found {
									if aApp, found := apps[application.Path]; !found {
										appModuleTmp := map[string]model.ApplicationModule{}
										for amPath, _ := range rapp.Modules {
											appModuleTmp[amPath] = application.Modules[amPath]
										}

										appTmp := model.Application{}
										help.DeepCopy(application, &appTmp)
										appTmp.Modules = appModuleTmp
										apps[application.Path] = appTmp
									} else {
										for rmPath := range rapp.Modules {
											if _, found := aApp.Modules[rmPath]; !found {
												aApp.Modules[rmPath] = application.Modules[rmPath]
											}
										}
									}
								}
							case predefined.RoleRightTypeApp:
								if _, found := role.Right.Applications[application.Path]; found {
									if _, found := apps[application.Path]; !found {
										apps[application.Path] = application
									}
								}
							case predefined.RoleRightTypePlatform:
								if _, found := apps[application.Path]; !found {
									apps[application.Path] = application
								}
							}
						}
					}
				}
			}
		}
	}

	appsSorted := self.SortApps(apps)
	permissionCode := self.GetPermissionCode(appsSorted)
	appRoute := self.GetRoute(appsSorted)
	appMapIs := help.M{
		"map":             apps,
		"sorted":          appsSorted,
		"permission_code": permissionCode,
		"route":           appRoute,
	}

	self.Cache.Set(cacheKey, appMapIs, self.Config.TokenExpireDuration.Duration)

	return appMapIs
}

// 自然顺序排序，值越大排越后面
func (self *Service) SortApps(apps map[string]model.Application) []help.M {
	aApps := []help.M{}
	for _, va := range apps {
		ma := help.NewStruct(va).StructToMap()
		aMods := []help.M{}
		for _, vm := range va.Modules {
			mm := help.NewStruct(vm).StructToMap()
			aActs := []help.M{}
			for _, vac := range vm.Actions {
				mac := help.NewStruct(vac).StructToMap()
				aActs = append(aActs, mac)
			}

			sort.Slice(aActs, func(i, j int) bool {
				return aActs[i]["order"].(int64) < aActs[j]["order"].(int64)
			})
			mm["actions"] = aActs
			aMods = append(aMods, mm)
		}

		sort.Slice(aMods, func(i, j int) bool {
			return aMods[i]["order"].(int64) < aMods[j]["order"].(int64)
		})
		ma["modules"] = aMods

		aApps = append(aApps, ma)

		sort.Slice(aApps, func(i, j int) bool {
			return aApps[i]["order"].(int64) < aApps[j]["order"].(int64)
		})
	}

	return aApps
}

// 前端通过name进行路由加载，所以name必须唯一，故在添加virtual application时，需要保证同级下的name必须唯一
func (self *Service) GetRoute(appsSorted []help.M) []help.M {
	routeList := []help.M{}
	if len(self.Config.DefaultRoute) > 0 {
		routeList = append(routeList, self.Config.DefaultRoute)
	}
	//app
	for _, application := range appsSorted {
		appRoute := help.M{
			"path":      application["path"],
			"name":      self.GetRoutePath(application["path"].(string), "", "", ""),
			"component": "LAYOUT",
			"meta":      self.transformMeta(application["meta"]),
			"icon":      application["meta"].(help.M)["icon"],
		}

		if defaultModuleAction, ok := application["default_module_action"].(string); ok && defaultModuleAction != "" {
			appRoute["redirect"] = application["path"].(string) + defaultModuleAction
		}

		if appPath, ok := appRoute["path"].(string); ok {
			if strings.HasPrefix(appPath, `http://`) || strings.HasPrefix(appPath, `https://`) {
				appRoute["component"] = "IFrame"
			} else {
				appRoute["path"] = strings.TrimRight(appPath, "/")
			}
		}

		if meta, ok := appRoute["meta"].(help.M); ok {
			if frI, found := meta["frame_src"]; found {
				if frs, ok := frI.(string); ok && frs != "" {
					appRoute["component"] = "IFrame"
				}
			}
		}

		//module
		if appModulesI, found := application["modules"]; found {
			if appModules, ok := appModulesI.([]help.M); ok {
				modulesRoute := []help.M{}
				for _, module := range appModules {
					if enable, ok := module["enable"].(bool); ok && enable {
						moduleRoute := help.M{
							"path": module["path"],
							"name": self.GetRoutePath(application["path"].(string), module["path"].(string), "", ""),
							"meta": self.transformMeta(module["meta"]),
							"icon": module["meta"].(help.M)["icon"],
						}

						if modulePath, ok := moduleRoute["path"].(string); ok {
							if strings.HasPrefix(modulePath, `http://`) || strings.HasPrefix(modulePath, `https://`) {
								moduleRoute["component"] = "IFrame"
							} else {
								moduleRoute["path"] = strings.TrimRight(modulePath, "/")
							}
						}

						if meta, ok := module["meta"].(help.M); ok {
							if frI, found := meta["frame_src"]; found {
								if frs, ok := frI.(string); ok && frs != "" {
									delete(module, "component")
								}
							}
						}

						//action
						if appActionsI, found := module["actions"]; found {
							if appActions, ok := appActionsI.([]help.M); ok {
								actionsRoute := []help.M{}
								for _, action := range appActions {
									if enable, ok := action["enable"].(bool); ok && enable {
										httpMethod, _ := action["method"].([]string)
										if httpMethod != nil && len(httpMethod) > 0 {
											if !help.NewSlice(httpMethod).CheckPartItem("GET", "") {
												continue
											}
										}

										//如果action隐藏菜单就不提供前端路由
										if action["meta"].(help.M)["hide_menu"].(bool) == true {
											continue
										}

										if aPath, ok := action["path"].(string); ok {
											//设置app级别的菜单默认跳转
											if aPath == "app_index" {
												if redirect, found := appRoute["redirect"]; !found || redirect == "" {
													appRoute["redirect"] = self.GetRoutePath(application["path"].(string), module["path"].(string), action["path"].(string), "")
												}
											}
										}
										actionRoute := help.M{
											"path":      action["path"],
											"name":      self.GetRoutePath(application["path"].(string), module["path"].(string), action["path"].(string), ""),
											"component": self.GetRoutePath(application["path"].(string), module["path"].(string), action["path"].(string), ""),
											"meta":      self.transformMeta(action["meta"]),
											"icon":      action["meta"].(help.M)["icon"],
										}

										if meta, ok := action["meta"].(help.M); ok {
											if frI, found := meta["frame_src"]; found {
												if frs, ok := frI.(string); ok && frs != "" {
													delete(actionRoute, "component")
												}
											}
										}

										actionsRoute = append(actionsRoute, actionRoute)

									}
								}

								if len(actionsRoute) > 0 {
									moduleRoute["children"] = actionsRoute
								} else {
									delete(moduleRoute, "children")
								}
							}
						}

						modulesRoute = append(modulesRoute, moduleRoute)
					}

				}

				if len(modulesRoute) > 0 {
					appRoute["children"] = modulesRoute
				} else {
					delete(appRoute, "children")
				}
			}
		}

		appendRoute := false
		if _, found := appRoute["children"]; found {
			appendRoute = true
		} else {
			if application["type"].(string) != predefined.ApplicationTypeCode {
				appendRoute = true
			}
		}

		if appendRoute {
			routeList = append(routeList, appRoute)
		}
	}

	return routeList
}

func (self *Service) GetPermissionCode(appsSorted []help.M) []string {
	permissionCode := []string{}

	//app
	for _, application := range appsSorted {
		//module
		if appModulesI, found := application["modules"]; found {
			if appModules, ok := appModulesI.([]help.M); ok {
				for _, module := range appModules {
					if enable, ok := module["enable"].(bool); ok && enable {
						//action
						if appActionsI, found := module["actions"]; found {
							if appActions, ok := appActionsI.([]help.M); ok {
								for _, action := range appActions {
									if enable, ok := action["enable"].(bool); ok && enable {
										httpMethods, _ := action["method"].([]string)
										if httpMethods != nil && len(httpMethods) > 0 {
											for _, httpMethod := range httpMethods {
												permissionCode = append(permissionCode, self.GetRoutePath(application["path"].(string), module["path"].(string), action["path"].(string), httpMethod))
											}
										} else {
											permissionCode = append(permissionCode, self.GetRoutePath(application["path"].(string), module["path"].(string), action["path"].(string), ""))
										}
									}
								}
							}
						}
					}

				}
			}
		}
	}

	return permissionCode
}

func (self *Service) GetRoutePath(appPath, modulePath, actionPath, method string) string {
	s := fmt.Sprintf("%s%s%s", appPath, modulePath, actionPath)
	if method != "" {
		s = method + "_" + s
	}

	return s
}

func (self *Service) transformMeta(metaI interface{}) interface{} {
	if meta, ok := metaI.(help.M); ok {
		newMeta := help.M{}

		for k, v := range meta {
			kk := k
			switch kk {
			case "frame_src":
				kk = "frameSrc"
			case "hide_breadcrumb":
				kk = "hideBreadcrumb"
			case "carry_param":
				kk = "carryParam"
			case "hide_children_in_menu":
				kk = "hideChildrenInMenu"
			case "hide_tab":
				kk = "hideTab"
			case "hide_menu":
				kk = "hideMenu"
			}

			newMeta[kk] = v
		}

		return newMeta
	} else {
		return metaI
	}
}
