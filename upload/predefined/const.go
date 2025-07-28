package predefined

const (
	UploadTypeImage       = "image"
	UploadTypeMedia       = "media"
	UploadTypeDocument    = "document"
	UploadTypeCompression = "compression"
	UploadTypeFile        = "file"
)

const (
	UploadStorageLocal   = "local"
	UploadStorageAliyun  = "aliyun"
	UploadStorageTencent = "tencent" // 添加腾讯云存储类型
	UploadStorageQiniu   = "qiniu"   // 添加七牛云存储类型
)

const (
	FormNameFile           = "file"
	FormNameFileBase64Data = "file_base64_data"
)

const (
	ErrorNotFoundUserIDFromRequestClaims = "not found user_id from request claims"
	ErrorNotFoundClaimsFromRequest       = "not found claims from request"
	ErrorFalidSaveFile                   = "Fail to save file"
	ErrorNotFoundEncodeData              = "Not Found encode data"
	ErrorNotFoundFormData                = "Not found form data"
	ErrorNotAllowUploadLocal             = "not allow upload local"
	ErrorNotAllowUploadAliyun            = "not allow upload aliyun"
	ErrorFileExtensionNotAllowed         = "file extension not allowed"
	ErrorNotAllowUploadTencent           = "not allow upload tencent" // 添加腾讯云错误信息
	ErrorNotAllowUploadQiniu             = "not allow upload qiniu"   // 添加七牛云错误信息
)
