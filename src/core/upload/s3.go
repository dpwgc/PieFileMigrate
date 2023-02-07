package upload

import (
	"PieFileMigrate/src/base"
	"time"
)

func NewS3UploadHandler() Handler {
	base.InitFtpConfig()
	return &S3UploadHandler{}
}

type S3UploadHandler struct{}

func (u *S3UploadHandler) UploadFiles(fileNames []string, filePaths []string, modTimes []time.Time) error {
	return nil
}
