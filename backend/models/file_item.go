package models

type FileItem struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size string `json:"size"`
	Date string `json:"date"`
}

type FileListResponse struct {
	Path  string     `json:"path"`
	Items []FileItem `json:"items"`
}

type CreateFolderRequest struct {
	Path string `json:"path"`
	Name string `json:"name" binding:"required"`
}

type RenameRequest struct {
	Path    string `json:"path"`
	OldName string `json:"oldName" binding:"required"`
	NewName string `json:"newName" binding:"required"`
}

type DeleteItemRequest struct {
	Path string `json:"path"`
	Name string `json:"name" binding:"required"`
}
