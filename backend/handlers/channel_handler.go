package handlers

import (
	"errors"
	"log"
	"net/http"

	"latestpack/models"
	"latestpack/services"
	"latestpack/utils"

	"github.com/gin-gonic/gin"
)

type ChannelHandler struct {
	svc *services.ChannelService
}

func NewChannelHandler(svc *services.ChannelService) *ChannelHandler {
	return &ChannelHandler{svc: svc}
}

func (h *ChannelHandler) List(c *gin.Context) {
	channels, err := h.svc.List(c.Request.Context())
	if err != nil {
		log.Printf("List channels error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	if channels == nil {
		channels = []models.ChannelDTO{}
	}
	utils.Success(c, http.StatusOK, channels)
}

func (h *ChannelHandler) Create(c *gin.Context) {
	var req models.CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	ch, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrChannelNameExists):
			utils.Error(c, http.StatusConflict, "Channel name already exists")
		case errors.Is(err, services.ErrCannotCreateLocal):
			utils.Error(c, http.StatusBadRequest, "Cannot create local channel")
		default:
			log.Printf("Create channel error: %v", err)
			utils.Error(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	utils.Success(c, http.StatusCreated, ch)
}

func (h *ChannelHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	ch, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrChannelNotFound):
			utils.Error(c, http.StatusNotFound, "Channel not found")
		case errors.Is(err, services.ErrChannelNameExists):
			utils.Error(c, http.StatusConflict, "Channel name already exists")
		case errors.Is(err, services.ErrLocalChannelLocked):
			utils.Error(c, http.StatusForbidden, "Local channel only allows enabled changes")
		default:
			log.Printf("Update channel error: %v", err)
			utils.Error(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	utils.Success(c, http.StatusOK, ch)
}

func (h *ChannelHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrChannelNotFound):
			utils.Error(c, http.StatusNotFound, "Channel not found")
		case errors.Is(err, services.ErrCannotDeleteLocal):
			utils.Error(c, http.StatusForbidden, "Cannot delete local channel")
		default:
			log.Printf("Delete channel error: %v", err)
			utils.Error(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "Deleted"})
}
