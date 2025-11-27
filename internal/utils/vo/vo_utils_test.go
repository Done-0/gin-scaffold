// Package vo provides common value objects test
// Author: Done-0
// Created: 2025-11-27
package vo

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Done-0/gin-scaffold/internal/types/errno"
	"github.com/Done-0/gin-scaffold/internal/utils/errorx"
)

func init() {
	gin.SetMode(gin.TestMode)
	errorx.Register(int32(errno.ErrInvalidParams), "invalid parameter: {{msg}}")
	errorx.Register(int32(errno.ErrInternalServer), "internal server error: {{msg}}")
}

func setupTestContext() *gin.Context {
	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)
	router.Use(requestid.New())
	c.Request = httptest.NewRequest("GET", "/test", nil)
	requestid.New()(c)
	return c
}

func TestSuccess(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		wantData any
	}{
		{"string data", "test data", "test data"},
		{"map data", map[string]string{"key": "value"}, map[string]string{"key": "value"}},
		{"error data", errors.New("error message"), "error message"},
		{"nil data", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestContext()
			result := Success(c, tt.data)

			assert.Equal(t, tt.wantData, result.Data)
			assert.Nil(t, result.Error)
			assert.NotEmpty(t, result.RequestId)
			assert.Greater(t, result.TimeStamp, int64(0))
		})
	}
}

func TestFail(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		err      error
		wantCode string
	}{
		{
			name:     "StatusError without params",
			data:     nil,
			err:      errorx.New(int32(errno.ErrInvalidParams)),
			wantCode: "10002",
		},
		{
			name:     "StatusError with params",
			data:     nil,
			err:      errorx.New(int32(errno.ErrInvalidParams), errorx.KV("msg", "username required")),
			wantCode: "10002",
		},
		{
			name:     "regular error",
			data:     nil,
			err:      errors.New("something went wrong"),
			wantCode: "10001",
		},
		{
			name:     "error data with regular error",
			data:     errors.New("data error"),
			err:      errors.New("system error"),
			wantCode: "10001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupTestContext()
			result := Fail(c, tt.data, tt.err)

			assert.NotNil(t, result.Error)
			assert.Equal(t, tt.wantCode, result.Error.Code)
			assert.NotEmpty(t, result.Error.Message)
			assert.NotEmpty(t, result.RequestId)
			assert.Greater(t, result.TimeStamp, int64(0))
		})
	}
}

func TestResultStructure(t *testing.T) {
	c := setupTestContext()

	t.Run("Success structure", func(t *testing.T) {
		result := Success(c, "test")

		assert.NotEmpty(t, result.RequestId)
		assert.Greater(t, result.TimeStamp, int64(0))
		assert.Nil(t, result.Error)
		assert.Equal(t, "test", result.Data)
	})

	t.Run("Fail structure", func(t *testing.T) {
		result := Fail(c, nil, errors.New("test error"))

		assert.NotEmpty(t, result.RequestId)
		assert.Greater(t, result.TimeStamp, int64(0))
		assert.NotNil(t, result.Error)
		assert.NotEmpty(t, result.Error.Code)
		assert.NotEmpty(t, result.Error.Message)
	})
}
