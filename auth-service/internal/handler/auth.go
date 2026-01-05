package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/phanthehoang2503/small-project/auth-service/internal/model"
	"github.com/phanthehoang2503/small-project/auth-service/internal/repo"
	logger "github.com/phanthehoang2503/small-project/internal/logger"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type registerReq struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Username string `json:"username" binding:"required,alphanum" example:"username123"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type loginReq struct {
	Login    string `json:"login" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type AuthHandler struct {
	Repo      repo.UserRepo
	jwtSecret []byte
	jwtExp    time.Duration
}

type registerResp struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type loginResp struct {
	Token    string `json:"token"`
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type errorResp struct {
	Error string `json:"error"`
}

func NewAuthHandler(r repo.UserRepo, secret []byte, expHours int) *AuthHandler {
	if expHours <= 0 {
		expHours = 72
	}

	return &AuthHandler{
		Repo:      r,
		jwtSecret: secret,
		jwtExp:    time.Duration(expHours) * time.Hour,
	}
}

func getTraceID(c *gin.Context) string {
	if tid := c.GetHeader("X-Request-ID"); tid != "" {
		return tid
	}

	return ""
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body registerReq true "Register payload"
// @Success 201 {object} registerResp
// @Failure 400 {object} errorResp
// @Failure 500 {object} errorResp
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	traceID := getTraceID(c)

	req := registerReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn(ctx, fmt.Sprintf("register: bad payload (trace_id=%s, err=%v)", traceID, err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check duplicate email
	if u, err := h.Repo.GetUser(req.Email); err == nil && u != nil {
		logger.Info(ctx, fmt.Sprintf("register: email already in use (trace_id=%s, email=%s)", traceID, req.Email))
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
		return
	} else if err != nil && err != repo.ErrNotFound {
		logger.Error(ctx, fmt.Sprintf("register: repo error checking email (trace_id=%s, email=%s, err=%v)", traceID, req.Email, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// Check duplicate username
	if u, err := h.Repo.GetUser(req.Username); err == nil && u != nil {
		logger.Info(ctx, fmt.Sprintf("register: username already in use (trace_id=%s, username=%s)", traceID, req.Username))
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already in use"})
		return
	} else if err != nil && err != repo.ErrNotFound {
		logger.Error(ctx, fmt.Sprintf("register: repo error checking username (trace_id=%s, username=%s, err=%v)", traceID, req.Username, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("register: failed to hash password (trace_id=%s, err=%v)", traceID, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashed),
	}

	if err := h.Repo.Create(user); err != nil {
		logger.Error(ctx, fmt.Sprintf("register: failed to create user (trace_id=%s, email=%s, username=%s, err=%v)", traceID, req.Email, req.Username, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	logger.Info(ctx, fmt.Sprintf("register: user created (trace_id=%s, id=%d, email=%s, username=%s)", traceID, user.ID, user.Email, user.Username))

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
}

// Login godoc
// @Summary Login with email or username
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body loginReq true "Login payload"
// @Success 200 {object} loginResp
// @Failure 400 {object} errorResp
// @Failure 401 {object} errorResp
// @Failure 500 {object} errorResp
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	traceID := getTraceID(c)

	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn(ctx, fmt.Sprintf("login: bad payload (trace_id=%s, err=%v)", traceID, err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.Repo.GetUser(req.Login)
	if err != nil {
		logger.Info(ctx, fmt.Sprintf("login: user not found or repo error (trace_id=%s, login=%s, err=%v)", traceID, req.Login, err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		logger.Info(ctx, fmt.Sprintf("login: invalid password (trace_id=%s, user_id=%d, login=%s)", traceID, u.ID, req.Login))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := middleware.GenerateToken(h.jwtSecret, u.ID, int(h.jwtExp.Hours()))
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("login: failed to create token (trace_id=%s, user_id=%d, err=%v)", traceID, u.ID, err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	logger.Info(ctx, fmt.Sprintf("login: success (trace_id=%s, user_id=%d, email=%s, username=%s)", traceID, u.ID, u.Email, u.Username))

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"id":       u.ID,
		"email":    u.Email,
		"username": u.Username,
	})
}
