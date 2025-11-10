package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/phanthehoang2503/small-project/auth-service/internal/model"
	"github.com/phanthehoang2503/small-project/auth-service/internal/repo"
	loggerclient "github.com/phanthehoang2503/small-project/internal/logger"
	"github.com/phanthehoang2503/small-project/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

type loginReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
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
		loggerclient.Warn(ctx, "register: bad payload", traceID, map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check duplicate email
	if u, err := h.Repo.GetUser(req.Email); err == nil && u != nil {
		loggerclient.Info(ctx, "register: email already in use", traceID, map[string]interface{}{
			"email": req.Email,
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
		return
	} else if err != nil && err != repo.ErrNotFound {
		loggerclient.Error(ctx, "register: repo error checking email", traceID, map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// Check duplicate username
	if u, err := h.Repo.GetUser(req.Username); err == nil && u != nil {
		loggerclient.Info(ctx, "register: username already in use", traceID, map[string]interface{}{
			"username": req.Username,
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already in use"})
		return
	} else if err != nil && err != repo.ErrNotFound {
		loggerclient.Error(ctx, "register: repo error checking username", traceID, map[string]interface{}{
			"error":    err.Error(),
			"username": req.Username,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	//hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		loggerclient.Error(ctx, "register: failed to hash password", traceID, map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashed),
	}

	if err := h.Repo.Create(user); err != nil {
		loggerclient.Error(ctx, "register: failed to create user", traceID, map[string]interface{}{
			"error":    err.Error(),
			"email":    req.Email,
			"username": req.Username,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	loggerclient.Info(ctx, "register: user created", traceID, map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})

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
		loggerclient.Warn(ctx, "login: bad payload", traceID, map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.Repo.GetUser(req.Login)
	if err != nil {
		loggerclient.Info(ctx, "login: user not found or repo error", traceID, map[string]interface{}{
			"login": req.Login,
			"error": err.Error(),
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		loggerclient.Info(ctx, "login: invalid password", traceID, map[string]interface{}{
			"user_id": u.ID,
			"login":   req.Login,
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := middleware.GenerateToken(h.jwtSecret, u.ID, int(h.jwtExp.Hours()))
	if err != nil {
		loggerclient.Error(ctx, "login: failed to create token", traceID, map[string]interface{}{
			"error":   err.Error(),
			"user_id": u.ID,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	loggerclient.Info(ctx, "login: success", traceID, map[string]interface{}{
		"user_id":  u.ID,
		"email":    u.Email,
		"username": u.Username,
	})

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"id":       u.ID,
		"email":    u.Email,
		"username": u.Username,
	})
}
