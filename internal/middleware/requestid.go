package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDHeader is the HTTP header name used to carry the request trace ID.
const RequestIDHeader = "X-Request-ID"

// RequestID is a middleware that ensures every request has a unique X-Request-ID.
// If the client provides one, it is reused; otherwise a new UUID is generated.
func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx.Set("request_id", requestID)
		ctx.Header(RequestIDHeader, requestID)
		ctx.Next()
	}
}
