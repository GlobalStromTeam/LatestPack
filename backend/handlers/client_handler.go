package handlers

import (
	"errors"
	"log"
	"net/http"

	"latestpack/services"
	"latestpack/utils"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	svc *services.ClientService
}

func NewClientHandler(svc *services.ClientService) *ClientHandler {
	return &ClientHandler{svc: svc}
}

func (h *ClientHandler) GetLatest(c *gin.Context) {
	resp, err := h.svc.GetLatest(c.Request.Context())
	if err != nil {
		log.Printf("Get latest error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	if resp == nil {
		utils.Error(c, http.StatusNotFound, "No version available")
		return
	}
	utils.Success(c, http.StatusOK, resp)
}

func (h *ClientHandler) GetUpdates(c *gin.Context) {
	from := c.Query("from")

	resp, err := h.svc.GetUpdates(c.Request.Context(), from)
	if err != nil {
		log.Printf("Get updates error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.Success(c, http.StatusOK, resp)
}

func (h *ClientHandler) Download(c *gin.Context) {
	version := c.Param("version")
	if version == "" {
		utils.Error(c, http.StatusBadRequest, "version parameter is required")
		return
	}

	path := c.Query("path")
	if path == "" {
		utils.Error(c, http.StatusBadRequest, "path parameter is required")
		return
	}

	ctx := c.Request.Context()

	if c.Request.Method == http.MethodHead {
		if err := h.svc.HeadFile(ctx, version, path, c.Writer); err != nil {
			switch {
			case errors.Is(err, services.ErrFileNotFound):
				utils.Error(c, http.StatusNotFound, "File not found")
			default:
				log.Printf("HEAD file error: %v", err)
				utils.Error(c, http.StatusInternalServerError, "Internal server error")
			}
		}
		return
	}

	if err := h.svc.DownloadFile(ctx, version, path, c.Writer, c.Request); err != nil {
		switch {
		case errors.Is(err, services.ErrFileNotFound):
			utils.Error(c, http.StatusNotFound, "File not found")
		default:
			log.Printf("Download file error: %v", err)
			utils.Error(c, http.StatusInternalServerError, "Internal server error")
		}
	}
}
