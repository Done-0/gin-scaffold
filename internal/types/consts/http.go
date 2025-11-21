// Package consts provides application constant definitions
// Author: Done-0
// Created: 2025-09-25
package consts

// HTTP method constants
const (
	MethodGET     = "GET"     // GET method
	MethodPOST    = "POST"    // POST method
	MethodPUT     = "PUT"     // PUT method
	MethodDELETE  = "DELETE"  // DELETE method
	MethodOPTIONS = "OPTIONS" // OPTIONS method
	MethodPATCH   = "PATCH"   // PATCH method
)

// HTTP header constants
const (
	// CORS related headers
	HeaderOrigin         = "Origin"           // Request origin
	HeaderContentType    = "Content-Type"     // Content type
	HeaderAccept         = "Accept"           // Accept type
	HeaderAcceptLanguage = "Accept-Language"  // Accept language
	HeaderAuthorization  = "Authorization"    // Authorization header
	HeaderXRequestedWith = "X-Requested-With" // AJAX request identifier
	HeaderContentLength  = "Content-Length"   // Content length

	// Network related headers
	HeaderRequestID     = "X-Request-ID"    // Request ID header
	HeaderXForwardedFor = "X-Forwarded-For" // Original client IP forwarded by proxy
	HeaderXRealIP       = "X-Real-IP"       // Real client IP
	HeaderXClientIP     = "X-Client-IP"     // Client IP (used by some proxies)
	HeaderUserAgent     = "User-Agent"      // User agent string
)
