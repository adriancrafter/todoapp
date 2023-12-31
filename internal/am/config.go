package am

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	configFileKey = "config.file"
)

var (
	fileFlag   string
	searchFlag bool
)

var (
	defaultFilePath    = "config.yml"
	defaultFilePathDev = filepath.FromSlash("configs/config.yml")
)

type (
	Config struct {
		name   string
		v      *viper.Viper
		file   string
		search bool
		status map[time.Time]string
	}
)

func init() {
	flag.StringVar(&fileFlag, "config-file", defaultFilePath, "path to configuration file")
	flag.BoolVar(&searchFlag, "config-search", false, "search for settings in common config places")
}

func NewConfig(name string) *Config {
	return &Config{
		name:   name,
		v:      viper.New(),
		search: false,
		status: map[time.Time]string{},
	}
}

func (cfg *Config) Load(file ...string) (updated *Config, err error) {
	if len(file) == 0 {
		flag.Parse()
		cfg.file = fileFlag
		cfg.search = searchFlag
	} else {
		cfg.file = file[0]
		cfg.search = false
	}

	cfg.cleanInput()
	cfg.updateLookupPaths()

	return cfg.load()
}

func (cfg *Config) File() string {
	return cfg.file
}

func (cfg *Config) SetFile(file string) {
	cfg.file = file
}

func (cfg *Config) load() (updated *Config, err error) {
	cfg.v.SetDefault(configFileKey, cfg.file)
	cfg.v.AutomaticEnv()
	cfg.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.v.SetTypeByDefaultValue(true)
	cfg.v.SetConfigFile(cfg.file)

	err = cfg.v.ReadInConfig()
	if _, ok := err.(*os.PathError); ok {
		cfg.addStatus(fmt.Sprintf("no config file at '%s', using default values", cfg.file))

	} else if err != nil {
		return cfg, fmt.Errorf("error reading config: %w", err)
	}

	cfg.v.OnConfigChange(cfg.defaultOnConfigChangeFunc())
	cfg.v.WatchConfig()

	return cfg, nil
}

func (cfg *Config) cleanInput() {
	cfg.file = strings.Trim(cfg.file, " ")
}

func (cfg *Config) updateLookupPaths() {
	file, _, noParent := cfg.baseAndPath()

	if file == "" || file == "." {
		file = defaultFilePathDev
	}

	cfg.v.SetConfigType("yaml")
	cfg.v.AddConfigPath(file)

	if cfg.search && noParent {
		cfg.v.AddConfigPath("/etc/" + cfg.name)
		cfg.v.AddConfigPath("$HOME/." + cfg.name)
		cfg.v.AddConfigPath(".")
		return
	}
}

func (cfg *Config) baseAndPath() (base, dir string, current bool) {
	base = filepath.Base(cfg.file)
	dir = filepath.Dir(cfg.file)
	return base, dir, dir == "."

}

func (cfg *Config) SetOnConfigChange(onConfigChangeFunc func(e fsnotify.Event)) {
	cfg.v.OnConfigChange(onConfigChangeFunc)
}

func (cfg *Config) Sub(branch string) *Config {
	return &Config{
		name:   cfg.name,
		v:      viper.Sub(branch),
		search: cfg.search,
	}
}

func (cfg *Config) Get(prop string) *Config {
	return &Config{
		name:   cfg.name,
		v:      viper.Sub(prop),
		search: cfg.search,
	}
}

// GetString returns the value associated with the key as a string.
func (cfg *Config) GetString(key string) string {
	return cfg.v.GetString(key)
}

// GetBool returns the value associated with the key as a boolean.
func (cfg *Config) GetBool(key string) bool {
	return cfg.v.GetBool(key)
}

// GetInt returns the value associated with the key as an integer.
func (cfg *Config) GetInt(key string) int {
	return cfg.v.GetInt(key)
}

// GetInt32 returns the value associated with the key as an integer.
func (cfg *Config) GetInt32(key string) int32 {
	return cfg.v.GetInt32(key)
}

// GetInt64 returns the value associated with the key as an integer.
func (cfg *Config) GetInt64(key string) int64 {
	return cfg.v.GetInt64(key)
}

// GetUint returns the value associated with the key as an unsigned integer.
func (cfg *Config) GetUint(key string) uint {
	return cfg.v.GetUint(key)
}

// GetUint16 returns the value associated with the key as an unsigned integer.
func (cfg *Config) GetUint16(key string) uint16 {
	return cfg.v.GetUint16(key)
}

// GetUint32 returns the value associated with the key as an unsigned integer.
func (cfg *Config) GetUint32(key string) uint32 {
	return cfg.v.GetUint32(key)
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (cfg *Config) GetUint64(key string) uint64 {
	return cfg.v.GetUint64(key)
}

// GetFloat64 returns the value associated with the key as a float64.
func (cfg *Config) GetFloat64(key string) float64 {
	return cfg.v.GetFloat64(key)
}

// GetTime returns the value associated with the key as time.
func (cfg *Config) GetTime(key string) time.Time {
	return cfg.v.GetTime(key)
}

// GetDuration returns the value associated with the key as a duration.
func (cfg *Config) GetDuration(key string) time.Duration {
	return cfg.v.GetDuration(key)
}

// GetIntSlice returns the value associated with the key as a slice of int values.
func (cfg *Config) GetIntSlice(key string) []int {
	return cfg.v.GetIntSlice(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (cfg *Config) GetStringSlice(key string) []string {
	return cfg.v.GetStringSlice(key)
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (cfg *Config) GetStringMap(key string) map[string]interface{} {
	return cfg.v.GetStringMap(key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (cfg *Config) GetStringMapString(key string) map[string]string {
	return cfg.v.GetStringMapString(key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (cfg *Config) GetStringMapStringSlice(key string) map[string][]string {
	return cfg.v.GetStringMapStringSlice(key)
}

// GetSizeInBytes returns the size of the value associated with the given key
func (cfg *Config) GetSizeInBytes(key string) uint {
	return cfg.v.GetSizeInBytes(key)

}

func (cfg *Config) List() string {
	var sb strings.Builder
	for k, v := range cfg.v.AllSettings() {
		entry := fmt.Sprintf("%s: %v\n", k, v)
		sb.WriteString(entry)
	}
	return sb.String()
}

func (cfg *Config) addStatus(message string) {
	cfg.status[time.Now()] = message
}

func (cfg *Config) defaultOnConfigChangeFunc() func(e fsnotify.Event) {
	return func(e fsnotify.Event) {
		cfg.addStatus(fmt.Sprintf("Config file updated: %s", e.Name))
	}
}
