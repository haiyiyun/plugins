package base

import (
	"github.com/haiyiyun/config"
)

type MongodbCfg struct {
	MongoDNS          string `json:"mongo_dns"`
	MongoDatabaseName string `json:"mongo_database_name"`
}

type CacheCfg struct {
	CacheType              string          `json:"cache_type"`
	CacheUrl               string          `json:"cache_url"`
	CacheShardCount        string          `json:"cache_shard_count"`
	CacheUStrictTypeCheck  string          `json:"cache_strict_type_check"`
	CacheDefaultExpiration config.Duration `json:"cache_default_expiration"`
	CacheCleanupInterval   config.Duration `json:"cache_cleanup_interval"`
}

type BaseCfg struct {
	DefaultRole                      string                     `json:"default_role"`
	DefaultTag                       string                     `json:"default_tag"`
	DefaultLevel                     int                        `json:"default_level"`
	CheckLogin                       bool                       `json:"check_login"`
	TokenByUrlQuery                  bool                       `json:"token_by_url_query"`
	TokenByUrlQueryName              string                     `json:"token_by_url_query_name"`
	IgnoreCheckLoginPath             map[string][]string        `json:"ignore_check_login_path"`
	TokenExpireDuration              config.Duration            `json:"token_expire_duration"`
	SpecifyUserIDTokenExpireDuration map[string]config.Duration `json:"specify_user_id_token_expire_duration"` //特别指定user_id的token过期时间
	OnlySingleLogin                  bool                       `json:"only_single_login"`                     //设置后，allow_multi_login，allow_multi_login_num不起作用
	OnlySingleLoginUserID            []string                   `json:"only_single_login_user_id"`             //不受only_single_login控制，设置后的user_id将使allow_multi_login，allow_multi_login_num不起作用
	OnlySingleLoginUserIDUnlimited   []string                   `json:"only_single_login_user_id_unlimited"`   //不受only_single_login控制，设置后的user_id受allow_multi_login，allow_multi_login_num影响
	AllowMultiLogin                  bool                       `json:"allow_multi_login"`
	AllowMultiLoginNum               int64                      `json:"allow_multi_login_num"`
	AllowMultiLoginUserIDUnlimited   []string                   `json:"allow_multi_login_user_id_unlimited"` //在允许allow_multi_login的情况下，设置后的user_id不受allow_multi_login_num限制
}

type Config struct {
	MongodbCfg
	CacheCfg
	BaseCfg
}
