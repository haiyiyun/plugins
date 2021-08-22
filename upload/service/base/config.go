package base

import "github.com/haiyiyun/config"

type MongodbCfg struct {
	MongoDNS          string `json:"mongo_dns"`
	MongoDatabaseName string `json:"mongo_database_name"`
}

type CacheCfg struct {
	CacheDefaultExpiration config.Duration `json:"cache_default_expiration"`
	CacheCleanupInterval   config.Duration `json:"cache_cleanup_interval"`
}

type BaseCfg struct {
	AllowUploadLocal           bool   `json:"allow_upload_local"`
	UploadDirectory            string `json:"upload_directory"`
	UploadImageDirectory       string `json:"upload_image_directory"`
	UploadMediaDirectory       string `json:"upload_media_directory"`
	UploadDocumentDirectory    string `json:"upload_document_directory"`
	UploadCompressionDirectory string `json:"upload_compression_directory"`
	UploadFileDirectory        string `json:"upload_file_directory"`
}

type Config struct {
	MongodbCfg
	CacheCfg
	BaseCfg
}
