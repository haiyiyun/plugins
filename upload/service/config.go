package service

import "github.com/haiyiyun/config"

type MongodbCfg struct {
	MongoDNS          string `json:"mongo_dns"`
	MongoDatabaseName string `json:"mongo_database_name"`
}

type CacheCfg struct {
	CacheDefaultExpiration config.Duration `json:"cache_default_expiration"`
	CacheCleanupInterval   config.Duration `json:"cache_cleanup_interval"`
}

type UploadFileCfg struct {
	WebRouter                  bool   `json:"web_router"`
	WebRouterRootPath          string `json:"web_router_root_path"`
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
	UploadFileCfg
}
