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
