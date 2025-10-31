package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/auth-service/internal/model"
	"github.com/phanthehoang2503/small-project/auth-service/internal/repo"
	"github.com/phanthehoang2503/small-project/pkg/auth"
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

func (h *AuthHandler) Register(c *gin.Context) {
	req := registerReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate mail or username
	if u, _ := h.Repo.GetUser(req.Email); u != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
		return
	}
	if u, _ := h.Repo.GetUser(req.Username); u != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already in use"})
		return
	}

	//hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashed),
	}

	if err := h.Repo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.Repo.GetUser(req.Login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(h.jwtSecret, u.ID, int(h.jwtExp.Hours()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"id":       u.ID,
		"email":    u.Email,
		"username": u.Username,
	})
}
