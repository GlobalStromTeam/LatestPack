package models

import "time"

type Channel struct {
	ID        string
	Name      string
	Type      string // "local", "webdav", "s3"
	Enabled   bool
	Weight    int
	Config    string // JSON-encoded ChannelConfig
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChannelConfig struct {
	Endpoint  string `json:"endpoint,omitempty"`
	Path      string `json:"path,omitempty"`
	AccessKey string `json:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty"`
	Bucket    string `json:"bucket,omitempty"`
	Region    string `json:"region,omitempty"`
	Mode      string `json:"mode,omitempty"` // "proxy" (default) or "openlist"
}

type CreateChannelRequest struct {
	Name    string        `json:"name" binding:"required"`
	Type    string        `json:"type" binding:"required"`
	Enabled *bool         `json:"enabled"`
	Weight  *int          `json:"weight"`
	Config  ChannelConfig `json:"config" binding:"required"`
}

type UpdateChannelRequest struct {
	Name    *string        `json:"name"`
	Enabled *bool          `json:"enabled"`
	Weight  *int           `json:"weight"`
	Config  *ChannelConfig `json:"config"`
}

type ChannelDTO struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Type    string        `json:"type"`
	Enabled bool          `json:"enabled"`
	Weight  int           `json:"weight"`
	Config  ChannelConfig `json:"config"`
}
