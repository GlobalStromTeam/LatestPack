package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"latestpack/models"
	"latestpack/services"
	"latestpack/utils"

	"github.com/gin-gonic/gin"
)

const maxPageSize = 100

type VersionHandler struct {
	svc *services.VersionService
}

func NewVersionHandler(svc *services.VersionService) *VersionHandler {
	return &VersionHandler{svc: svc}
}

func (h *VersionHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 5
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	resp, err := h.svc.List(c.Request.Context(), page, pageSize)
	if err != nil {
		log.Printf("List versions error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.Success(c, http.StatusOK, resp)
}

func (h *VersionHandler) Create(c *gin.Context) {
	var req models.CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	v, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrVersionExists) {
			utils.Error(c, http.StatusBadRequest, "Version already exists")
			return
		}
		log.Printf("Create version error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.Success(c, http.StatusCreated, v)
}

func (h *VersionHandler) Delete(c *gin.Context) {
	version := c.Param("version")

	err := h.svc.Delete(c.Request.Context(), version)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrVersionNotFound):
			utils.Error(c, http.StatusNotFound, "Version not found")
		case errors.Is(err, services.ErrNotLatestVersion):
			utils.Error(c, http.StatusForbidden, "Only the latest version can be deleted")
		default:
			log.Printf("Delete version error: %v", err)
			utils.Error(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "Deleted"})
}
