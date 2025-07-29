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
	StorageType              string   `json:"storage_type"`
	AppendFileExt            bool     `json:"append_file_ext"`
	AllowedFileExtensions    []string `json:"allowed_file_extensions"`
	DisallowedFileExtensions []string `json:"disallowed_file_extensions"`
}

type LocalConfig struct {
	AllowUploadLocal           bool   `json:"allow_upload_local"`
	UploadDirectory            string `json:"upload_directory"`
	UploadImageDirectory       string `json:"upload_image_directory"`
	UploadMediaDirectory       string `json:"upload_media_directory"`
	UploadDocumentDirectory    string `json:"upload_document_directory"`
	UploadCompressionDirectory string `json:"upload_compression_directory"`
	UploadFileDirectory        string `json:"upload_file_directory"`
}

type AliyunConfig struct {
	AliyunEndpoint          string `json:"aliyun_endpoint"`
	AliyunAccessKeyID       string `json:"aliyun_access_key_id"`
	AliyunAccessKeySecret   string `json:"aliyun_access_key_secret"`
	AliyunBucketName        string `json:"aliyun_bucket_name"`
	AliyunBaseURL           string `json:"aliyun_base_url"`
	AliyunUseInternal       bool   `json:"aliyun_use_internal"`        // 是否使用内网访问
	AliyunSecure            bool   `json:"aliyun_secure"`              // 是否使用HTTPS
	AliyunDisableUpload     bool   `json:"aliyun_disable_upload"`      // 是否禁用上传
	AliyunDisableDownload   bool   `json:"aliyun_disable_download"`    // 是否禁用下载
	AliyunDisableBucketCRUD bool   `json:"aliyun_disable_bucket_crud"` // 是否禁用Bucket的CRUD操作
	AliyunAppendFileExt     bool   `json:"aliyun_append_file_ext"`     // 是否添加文件扩展名
}

type TencentConfig struct {
	TencentEndpoint          string `json:"tencent_endpoint"`
	TencentAccessKeyID       string `json:"tencent_access_key_id"`
	TencentAccessKeySecret   string `json:"tencent_access_key_secret"`
	TencentBucketName        string `json:"tencent_bucket_name"`
	TencentBaseURL           string `json:"tencent_base_url"`
	TencentUseInternal       bool   `json:"tencent_use_internal"`        // 是否使用内网访问
	TencentSecure            bool   `json:"tencent_secure"`              // 是否使用HTTPS
	TencentDisableUpload     bool   `json:"tencent_disable_upload"`      // 是否禁用上传
	TencentDisableDownload   bool   `json:"tencent_disable_download"`    // 是否禁用下载
	TencentDisableBucketCRUD bool   `json:"tencent_disable_bucket_crud"` // 是否禁用Bucket的CRUD操作
	TencentAppendFileExt     bool   `json:"tencent_append_file_ext"`     // 是否添加文件扩展名
}

type QiniuConfig struct {
	QiniuAccessKey       string `json:"qiniu_access_key"`
	QiniuSecretKey       string `json:"qiniu_secret_key"`
	QiniuBucketName      string `json:"qiniu_bucket_name"`
	QiniuBaseURL         string `json:"qiniu_base_url"`
	QiniuUseHTTPS        bool   `json:"qiniu_use_https"`        // 是否使用HTTPS
	QiniuUseCdnDomains   bool   `json:"qiniu_use_cdn_domains"`  // 是否使用CDN加速上传
	QiniuDisableUpload   bool   `json:"qiniu_disable_upload"`   // 是否禁用上传
	QiniuDisableDownload bool   `json:"qiniu_disable_download"` // 是否禁用下载
}

type Config struct {
	MongodbConfig
	CacheConfig
	BaseConfig
	LocalConfig
	AliyunConfig
	TencentConfig
	QiniuConfig
}
