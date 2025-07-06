package models

type FileMetaData struct {
	Filename string `json:"filename"`
	// SHA256   string `json:"sha256"`
	Size int64 `json:"size"`
}
