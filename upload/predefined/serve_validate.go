package predefined

type RequestServeUploadID struct {
	ID string `validate:"required,bson_object_id"`
}

type RequestServeFileType struct {
	FileType string `form:"file_type" validate:"required"`
}

type RequestServeFile struct {
	RequestServeFileType
	Remark string
}
