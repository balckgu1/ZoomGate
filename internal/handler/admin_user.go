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

// AdminAuthHandler handles admin login and JWT token generation.
type AdminAuthHandler struct {
	dataStore store.Store
	appConfig *config.Config
	logger    *zap.Logger
}

// LoginRequest represents the JSON body for admin login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the JSON response after successful login.
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// NewAdminAuthHandler creates a new AdminAuthHandler with the given dependencies.
func NewAdminAuthHandler(dataStore store.Store, appConfig *config.Config, logger *zap.Logger) *AdminAuthHandler {
	return &AdminAuthHandler{dataStore: dataStore, appConfig: appConfig, logger: logger}
}

// Login validates user credentials and returns a JWT token on success.
func (handler *AdminAuthHandler) Login(ctx *gin.Context) {
	var loginReq LoginRequest
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Look up user by username
	user, err := handler.dataStore.GetUserByUsername(loginReq.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verify password against stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token with user claims
	token, err := auth.GenerateToken(
		handler.appConfig.Auth.JWTSecret,
		user.ID,
		user.Username,
		string(user.Role),
		handler.appConfig.Auth.TokenExpiry,
	)
	if err != nil {
		handler.logger.Error("failed to generate user token", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		Token:    token,
		Username: user.Username,
		Role:     string(user.Role),
	})
}

// UserHandler handles user CRUD operations for the admin API.
type UserHandler struct {
	dataStore store.Store
	appConfig *config.Config
	logger    *zap.Logger
}

// CreateUserRequest represents the JSON body for creating a new user.
type CreateUserRequest struct {
	Username  string     `json:"username" binding:"required"`
	Password  string     `json:"password" binding:"required,min=6"`
	Role      model.Role `json:"role" binding:"required"`
	Email     string     `json:"email"`
	RateLimit int        `json:"rate_limit"`
}

// NewUserHandler creates a new UserHandler with the given dependencies.
func NewUserHandler(dataStore store.Store, appConfig *config.Config, logger *zap.Logger) *UserHandler {
	return &UserHandler{dataStore: dataStore, appConfig: appConfig, logger: logger}
}

// List returns a paginated list of all users.
func (handler *UserHandler) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	users, total, err := handler.dataStore.ListUsers(page, pageSize)
	if err != nil {
		handler.logger.Error("failed to list users", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  page,
	})
}

// Create registers a new user with hashed password.
func (handler *UserHandler) Create(ctx *gin.Context) {
	var createReq CreateUserRequest
	if err := ctx.ShouldBindJSON(&createReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(createReq.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	rateLimit := createReq.RateLimit
	if rateLimit == 0 {
		rateLimit = handler.appConfig.Security.RateLimit.DefaultRPM
	}

	newUser := model.User{
		Username:     createReq.Username,
		PasswordHash: string(passwordHash),
		Role:         createReq.Role,
		Email:        createReq.Email,
		RateLimit:    rateLimit,
	}

	if err := handler.dataStore.CreateUser(&newUser); err != nil {
		handler.logger.Error("failed to create user", zap.Error(err))
		ctx.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	ctx.JSON(http.StatusCreated, newUser)
}

// UpdateUserRequest represents the JSON body for updating an existing user.
// All fields are optional pointers so that only provided fields are updated.
type UpdateUserRequest struct {
	Email     *string     `json:"email"`
	Role      *model.Role `json:"role"`
	RateLimit *int        `json:"rate_limit"`
	Password  *string     `json:"password"`
}

// Update modifies an existing user's fields. Only non-nil fields are applied.
func (handler *UserHandler) Update(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	existingUser, err := handler.dataStore.GetUserByID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var updateReq UpdateUserRequest
	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateReq.Email != nil {
		existingUser.Email = *updateReq.Email
	}
	if updateReq.Role != nil {
		existingUser.Role = *updateReq.Role
	}
	if updateReq.RateLimit != nil {
		existingUser.RateLimit = *updateReq.RateLimit
	}
	if updateReq.Password != nil && *updateReq.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*updateReq.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		existingUser.PasswordHash = string(passwordHash)
	}

	if err := handler.dataStore.UpdateUser(existingUser); err != nil {
		handler.logger.Error("failed to update user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, existingUser)
}

// Delete removes a user by their ID.
func (handler *UserHandler) Delete(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := handler.dataStore.DeleteUser(uint(userID)); err != nil {
		handler.logger.Error("failed to delete user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GenerateAPIKey creates a new API key for a user. The full key is returned
// only once in the response; subsequent retrievals are not possible.
func (handler *UserHandler) GenerateAPIKey(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	existingUser, err := handler.dataStore.GetUserByID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	fullKey, keyHash, keyPrefix, err := auth.GenerateAPIKey()
	if err != nil {
		handler.logger.Error("failed to generate API key", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	existingUser.APIKeyHash = keyHash
	existingUser.APIKeyPrefix = keyPrefix
	if err := handler.dataStore.UpdateUser(existingUser); err != nil {
		handler.logger.Error("failed to save API key", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// Return the full key only once; it cannot be retrieved again
	ctx.JSON(http.StatusOK, gin.H{
		"api_key": fullKey,
		"message": "Save this key securely. It will not be shown again.",
	})
}
