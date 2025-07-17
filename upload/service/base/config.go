package base

import "github.com/haiyiyun/config"

type MongodbConfig struct {
	MongoDNS          string `json:"mongo_dns"`
	MongoDatabaseName string `json:"mongo_database_name"`
}

type CacheConfig struct {
	CacheType              string          `json:"cache_type"`
	CacheUrl               string          `json:"cache_url"`
	CacheShardCount        string          `json:"cache_shard_count"`
	CacheUStrictTypeCheck  string          `json:"cache_strict_type_check"`
	CacheDefaultExpiration config.Duration `json:"cache_default_expiration"`
	CacheCleanupInterval   config.Duration `json:"cache_cleanup_interval"`
}

type BaseConfig struct {
	AppendFileExt              bool   `json:"append_file_ext"`
	AllowUploadLocal           bool   `json:"allow_upload_local"`
	UploadDirectory            string `json:"upload_directory"`
	UploadImageDirectory       string `json:"upload_image_directory"`
	UploadMediaDirectory       string `json:"upload_media_directory"`
	UploadDocumentDirectory    string `json:"upload_document_directory"`
	UploadCompressionDirectory string `json:"upload_compression_directory"`
	UploadFileDirectory        string `json:"upload_file_directory"`
}

type Config struct {
	MongodbConfig
	CacheConfig
	BaseConfig
}
