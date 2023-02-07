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

func (u *FTPUploadHandler) UploadFiles(fileNames []string, filePaths []string, modTimes []time.Time) error {
	return nil
}
