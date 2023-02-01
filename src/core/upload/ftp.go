package upload

import (
	"PieFileMigrate/src/base"
	"time"
)

func NewFTPUploadHandler() Handler {
	base.InitFtpConfig()
	return &FTPUploadHandler{}
}

type FTPUploadHandler struct{}

func (u *FTPUploadHandler) UploadFile(fileName string, filePath string, modTime time.Time) error {
	return nil
}
