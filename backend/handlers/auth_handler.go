package handlers

import (
	"errors"
	"log"
	"net/http"

	"latestpack/services"
	"latestpack/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	token, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCreds) {
			utils.Error(c, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		log.Printf("Login error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	utils.Success(c, http.StatusOK, loginResponse{Token: token, Username: req.Username})
}

type updateUsernameRequest struct {
	Username string `json:"username" binding:"required"`
}

func (h *AuthHandler) UpdateUsername(c *gin.Context) {
	var req updateUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	currentUsername := c.GetString("username")
	newUsername, err := h.svc.UpdateUsername(c.Request.Context(), currentUsername, req.Username)
	if err != nil {
		if errors.Is(err, services.ErrUsernameExists) {
			utils.Error(c, http.StatusConflict, "Username already exists")
			return
		}
		log.Printf("Update username error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"username": newUsername})
}

type updatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	var req updatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	username := c.GetString("username")
	err := h.svc.UpdatePassword(c.Request.Context(), username, req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCurrentPassword):
			utils.Error(c, http.StatusBadRequest, "Current password is incorrect")
		case errors.Is(err, services.ErrPasswordTooShort):
			utils.Error(c, http.StatusBadRequest, "Password must be at least 6 characters")
		default:
			log.Printf("Update password error: %v", err)
			utils.Error(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "Password updated"})
}
