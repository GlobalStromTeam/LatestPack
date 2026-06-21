package handlers

import (
	"log"
	"net/http"

	"latestpack/services"
	"latestpack/utils"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	versionSvc *services.VersionService
	statsSvc   *services.StatsService
}

func NewDashboardHandler(versionSvc *services.VersionService, statsSvc *services.StatsService) *DashboardHandler {
	return &DashboardHandler{versionSvc: versionSvc, statsSvc: statsSvc}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.statsSvc.GetStats(c.Request.Context())
	if err != nil {
		log.Printf("Get stats error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.Success(c, http.StatusOK, stats)
}

func (h *DashboardHandler) GetLatestVersion(c *gin.Context) {
	v, err := h.versionSvc.GetLatest(c.Request.Context())
	if err != nil {
		log.Printf("Get latest version error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	if v == nil {
		utils.Error(c, http.StatusNotFound, "No versions found")
		return
	}
	utils.Success(c, http.StatusOK, v)
}
