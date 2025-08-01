// Package db 提供数据库连接和管理功能
// 创建者：Done-0
// 创建时间：2025-07-01
package db

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/Done-0/gin-scaffold/configs"
	"github.com/Done-0/gin-scaffold/internal/global"
)

// 数据库类型常量
const (
	DIALECT_POSTGRES = "postgres" // PostgreSQL 数据库
	DIALECT_SQLITE   = "sqlite"   // SQLite 数据库
	DIALECT_MYSQL    = "mysql"    // MySQL 数据库
)

// New 初始化数据库连接
// 参数：
//
//	config: 配置信息
func New(config *configs.Config) {
	dialect := config.DBConfig.DBDialect

	// 处理不同数据库类型的初始化
	switch dialect {
	case DIALECT_SQLITE:
		if err := ensureDBExists(nil, config); err != nil {
			global.SysLog.Fatalf("数据库不存在且创建失败: %v", err)
		}
	case DIALECT_POSTGRES, DIALECT_MYSQL:
		systemDBName := getSystemDBName(dialect)
		systemDB, err := connectToDB(config, systemDBName)
		if err != nil {
			global.SysLog.Fatalf("连接系统数据库失败: %v", err)
		}

		if err := ensureDBExists(systemDB, config); err != nil {
			global.SysLog.Fatalf("数据库不存在且创建失败: %v", err)
		}

		if sqlDB, err := systemDB.DB(); err == nil {
			sqlDB.Close()
		}
	default:
		global.SysLog.Fatalf("不支持的数据库类型: %s", dialect)
	}

	// 连接目标数据库
	var err error
	global.DB, err = connectToDB(config, config.DBConfig.DBName)
	if err != nil {
		global.SysLog.Fatalf("连接数据库失败: %v", err)
	}

	log.Printf("「%s」数据库连接成功...", config.DBConfig.DBName)
	global.SysLog.Infof("「%s」数据库连接成功...", config.DBConfig.DBName)

	// 执行数据库迁移
	if err = autoMigrate(); err != nil {
		global.SysLog.Fatalf("数据库自动迁移失败: %v", err)
	}
}

// Close 关闭数据库连接
// 返回值：
//
//	error: 错误信息
func Close() error {
	if global.DB == nil {
		return nil
	}

	sqlDB, err := global.DB.DB()
	if err != nil {
		return fmt.Errorf("获取 SQL 数据库实例失败: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接失败: %w", err)
	}

	global.SysLog.Info("数据库连接已关闭")
	return nil
}

// getSystemDBName 获取系统数据库名称
// 参数：
//
//	dialect: 数据库类型
//
// 返回值：
//
//	string: 系统数据库名称
func getSystemDBName(dialect string) string {
	switch dialect {
	case DIALECT_POSTGRES:
		return "postgres"
	case DIALECT_MYSQL:
		return "information_schema"
	default:
		return ""
	}
}

// connectToDB 连接到指定数据库
// 参数：
//
//	config: 配置信息
//	dbName: 数据库名称
//
// 返回值：
//
//	*gorm.DB: 数据库连接实例
//	error: 错误信息
func connectToDB(config *configs.Config, dbName string) (*gorm.DB, error) {
	dialector, err := getDialector(config, dbName)
	if err != nil {
		return nil, fmt.Errorf("获取数据库驱动器失败: %v", err)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	return db, nil
}
