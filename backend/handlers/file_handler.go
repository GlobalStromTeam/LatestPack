package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"

	"latestpack/models"
	"latestpack/services"
	"latestpack/utils"

	"github.com/gin-gonic/gin"
)

const maxUploadSize int64 = 500 << 20 // 500 MB

type FileHandler struct {
	svc *services.FileService
}

func NewFileHandler(svc *services.FileService) *FileHandler {
	return &FileHandler{svc: svc}
}

func (h *FileHandler) List(c *gin.Context) {
	path := c.Query("path")

	resp, err := h.svc.List(path)
	if err != nil {
		if errors.Is(err, services.ErrInvalidPath) {
			utils.Error(c, http.StatusBadRequest, "Invalid path")
			return
		}
		log.Printf("List files error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	utils.Success(c, http.StatusOK, resp)
}

func (h *FileHandler) CreateFolder(c *gin.Context) {
	var req models.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	item, err := h.svc.CreateFolder(req.Path, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrItemExists):
			utils.Error(c, http.StatusConflict, "Item already exists")
		case errors.Is(err, services.ErrInvalidPath):
			utils.Error(c, http.StatusBadRequest, "Invalid path")
		case errors.Is(err, services.ErrInvalidName):
			utils.Error(c, http.StatusBadRequest, "Invalid name")
		default:
			log.Printf("Create folder error: %v", err)
			utils.Error(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *FileHandler) Upload(c *gin.Context) {
	path := c.PostForm("path")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "File is required")
		return
	}
	defer file.Close()

	if header.Size > maxUploadSize {
		utils.Error(c, http.StatusRequestEntityTooLarge, "File too large")
		return
	}

	data, err := io.ReadAll(io.LimitReader(file, maxUploadSize+1))
	if err != nil {
		log.Printf("Read upload error: %v", err)
		utils.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	if int64(len(data)) > maxUploadSize {
		utils.Error(c, http.StatusRequestEntityTooLarge, "File too large")
		return
	}

	item, err := h.svc.SaveFile(path, header.Filename, data)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrItemExists):
			utils.Error(c, http.StatusConflict, "File already exists")
		case errors.Is(err, services.ErrInvalidPath):
			utils.Error(c, http.StatusBadRequest, "Invalid path")
		case errors.Is(err, services.ErrInvalidName):
			utils.Error(c, http.StatusBadRequest, "Invalid filename")
		case errors.Is(err, services.ErrFileTooLarge):
			utils.Error(c, http.StatusRequestEntityTooLarge, "File too large")
		default:
			log.Printf("Save file error: %v", err)
			utils.Error(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	utils.Success(c, http.StatusCreated, item)
}

func (h *FileHandler) Rename(c *gin.Context) {
	var req models.RenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	item, err := h.svc.Rename(req.Path, req.OldName, req.NewName)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrItemNotFound):
			utils.Error(c, http.StatusNotFound, "Item not found")
		case errors.Is(err, services.ErrItemExists):
			utils.Error(c, http.StatusConflict, "Item already exists")
		case errors.Is(err, services.ErrInvalidPath):
			utils.Error(c, http.StatusBadRequest, "Invalid path")
		case errors.Is(err, services.ErrInvalidName):
			utils.Error(c, http.StatusBadRequest, "Invalid name")
		default:
			log.Printf("Rename error: %v", err)
			utils.Error(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	utils.Success(c, http.StatusOK, item)
}

func (h *FileHandler) Delete(c *gin.Context) {
	var req models.DeleteItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	err := h.svc.Delete(req.Path, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrItemNotFound):
			utils.Error(c, http.StatusNotFound, "Item not found")
		case errors.Is(err, services.ErrInvalidPath):
			utils.Error(c, http.StatusBadRequest, "Invalid path")
		case errors.Is(err, services.ErrInvalidName):
			utils.Error(c, http.StatusBadRequest, "Invalid name")
		default:
			log.Printf("Delete file error: %v", err)
			utils.Error(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "Deleted"})
}
