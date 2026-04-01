package handler

import (
	"net/http"
	"strconv"

	"zoomgate/internal/auth"
	"zoomgate/internal/config"
	"zoomgate/internal/model"
	"zoomgate/internal/store"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AdminAuthHandler struct {
	store  store.Store
	cfg    *config.Config
	logger *zap.Logger
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func NewAdminAuthHandler(s store.Store, cfg *config.Config, l *zap.Logger) *AdminAuthHandler {
	return &AdminAuthHandler{store: s, cfg: cfg, logger: l}
}

// Login handles user login and returns a JWT token.
func (h *AdminAuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	// JSON -> GO struct
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Get user by username
	userctx, err := h.store.GetUserByUsername(req.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Validate user password
	if err := bcrypt.CompareHashAndPassword([]byte(userctx.PasswordHash), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(h.cfg.Auth.JWTSecret, userctx.ID, userctx.Username, string(userctx.Role), h.cfg.Auth.TokenExpiry)
	if err != nil {
		h.logger.Error("failed to generate user token", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		Token:    token,
		Username: userctx.Username,
		Role:     string(userctx.Role),
	})
}

// UserHandler handles user CRUD operations.
type UserHandler struct {
	store  store.Store
	cfg    *config.Config
	logger *zap.Logger
}

type CreateUserRequest struct {
	Username  string     `json:"username" binding:"required"`
	Password  string     `json:"password" binding:"required,min=6"`
	Role      model.Role `json:"role" binding:"required"`
	Email     string     `json:"email"`
	RateLimit int        `json:"rate_limit"`
}

func NewUserHandler(s store.Store, cfg *config.Config, l *zap.Logger) *UserHandler {
	return &UserHandler{store: s, cfg: cfg, logger: l}
}

// List lists all users.
func (h *UserHandler) List(ctx *gin.Context) {
	// Extract page and page size from request query parameters, if not provided use default values
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	users, total, err := h.store.ListUsers(page, pageSize)
	if err != nil {
		h.logger.Error("failed to list users", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  page,
	})
}

// Create creates a new user.
func (h *UserHandler) Create(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	rateLimit := req.RateLimit
	if rateLimit == 0 {
		rateLimit = h.cfg.Security.RateLimit.DefaultRPM
	}

	userInfo := model.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Role:         req.Role,
		Email:        req.Email,
		RateLimit:    rateLimit,
	}

	err = h.store.CreateUser(&userInfo)
	if err != nil {
		h.logger.Error("failed to create user", zap.Error(err))
		ctx.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	ctx.JSON(http.StatusCreated, userInfo)
}

// Update updates an existing user.
func (h *UserHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.store.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req struct {
		Email     *string     `json:"email"`
		Role      *model.Role `json:"role"`
		RateLimit *int        `json:"rate_limit"`
		Password  *string     `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.RateLimit != nil {
		user.RateLimit = *req.RateLimit
	}
	if req.Password != nil && *req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		user.PasswordHash = string(hash)
	}

	if err := h.store.UpdateUser(user); err != nil {
		h.logger.Error("failed to update user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Delete deletes an existing user.
func (h *UserHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.store.DeleteUser(uint(id)); err != nil {
		h.logger.Error("failed to delete user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *UserHandler) GenerateAPIKey(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.store.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	fullKey, hash, prefix, err := auth.GenerateAPIKey()
	if err != nil {
		h.logger.Error("failed to generate API key", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	user.APIKeyHash = hash
	user.APIKeyPrefix = prefix
	if err := h.store.UpdateUser(user); err != nil {
		h.logger.Error("failed to save API key", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// Return the full key only once; it won't be retrievable again
	c.JSON(http.StatusOK, gin.H{
		"api_key": fullKey,
		"message": "Save this key securely. It will not be shown again.",
	})
}
