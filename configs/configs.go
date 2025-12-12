// Package configs provides application configuration loading and updating functionality
// Author: Done-0
// Created: 2025-09-25
package configs

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// AppConfig application configuration
type AppConfig struct {
	AppName    string      `mapstructure:"APP_NAME"` // Application name
	AppHost    string      `mapstructure:"APP_HOST"` // Application host
	AppPort    string      `mapstructure:"APP_PORT"` // Application port
	CORSConfig CORSConfig  `mapstructure:"CORS"`     // CORS configuration
	Email      EmailConfig `mapstructure:"EMAIL"`    // Email configuration
	JWT        JWTConfig   `mapstructure:"JWT"`      // JWT authentication configuration
	User       UserConfig  `mapstructure:"USER"`     // User related configuration
}

// EmailConfig email configuration
type EmailConfig struct {
	EmailType string `mapstructure:"EMAIL_TYPE"` // Email type
	FromEmail string `mapstructure:"FROM_EMAIL"` // Sender email address
	EmailSmtp string `mapstructure:"EMAIL_SMTP"` // Email SMTP server
}

// JWTConfig JWT authentication configuration
type JWTConfig struct {
	Secret        string `mapstructure:"SECRET"`         // JWT signing secret
	ExpireTime    int64  `mapstructure:"EXPIRE_TIME"`    // Token expiration time (hours)
	RefreshExpire int64  `mapstructure:"REFRESH_EXPIRE"` // Refresh token expiration time (hours)
}

// UserConfig user related configuration
type UserConfig struct {
	SuperAdminEmail    string `mapstructure:"SUPER_ADMIN_EMAIL"`    // Administrator email
	SuperAdminPassword string `mapstructure:"SUPER_ADMIN_PASSWORD"` // Administrator password
	SuperAdminNickname string `mapstructure:"SUPER_ADMIN_NICKNAME"` // Administrator nickname
}

// DatabaseConfig database configuration
type DatabaseConfig struct {
	DBDialect  string `mapstructure:"DB_DIALECT"` // Database type
	DBName     string `mapstructure:"DB_NAME"`    // Database name
	DBHost     string `mapstructure:"DB_HOST"`    // Database host
	DBPort     string `mapstructure:"DB_PORT"`    // Database port
	DBUser     string `mapstructure:"DB_USER"`    // Database user
	DBPassword string `mapstructure:"DB_PSW"`     // Database password
	DBPath     string `mapstructure:"DB_PATH"`    // Database path
}

// LogConfig logging configuration
type LogConfig struct {
	LogFilePath     string `mapstructure:"LOG_FILE_PATH"`     // Log file path
	LogFileName     string `mapstructure:"LOG_FILE_NAME"`     // Log file name
	LogTimestampFmt string `mapstructure:"LOG_TIMESTAMP_FMT"` // Log timestamp format
	LogMaxAge       int64  `mapstructure:"LOG_MAX_AGE"`       // Log retention days
	LogRotationTime int64  `mapstructure:"LOG_ROTATION_TIME"` // Log rotation time (hours)
	LogLevel        string `mapstructure:"LOG_LEVEL"`         // Log level
}

// RedisConfig Redis configuration
type RedisConfig struct {
	RedisHost     string `mapstructure:"REDIS_HOST"` // Redis server address
	RedisPort     string `mapstructure:"REDIS_PORT"` // Redis server port
	RedisPassword string `mapstructure:"REDIS_PSW"`  // Redis password
	RedisDB       string `mapstructure:"REDIS_DB"`   // Redis database index

	// Connection pool settings
	PoolSize     int `mapstructure:"POOL_SIZE"`      // Maximum connection pool size
	MinIdleConns int `mapstructure:"MIN_IDLE_CONNS"` // Minimum idle connections
	DialTimeout  int `mapstructure:"DIAL_TIMEOUT"`   // Connection timeout in seconds
	ReadTimeout  int `mapstructure:"READ_TIMEOUT"`   // Read timeout in seconds
	WriteTimeout int `mapstructure:"WRITE_TIMEOUT"`  // Write timeout in seconds
}

// KafkaConfig Kafka configuration
type KafkaConfig struct {
	Brokers       []string `mapstructure:"BROKERS"`        // Kafka broker addresses
	ConsumerGroup string   `mapstructure:"CONSUMER_GROUP"` // Consumer group name
}

// CORSConfig CORS cross-origin configuration
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"ALLOW_ORIGINS"`     // Allowed origins
	AllowMethods     []string `mapstructure:"ALLOW_METHODS"`     // Allowed HTTP methods
	AllowHeaders     []string `mapstructure:"ALLOW_HEADERS"`     // Allowed headers
	ExposeHeaders    []string `mapstructure:"EXPOSE_HEADERS"`    // Exposed headers
	AllowCredentials bool     `mapstructure:"ALLOW_CREDENTIALS"` // Whether to allow credentials
	MaxAge           int64    `mapstructure:"MAX_AGE"`           // Preflight request cache time (hours)
}

// ProviderConfig generic provider configuration
type ProviderConfig struct {
	Enabled   bool                     `mapstructure:"ENABLED"`   // Whether provider is enabled
	Instances []ProviderInstanceConfig `mapstructure:"INSTANCES"` // Provider instances
}

// ProviderInstanceConfig individual provider instance configuration
type ProviderInstanceConfig struct {
	Name        string   `mapstructure:"NAME"`        // Instance name
	Enabled     bool     `mapstructure:"ENABLED"`     // Whether instance is enabled
	BaseURL     string   `mapstructure:"BASE_URL"`    // API base URL
	Keys        []string `mapstructure:"KEYS"`        // API key list
	Models      []string `mapstructure:"MODELS"`      // Available model list
	MaxTokens   int      `mapstructure:"MAX_TOKENS"`  // Maximum output tokens, controls response length
	Temperature float32  `mapstructure:"TEMPERATURE"` // Sampling temperature (0.0-2.0), higher=more creative, lower=more focused
	TopP        float32  `mapstructure:"TOP_P"`       // Nucleus sampling (0.0-1.0), controls vocabulary diversity, typically 0.9-0.95
	TopK        int      `mapstructure:"TOP_K"`       // Limits candidate words, 0=unlimited, typically 40-100
	Timeout     int      `mapstructure:"TIMEOUT"`     // Request timeout (seconds)
	MaxRetries  int      `mapstructure:"MAX_RETRIES"` // Maximum retry attempts
	RateLimit   string   `mapstructure:"RATE_LIMIT"`  // Rate limit (e.g., "60/min", "1/s")
}

// PromptConfig prompt template configuration
type PromptConfig struct {
	Dir string `mapstructure:"DIR"` // Prompt templates directory
}

// AIConfig AI service configuration
type AIConfig struct {
	Providers map[string]ProviderConfig `mapstructure:"PROVIDERS"` // Provider configurations
	Prompt    PromptConfig              `mapstructure:"PROMPT"`    // Prompt template configuration
}

// Config main configuration structure
type Config struct {
	AppConfig   AppConfig      `mapstructure:"APP"`      // Application configuration
	DBConfig    DatabaseConfig `mapstructure:"DATABASE"` // Database configuration
	LogConfig   LogConfig      `mapstructure:"LOG"`      // Logging configuration
	RedisConfig RedisConfig    `mapstructure:"REDIS"`    // Redis configuration
	KafkaConfig KafkaConfig    `mapstructure:"KAFKA"`    // Kafka configuration
	AI          AIConfig       `mapstructure:"AI"`       // AI service configuration
}

// Configuration file path constants
const (
	DefaultConfigPath = "./configs/config.local.yml" // Default configuration file path
	LocalConfigPath   = "./configs/config.local.yml" // Local development configuration
	ProdConfigPath    = "./configs/config.prod.yml"  // Production environment configuration
)

var (
	instance *Config      // Global configuration instance
	mu       sync.RWMutex // Configuration read-write lock
	v        *viper.Viper // Viper instance
)

// New initializes configuration
// - ENV=prod uses production configuration
func New() error {
	v = viper.New()

	env := os.Getenv("ENV")
	var configPath string
	switch env {
	case "prod", "production":
		configPath = ProdConfigPath
	case "local", "development":
		configPath = LocalConfigPath
	default:
		configPath = DefaultConfigPath
	}

	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	instance = &config
	go monitorConfigChanges()
	return nil
}

// GetConfig gets configuration
func GetConfig() (*Config, error) {
	mu.RLock()
	defer mu.RUnlock()

	if instance == nil {
		return nil, fmt.Errorf("config not initialized")
	}

	configCopy := *instance
	return &configCopy, nil
}

// monitorConfigChanges monitors configuration changes
func monitorConfigChanges() {
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		var newConfig Config
		if err := v.Unmarshal(&newConfig); err != nil {
			log.Printf("failed to unmarshal new config: %v", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		oldConfig := *instance
		changes := make(map[string][2]any)

		if !compareStructs(oldConfig, newConfig, "", changes) {
			log.Printf("config type mismatch, changes blocked")
			return
		}

		instance = &newConfig

		for path, values := range changes {
			log.Printf("config item [%s] changed: %v -> %v", path, values[0], values[1])
		}
	})
}

// compareStructs compares structs and collects changes
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

// UpdateField updates configuration field
func UpdateField(updateFunc func(*Config)) error {
	mu.Lock()
	defer mu.Unlock()

	oldConfig := *instance
	updateFunc(instance)

	configFile := v.ConfigFileUsed()
	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	newContent := string(content)

	var updateContent func(reflect.Value, reflect.Value, reflect.Type)
	updateContent = func(oldVal, newVal reflect.Value, t reflect.Type) {
		for i := 0; i < oldVal.NumField(); i++ {
			oldField, newField := oldVal.Field(i), newVal.Field(i)
			if tag := t.Field(i).Tag.Get("mapstructure"); tag != "" {
				if oldField.Kind() == reflect.Struct {
					updateContent(oldField, newField, oldField.Type())
				} else if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
					var old, new string
					if oldField.Kind() == reflect.Slice || oldField.Kind() == reflect.Array {
						// Array type
						var oldElems, newElems []string
						for i := 0; i < oldField.Len(); i++ {
							elem := oldField.Index(i)
							if elem.Kind() == reflect.String {
								oldElems = append(oldElems, fmt.Sprintf(`"%s"`, elem.String()))
							} else {
								oldElems = append(oldElems, fmt.Sprintf("%v", elem.Interface()))
							}
						}
						for i := 0; i < newField.Len(); i++ {
							elem := newField.Index(i)
							if elem.Kind() == reflect.String {
								newElems = append(newElems, fmt.Sprintf(`"%s"`, elem.String()))
							} else {
								newElems = append(newElems, fmt.Sprintf("%v", elem.Interface()))
							}
						}
						old, new = fmt.Sprintf("[%s]", strings.Join(oldElems, ", ")), fmt.Sprintf("[%s]", strings.Join(newElems, ", "))

						for _, pattern := range []string{fmt.Sprintf(`%s: %s`, tag, old), fmt.Sprintf(`%s: []`, tag)} {
							if strings.Contains(newContent, pattern) {
								newContent = strings.ReplaceAll(newContent, pattern, fmt.Sprintf(`%s: %s`, tag, new))
								break
							}
						}
					} else {
						// Non-array type
						old, new = fmt.Sprintf("%v", oldField.Interface()), fmt.Sprintf("%v", newField.Interface())
						var newFormatted string
						switch newField.Kind() {
						case reflect.String:
							newFormatted = fmt.Sprintf(`"%s"`, new)
						default:
							newFormatted = new
						}
						for _, pattern := range []string{
							fmt.Sprintf(`%s: "%s"`, tag, old),
							fmt.Sprintf(`%s: %s`, tag, old),
							fmt.Sprintf(`%s: ""`, tag),
						} {
							if strings.Contains(newContent, pattern) {
								newContent = strings.ReplaceAll(newContent, pattern, fmt.Sprintf(`%s: %s`, tag, newFormatted))
								break
							}
						}
					}
				}
			}
		}
	}

	updateContent(reflect.ValueOf(oldConfig), reflect.ValueOf(*instance), reflect.TypeOf(oldConfig))

	if newContent != string(content) {
		return os.WriteFile(configFile, []byte(newContent), 0644)
	}

	return nil
}
