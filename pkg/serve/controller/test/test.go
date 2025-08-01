package test

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	bizErr "github.com/Done-0/gin-scaffold/internal/error"
	"github.com/Done-0/gin-scaffold/internal/global"
	"github.com/Done-0/gin-scaffold/internal/utils"
	"github.com/Done-0/gin-scaffold/pkg/vo"
)

// TestPing godoc
// @Summary Ping接口
// @Description 测试接口
// @Tags test
// @Accept json
// @Produce plain
// @Success 200 {string} string "Pong successfully!"
// @Router /test/testPing [get]
func TestPing(c *gin.Context) {
	utils.BizLogger(c).Info("Ping...")
	c.String(http.StatusOK, "Pong successfully!")
}

// TestHello godoc
// @Summary Hello接口
// @Description 测试接口
// @Tags test
// @Accept json
// @Produce plain
// @Success 200 {string} string "Hello, gin-scaffold 🎉!"
// @Router /test/testHello [get]
func TestHello(c *gin.Context) {
	utils.BizLogger(c).Info("Hello, gin-scaffold!")
	c.String(http.StatusOK, "Hello, gin-scaffold! 🎉!\n")
}

// TestLogger godoc
// @Summary 测试日志接口
// @Description 测试日志功能
// @Tags test
// @Accept json
// @Produce plain
// @Success 200 {string} string "测试日志成功"
// @Router /test/testLogger [get]
func TestLogger(c *gin.Context) {
	utils.BizLogger(c).Info("测试日志...")
	c.String(http.StatusOK, "测试日志成功")
}

// TestRedis godoc
// @Summary 测试Redis接口
// @Description 测试Redis功能
// @Tags test
// @Accept json
// @Produce plain
// @Success 200 {string} string "测试缓存功能完成"
// @Router /test/testRedis [get]
func TestRedis(c *gin.Context) {
	utils.BizLogger(c).Info("开始写入缓存...")
	err := global.RedisClient.Set(c.Request.Context(), "TEST:", "测试 value", 0).Err()
	if err != nil {
		utils.BizLogger(c).Error("测试写入缓存失败:", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	utils.BizLogger(c).Info("开始读取缓存...")
	val, err := global.RedisClient.Get(c.Request.Context(), "TEST:").Result()
	if err != nil {
		utils.BizLogger(c).Error("测试读取缓存失败:", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	utils.BizLogger(c).Infof("读取缓存成功, key: %s, value: %s", "TEST:", val)
	c.String(http.StatusOK, "测试缓存功能完成")
}

// TestSuccRes godoc
// @Summary 测试成功响应接口
// @Description 测试成功响应
// @Tags test
// @Accept json
// @Produce json
// @Success 200 {object} vo.Result "测试成功响应"
// @Router /test/testSuccessRes [get]
func TestSuccRes(c *gin.Context) {
	utils.BizLogger(c).Info("测试成功响应...")
	c.JSON(http.StatusOK, vo.Success(c, "测试成功响应成功"))
}

// TestErrRes godoc
// @Summary 测试错误响应接口
// @Description 测试错误响应
// @Tags test
// @Accept json
// @Produce json
// @Success 500 {object} vo.Result "测试错误响应"
// @Router /test/testErrRes [get]
func TestErrRes(c *gin.Context) {
	utils.BizLogger(c).Info("测试失败响应...")
	c.JSON(http.StatusInternalServerError, vo.Fail(c, nil, bizErr.New(bizErr.SYSTEM_ERROR)))
}

// TestErrorMiddleware godoc
// @Summary 测试错误处理中间件接口
// @Description 测试错误中间件
// @Tags test
// @Accept json
// @Produce json
// @Success 500 {object} vo.Result "测试错误中间件"
// @Router /test/testErrorMiddleware [get]
func TestErrorMiddleware(c *gin.Context) {
	utils.BizLogger(c).Info("测试错误处理中间件...")
	panic("测试错误处理中间件...")
}

// TestLongReq godoc
// @Summary 长时间请求接口
// @Description 模拟耗时请求
// @Tags test
// @Accept json
// @Produce plain
// @Success 200 {string} string "模拟耗时请求处理完成"
// @Router /test/testLongReq [get]
func TestLongReq(c *gin.Context) {
	utils.BizLogger(c).Info("开始测试耗时请求...")
	time.Sleep(20 * time.Second)
	c.String(http.StatusOK, "模拟耗时请求处理完成")
}
