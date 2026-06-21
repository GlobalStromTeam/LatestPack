package models

import "time"

type Version struct {
	Version   string
	Date      string
	Size      string
	Notes     string // JSON-encoded []string
	CreatedAt time.Time
}

type VersionChange struct {
	ID      int64
	Version string
	Action  string // "add", "modify", "delete"
	Path    string
}

type CreateVersionRequest struct {
	Version string   `json:"version" binding:"required"`
	Notes   []string `json:"notes" binding:"required"`
}

type ChangeEntry struct {
	Action string `json:"action" binding:"required"`
	Path   string `json:"path" binding:"required"`
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

// Client API DTOs

type ClientLatestResponse struct {
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

type ClientUpdateVersion struct {
	Version   string        `json:"version"`
	Timestamp int64         `json:"timestamp"`
	Changes   []ChangeEntry `json:"changes"`
}

type ClientUpdatesResponse struct {
	Versions []ClientUpdateVersion `json:"versions"`
}
