package models

type VersionFileSnapshot struct {
	ID      int64
	Version string
	Path    string
	Hash    string
}
