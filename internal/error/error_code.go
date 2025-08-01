// Package biz_err 提供错误码和错误信息定义
// 创建者：Done-0
// 创建时间：2025-07-01
package biz_err

// 错误码常量定义 - 基于 32 位无符号整型，采用 16 进制表示
// 格式：0x + 应用标识码(4位) + 错误类型码(4位)
// 应用标识码：0001 (metaphysics 应用)
const (
	SUCCESS = 0x00000000 // 成功

	// 参数错误 (0x00010500 - 0x000105FF)
	PARAM_ERROR   = 0x00010501 // 参数错误
	PARAM_MISSING = 0x00010502 // 必填参数为空
	PARAM_FORMAT  = 0x00010503 // 参数格式不正确
	PARAM_RANGE   = 0x00010504 // 参数范围不正确
	JSON_PARSE    = 0x00010505 // JSON格式解析失败

	// 认证授权错误 (0x00010100 - 0x000101FF)
	AUTH_ERROR    = 0x00010101 // 认证授权错误
	TOKEN_INVALID = 0x00010102 // 令牌无效
	TOKEN_EXPIRED = 0x00010103 // 令牌已过期
	UNAUTHORIZED  = 0x00010104 // 未授权访问
	FORBIDDEN     = 0x00010105 // 权限不足

	// 系统错误 (0x00010200 - 0x000102FF)
	SYSTEM_ERROR = 0x00010201 // 系统错误
	SERVICE_BUSY = 0x00010202 // 服务繁忙
	TIMEOUT      = 0x00010203 // 服务超时

	// 中间件错误 (0x00010300 - 0x000103FF)
	DB_ERROR    = 0x00010301 // 数据库错误
	CACHE_ERROR = 0x00010302 // 缓存错误
	MQ_ERROR    = 0x00010303 // 消息队列错误

	// HTTP错误 (0x00010400 - 0x000104FF)
	NOT_FOUND          = 0x00010404 // 资源不存在
	METHOD_NOT_ALLOWED = 0x00010405 // 方法不允许

	// 业务错误 (0x00010600 - 0x000106FF)
	BIZ_ERROR          = 0x00010601 // 业务错误
	RESOURCE_NOT_FOUND = 0x00010602 // 资源不存在
	CONFLICT           = 0x00010603 // 资源冲突
	RATE_LIMIT         = 0x00010607 // 请求频率超限

	// 未知错误
	UNKNOWN_ERROR = 0x00010801 // 未知错误
)

// CodeMsg 错误码对应的错误信息
var CodeMsg = map[int]string{
	SUCCESS: "操作成功",

	PARAM_ERROR:   "参数错误",
	PARAM_MISSING: "必填参数为空",
	PARAM_FORMAT:  "参数格式不正确",
	PARAM_RANGE:   "参数范围不正确",
	JSON_PARSE:    "JSON格式解析失败",

	AUTH_ERROR:    "认证失败",
	TOKEN_INVALID: "令牌无效",
	TOKEN_EXPIRED: "令牌已过期",
	UNAUTHORIZED:  "未授权访问",
	FORBIDDEN:     "权限不足",

	SYSTEM_ERROR: "系统内部错误",
	SERVICE_BUSY: "服务繁忙",
	TIMEOUT:      "服务超时",

	DB_ERROR:    "数据库错误",
	CACHE_ERROR: "缓存错误",
	MQ_ERROR:    "消息队列错误",

	NOT_FOUND:          "资源不存在",
	METHOD_NOT_ALLOWED: "方法不允许",

	BIZ_ERROR:          "业务处理失败",
	RESOURCE_NOT_FOUND: "资源不存在",
	CONFLICT:           "资源冲突",
	RATE_LIMIT:         "请求频率超限",

	UNKNOWN_ERROR: "未知错误",
}

// GetMessage 根据错误码获取对应的错误信息
// 参数：
//
//	code: 错误码
//
// 返回值：
//
//	string: 错误信息
func GetMessage(code int) string {
	if msg, ok := CodeMsg[code]; ok {
		return msg
	}
	return CodeMsg[UNKNOWN_ERROR]
}
