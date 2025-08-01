// Package cmd 提供应用程序的启动和运行入口
// 创建者：Done-0
// 创建时间：2025-07-01
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/configs"
	"github.com/Done-0/gin-scaffold/internal/db"
	"github.com/Done-0/gin-scaffold/internal/global"
	"github.com/Done-0/gin-scaffold/internal/logger"
	"github.com/Done-0/gin-scaffold/internal/middleware"
	"github.com/Done-0/gin-scaffold/internal/redis"
	"github.com/Done-0/gin-scaffold/pkg/router"
)

// Start 启动服务
func Start() {
	if err := configs.New(configs.DefaultConfigPath); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	config, err := configs.GetConfig()
	if err != nil {
		log.Fatalf("获取配置失败: %v", err)
	}

	// 初始化 Logger
	logger.New(config)

	// 初始化数据库连接并自动迁移模型
	db.New(config)

	// 初始化 Redis 连接
	redis.New(config)

	// 初始化 gin 实例
	app := gin.New()

	// 初始化中间件
	middleware.New(app)

	// 注册路由
	router.New(app)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.AppConfig.AppHost, config.AppConfig.AppPort),
		Handler: app,
	}

	// 上下文控制
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			global.SysLog.Errorf("服务启动失败: %v", err)
		}
	}()
	log.Printf("⇨ http server started on %s", srv.Addr)
	global.SysLog.Infof("⇨ http server started on %s", srv.Addr)

	// 等待中断信号
	<-ctx.Done()
	global.SysLog.Info("正在优雅关闭服务...")

	// 创建关闭超时上下文
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 优雅关闭
	if err := srv.Shutdown(shutdownCtx); err != nil {
		global.SysLog.Errorf("服务关闭异常: %v", err)
	}

	// 清理资源
	if err := db.Close(); err != nil {
		global.SysLog.Errorf("数据库关闭异常: %v", err)
	}

	if err := redis.Close(); err != nil {
		global.SysLog.Errorf("Redis 关闭异常: %v", err)
	}

	global.SysLog.Info("服务已优雅关闭")
}
