package models

import "time"

type Version struct {
	Version   string
	Date      string
	Size      string
	Notes     string // JSON-encoded []string
	CreatedAt time.Time
}

type CreateVersionRequest struct {
	Version string   `json:"version" binding:"required"`
	Notes   []string `json:"notes" binding:"required"`
}

type VersionListResponse struct {
	Items    []VersionDTO `json:"items"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
}

type VersionDTO struct {
	Version string   `json:"version"`
	Date    string   `json:"date"`
	Size    string   `json:"size"`
	Notes   []string `json:"notes"`
}
