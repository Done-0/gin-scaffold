// Package controller provides test controller
// Author: Done-0
// Created: 2025-09-25
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/internal/sse"
	"github.com/Done-0/gin-scaffold/internal/types/errno"
	"github.com/Done-0/gin-scaffold/internal/utils/errorx"
	"github.com/Done-0/gin-scaffold/internal/utils/validator"
	"github.com/Done-0/gin-scaffold/internal/utils/vo"
	"github.com/Done-0/gin-scaffold/pkg/serve/controller/dto"
	"github.com/Done-0/gin-scaffold/pkg/serve/service"
)

// TestController test HTTP controller
type TestController struct {
	testService service.TestService
	sseManager  sse.SSEManager
}

// NewTestController creates test controller
func NewTestController(testService service.TestService, sseManager sse.SSEManager) *TestController {
	return &TestController{
		testService: testService,
		sseManager:  sseManager,
	}
}

// TestPing handles ping test endpoint
// @Router /api/v1/test/testPing [get]
func (tc *TestController) TestPing(c *gin.Context) {
	response, err := tc.testService.TestPing(c)
	if err != nil {
		c.JSON(500, vo.Fail(c, err, errorx.New(errno.ErrInternalServer)))
		return
	}

	c.JSON(200, vo.Success(c, response))
}

// TestHello handles hello test endpoint
// @Router /api/v1/test/testHello [get]
func (tc *TestController) TestHello(c *gin.Context) {
	response, err := tc.testService.TestHello(c)
	if err != nil {
		c.JSON(500, vo.Fail(c, err, errorx.New(errno.ErrInternalServer)))
		return
	}

	c.JSON(200, vo.Success(c, response))
}

// TestLogger handles logger test endpoint
// @Router /api/v1/test/testLogger [get]
func (tc *TestController) TestLogger(c *gin.Context) {
	response, err := tc.testService.TestLogger(c)
	if err != nil {
		c.JSON(500, vo.Fail(c, err, errorx.New(errno.ErrInternalServer)))
		return
	}

	c.JSON(200, vo.Success(c, response))
}

// TestRedis handles redis test endpoint
// @Router /api/v1/test/testRedis [post]
func (tc *TestController) TestRedis(c *gin.Context) {
	req := &dto.TestRedisRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, vo.Fail(c, err, errorx.New(errno.ErrInvalidParams, errorx.KV("msg", "bind JSON failed"))))
		return
	}

	response, err := tc.testService.TestRedis(c, req)
	if err != nil {
		c.JSON(500, vo.Fail(c, err, errorx.New(errno.ErrInternalServer)))
		return
	}

	c.JSON(200, vo.Success(c, response))
}

// TestSuccess handles success test endpoint
// @Router /api/v1/test/testSuccessRes [get]
func (tc *TestController) TestSuccess(c *gin.Context) {
	response, err := tc.testService.TestSuccess(c)
	if err != nil {
		c.JSON(500, vo.Fail(c, err, errorx.New(errno.ErrInternalServer)))
		return
	}

	c.JSON(200, vo.Success(c, response))
}

// TestError handles error test endpoint
// @Router /api/v1/test/testErrRes [get]
func (tc *TestController) TestError(c *gin.Context) {
	response, err := tc.testService.TestError(c)
	if err != nil {
		c.JSON(500, vo.Fail(c, err, errorx.New(errno.ErrInternalServer)))
		return
	}

	c.JSON(200, vo.Success(c, response))
}

// TestErrorMiddleware handles error middleware test endpoint
// @Router /api/v1/test/testErrorMiddleware [get]
func (tc *TestController) TestErrorMiddleware(c *gin.Context) {
	// This will trigger the recovery middleware
	panic("Test panic for recovery middleware")
}

// TestLong handles long request test endpoint
// @Router /api/v2/test/testLongReq [post]
func (tc *TestController) TestLong(c *gin.Context) {
	req := &dto.TestLongRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, vo.Fail(c, err, errorx.New(errno.ErrInvalidParams, errorx.KV("msg", "bind JSON failed"))))
		return
	}

	errors := validator.Validate(req)
	if errors != nil {
		c.JSON(http.StatusBadRequest, vo.Fail(c, errors, errorx.New(errno.ErrInvalidParams, errorx.KV("msg", "validation failed"))))
		return
	}

	response, err := tc.testService.TestLong(c, req)
	if err != nil {
		c.JSON(500, vo.Fail(c, err, errorx.New(errno.ErrInternalServer)))
		return
	}

	c.JSON(200, vo.Success(c, response))
}

// TestI18n handles i18n test endpoint
// @Router /api/v1/test/testI18n [get]
func (tc *TestController) TestI18n(c *gin.Context) {
	response, err := tc.testService.TestI18n(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Fail(c, err, errorx.New(errno.ErrInternalServer, errorx.KV("msg", "i18n test failed"))))
		return
	}

	c.JSON(http.StatusOK, vo.Success(c, response))
}

// TestStream handles simple SSE streaming test endpoint
// @Router /api/v1/test/testStream [post]
func (tc *TestController) TestStream(c *gin.Context) {
	req := &dto.TestStreamRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, vo.Fail(c, err, errorx.New(errno.ErrInvalidParams, errorx.KV("msg", "bind JSON failed"))))
		return
	}

	errors := validator.Validate(req)
	if errors != nil {
		c.JSON(http.StatusBadRequest, vo.Fail(c, errors, errorx.New(errno.ErrInvalidParams, errorx.KV("msg", "validation failed"))))
		return
	}

	events, err := tc.testService.TestStream(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Fail(c, err, errorx.New(errno.ErrInternalServer)))
		return
	}

	_ = tc.sseManager.StreamToClient(c, events)
}
