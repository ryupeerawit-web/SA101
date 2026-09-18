package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ryu111/special-area-booking/internal/config"
	"github.com/ryu111/special-area-booking/internal/dto"
	"github.com/ryu111/special-area-booking/internal/middleware"
	"github.com/ryu111/special-area-booking/internal/models"
	"github.com/ryu111/special-area-booking/internal/utils"
)

type Auth struct {
	DB  *gorm.DB
	Cfg config.Config
}

// Register handles new user account creation.
func (a Auth) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, "invalid request")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var count int64
	a.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		Fail(c, http.StatusConflict, "email already registered")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		Fail(c, http.StatusInternalServerError, "could not hash password")
		return
	}

	user := models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		Role:         models.RoleMember,
		Status:       models.UserActive,
	}

	if err := a.DB.Create(&user).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "could not create user")
		return
	}

	token, _ := utils.Token(user.ID, user.Role, a.Cfg)
	OK(c, http.StatusCreated, gin.H{
		"token": token,
		"user":  user,
	})
}

// Login handles user authentication and JWT token issuance.
func (a Auth) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, "invalid request")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := a.DB.Where("email = ?", email).First(&user).Error; err != nil || !utils.CheckPassword(user.PasswordHash, req.Password) {
		Fail(c, http.StatusUnauthorized, "email or password is incorrect")
		return
	}

	if user.Status != models.UserActive {
		Fail(c, http.StatusForbidden, "user is suspended")
		return
	}

	token, _ := utils.Token(user.ID, user.Role, a.Cfg)
	OK(c, http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// Me retrieves the profile of the currently authenticated user.
func (a Auth) Me(c *gin.Context) {
	var user models.User
	if err := a.DB.First(&user, middleware.UserID(c)).Error; err != nil {
		Fail(c, http.StatusNotFound, "user not found")
		return
	}

	OK(c, http.StatusOK, user)
}
