// Package configs 提供应用程序配置加载和更新功能
// 创建者：Done-0
// 创建时间：2025-07-01
package configs

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// AppConfig 应用配置
type AppConfig struct {
	AppName string `mapstructure:"APP_NAME"`
	AppHost string `mapstructure:"APP_HOST"`
	AppPort string `mapstructure:"APP_PORT"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DBDialect  string `mapstructure:"DB_DIALECT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PSW"`
	DBPath     string `mapstructure:"DB_PATH"`
}

// LogConfig 日志配置
type LogConfig struct {
	LogFilePath     string `mapstructure:"LOG_FILE_PATH"`
	LogFileName     string `mapstructure:"LOG_FILE_NAME"`
	LogTimestampFmt string `mapstructure:"LOG_TIMESTAMP_FMT"`
	LogMaxAge       int64  `mapstructure:"LOG_MAX_AGE"`
	LogRotationTime int64  `mapstructure:"LOG_ROTATION_TIME"`
	LogLevel        string `mapstructure:"LOG_LEVEL"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PSW"`
	RedisDB       string `mapstructure:"REDIS_DB"`
}

// Config 总配置结构
type Config struct {
	AppConfig   AppConfig      `mapstructure:"APP"`
	DBConfig    DatabaseConfig `mapstructure:"DATABASE"`
	LogConfig   LogConfig      `mapstructure:"LOG"`
	RedisConfig RedisConfig    `mapstructure:"REDIS"`
}

// DefaultConfigPath 默认配置文件路径
const DefaultConfigPath = "./configs/configs.yaml"

var (
	configInstance  *Config      // 全局配置实例
	configMutex     sync.RWMutex // 配置读写锁
	viperController *viper.Viper // viper实例
)

// New 初始化配置
// 参数：
//
//	configPath: 配置文件路径
//
// 返回值：
//
//	error: 初始化过程中的错误
func New(configPath string) error {
	viperController = viper.New()
	viperController.SetConfigFile(configPath)

	if err := viperController.ReadInConfig(); err != nil {
		return fmt.Errorf("配置文件读取失败: %w", err)
	}

	var config Config
	if err := viperController.Unmarshal(&config); err != nil {
		return fmt.Errorf("配置解析失败: %w", err)
	}

	configInstance = &config
	go monitorConfigChanges()
	return nil
}

// GetConfig 获取配置
// 返回值：
//
//	*Config: 配置副本
//	error: 获取过程中的错误
func GetConfig() (*Config, error) {
	configMutex.RLock()
	defer configMutex.RUnlock()

	if configInstance == nil {
		return nil, fmt.Errorf("配置未初始化")
	}

	configCopy := *configInstance
	return &configCopy, nil
}

// monitorConfigChanges 监听配置变更
func monitorConfigChanges() {
	viperController.WatchConfig()
	viperController.OnConfigChange(func(e fsnotify.Event) {
		var newConfig Config
		if err := viperController.Unmarshal(&newConfig); err != nil {
			log.Printf("新配置解析失败: %v", err)
			return
		}

		configMutex.Lock()
		defer configMutex.Unlock()

		oldConfig := *configInstance
		changes := make(map[string][2]any)

		if !compareStructs(oldConfig, newConfig, "", changes) {
			log.Printf("配置类型不一致，变更被阻止")
			return
		}

		configInstance = &newConfig

		for path, values := range changes {
			log.Printf("配置项 [%s] 发生变化: %v -> %v", path, values[0], values[1])
		}
	})
}

// compareStructs 比较结构体并收集变更
// 参数：
//
//	oldObj: 旧结构体
//	newObj: 新结构体
//	prefix: 字段路径前缀
//	changes: 记录变更的映射
//
// 返回值：
//
//	bool: 结构体类型是否一致
func compareStructs(oldObj, newObj any, prefix string, changes map[string][2]any) bool {
	oldVal := reflect.ValueOf(oldObj)
	newVal := reflect.ValueOf(newObj)

	if oldVal.Type() != newVal.Type() {
		return false
	}

	if oldVal.Kind() != reflect.Struct {
		return true
	}

	for i := 0; i < oldVal.NumField(); i++ {
		oldField := oldVal.Field(i)
		newField := newVal.Field(i)
		fieldName := oldVal.Type().Field(i).Name
		fullName := prefix + fieldName

		if oldField.Kind() == reflect.Struct {
			if !compareStructs(oldField.Interface(), newField.Interface(), fullName+".", changes) {
				return false
			}
			continue
		}

		if oldField.Kind() != newField.Kind() {
			return false
		}

		if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			changes[fullName] = [2]any{oldField.Interface(), newField.Interface()}
		}
	}

	return true
}
