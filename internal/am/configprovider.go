package am

import (
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ConfigProvider is an interface that replicates the methods provided by viper.Viper.
type ConfigProvider interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	SetDefault(key string, value interface{})
	SetEnvPrefix(in string)
	AutomaticEnv()
	SetEnvKeyReplacer(replacer *strings.Replacer)
	SetTypeByDefaultValue(enable bool)
	SetConfigFile(in string)
	ReadInConfig() error
	OnConfigChange(cb func(e fsnotify.Event))
	WatchConfig()
	Sub(key string) ConfigProvider
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetUint(key string) uint
	GetUint16(key string) uint16
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetSizeInBytes(key string) uint
}
