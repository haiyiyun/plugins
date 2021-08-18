package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApplicationMeta struct {
	Title              string `bson:"title" json:"title" map:"title"`
	Affix              bool   `bson:"affix" json:"affix" map:"affix"`
	Icon               string `bson:"icon" json:"icon,omitempty" map:"icon"`
	FrameSrc           string `bson:"frame_src" json:"frame_src,omitempty" map:"frame_src"`
	HideBreadcrumb     bool   `bson:"hide_breadcrumb" json:"hide_breadcrumb,omitempty" map:"hide_breadcrumb"`
	CarryParam         bool   `bson:"carry_param" json:"carry_param,omitempty" map:"carry_param"`
	HideChildrenInMenu bool   `bson:"hide_children_in_menu" json:"hide_children_in_menu,omitempty" map:"hide_children_in_menu"`
	HideTab            bool   `bson:"hide_tab" json:"hide_tab,omitempty" map:"hide_tab"`
	HideMenu           bool   `bson:"hide_menu" json:"hide_menu,omitempty" map:"hide_menu"`
}

type ApplicationModuleAction struct {
	Type   string          `bson:"type" json:"type" map:"type"`
	Name   string          `bson:"name" json:"name" map:"name"`
	Path   string          `bson:"path" json:"path" map:"path"`
	Order  int64           `bson:"order" json:"order" map:"order"`
	Meta   ApplicationMeta `bson:"meta" json:"meta" map:"meta"`
	Enable bool            `bson:"enable" json:"enable" map:"enable"`
	Method []string        `bson:"method" json:"method" map:"method"`
}

type ApplicationModule struct {
	Type    string                             `bson:"type" json:"type" map:"type"`
	Name    string                             `bson:"name" json:"name" map:"name"`
	Path    string                             `bson:"path" json:"path" map:"path"`
	Order   int64                              `bson:"order" json:"order" map:"order"`
	Meta    ApplicationMeta                    `bson:"meta" json:"meta" map:"meta"`
	Enable  bool                               `bson:"enable" json:"enable" map:"enable"`
	Actions map[string]ApplicationModuleAction `bson:"actions" json:"actions" map:"actions"`
}

///basic/application/index
type Application struct {
	ID                  primitive.ObjectID           `bson:"_id,omitempty" json:"_id" map:"_id"`
	Type                string                       `bson:"type" json:"type" map:"type"`
	Name                string                       `bson:"name" json:"name" map:"name"`
	Path                string                       `bson:"path" json:"path" map:"path"`
	Order               int64                        `bson:"order" json:"order" map:"order"`
	Meta                ApplicationMeta              `bson:"meta" json:"meta" map:"meta"`
	Enable              bool                         `bson:"enable" json:"enable" map:"enable"`
	DefaultModuleAction string                       `bson:"default_module_action" json:"default_module_action" map:"default_module_action"`
	Modules             map[string]ApplicationModule `bson:"modules" json:"modules" map:"modules"`
	CreateTime          time.Time                    `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime          time.Time                    `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
