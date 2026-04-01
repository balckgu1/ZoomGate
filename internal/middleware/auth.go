package middleware

import (
	"net/http"
	"strings"

	"zoomgate/internal/auth"
	"zoomgate/internal/store"

	"github.com/gin-gonic/gin"
)

// Context keys used to store authenticated user info in the Gin context.
const (
	CtxUserID   = "user_id"
	CtxUsername = "username"
	CtxRole     = "role"
)

// JWTAuth returns a middleware that validates JWT tokens from the Authorization header.
// It is intended for admin UI routes. On success it injects user_id, username, and role
// into the Gin context.
func JWTAuth(jwtSecret string, _ store.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := extractBearerToken(ctx)
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
			return
		}

		claims, err := auth.ParseToken(jwtSecret, tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		ctx.Set(CtxUserID, claims.UserID)
		ctx.Set(CtxUsername, claims.Username)
		ctx.Set(CtxRole, claims.Role)
		ctx.Next()
	}
}

// APIKeyAuth returns a middleware that validates API keys from the Authorization header.
// It is intended for LLM proxy routes (e.g. /v1/chat/completions). The API key is hashed
// with SHA-256 and looked up in the database.
func APIKeyAuth(dataStore store.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := extractBearerToken(ctx)
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
			return
		}

		apiKeyHash := auth.HashAPIKey(apiKey)
		user, err := dataStore.GetUserByAPIKeyHash(apiKeyHash)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			return
		}

		ctx.Set(CtxUserID, user.ID)
		ctx.Set(CtxUsername, user.Username)
		ctx.Set(CtxRole, string(user.Role))
		ctx.Next()
	}
}

// RequireRole returns a middleware that checks the authenticated user has the required role.
// Admin users are granted access to all roles.
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get(CtxRole)
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		if !auth.HasRole(userRole.(string), requiredRole) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		ctx.Next()
	}
}

// extractBearerToken extracts the token value from an "Authorization: Bearer <token>" header.
// It also accepts a raw token without the "Bearer" prefix for convenience.
func extractBearerToken(ctx *gin.Context) string {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(header)
}
