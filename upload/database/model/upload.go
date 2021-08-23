package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Upload struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"_id" map:"_id"`
	Type             string             `bson:"type" json:"type" map:"type"`
	Storage          string             `bson:"storage" json:"storage" map:"storage"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id" map:"user_id"`
	ContentType      string             `bson:"content_type" json:"content_type" map:"content_type"`
	OriginalFileName string             `bson:"original_file_name" json:"original_file_name" map:"original_file_name"`
	FileName         string             `bson:"file_name" json:"file_name" map:"file_name"`
	FileExt          string             `bson:"file_ext" json:"file_ext" map:"file_ext"`
	Path             string             `bson:"path" json:"path" map:"path"` //存储的时相对path。实际文件地址需要使用base.Config.UploadDirectory来拼接
	URL              string             `bson:"url" json:"url" map:"url"`    //存储的是相对url。使用内置文件服务时，实际地址需要使用Config.DownloadLocalUrlDirectory来手动拼接；使用非内置服务时，需自行配置拼接
	Size             int64              `bson:"size" json:"size" map:"size"`
	Remark           string             `bson:"remark" json:"remark" map:"remark"`
	CreateTime       time.Time          `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime       time.Time          `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
