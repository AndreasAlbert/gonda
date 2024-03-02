package metadata

import (
	"time"
)

type UploadStatus string

const (
	Pending UploadStatus = "PENDING"
	Done    UploadStatus = "DONE"
)

type IndexEntry struct {
	Arch        string   `json:"arch"`
	Build       string   `json:"build"`
	BuildNumber int      `json:"build_number"`
	Depends     []string `json:"depends"`
	License     string   `json:"license"`
	Name        string   `json:"name"`
	Platform    string   `json:"platform"`
	Subdir      string   `json:"subdir"`
	Timestamp   int      `json:"timestamp"`
	Version     string   `json:"version"`
}
type FileUpload struct {
	FileName   string
	UploadTime time.Time
	Status     UploadStatus
}

type File struct {
	Name     string
	Arch     string
	Build    string
	Platform string
	Sha256   string
	Md5      string
	Version  string
	Depends  []string
}
type MetaDataStore interface {
	CreateFileUpload(FileName string, UploadTime time.Time) error
	GetFileUploadByName(FileName string) (FileUpload, error)
	UpdateFileUploadStatus(FileName string, Status UploadStatus) error
}
