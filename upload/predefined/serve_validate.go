package predefined

import "go.mongodb.org/mongo-driver/bson/primitive"

type RequestServeUploadID struct {
	ID primitive.ObjectID `form:"id" validate:"required"`
}

type RequestServeFileType struct {
	FileType string `form:"file_type" validate:"required"`
}

type RequestServeFile struct {
	RequestServeFileType
	Remark string `form:"remark"`
}
