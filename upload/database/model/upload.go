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
	Path             string             `bson:"path" json:"path" map:"path"`
	URL              string             `bson:"url" json:"url" map:"url"`
	Size             int64              `bson:"size" json:"size" map:"size"`
	CreateTime       time.Time          `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime       time.Time          `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
