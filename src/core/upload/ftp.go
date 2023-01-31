package upload

import "PieFileMigrate/src/base"

func NewFTPUploadHandler() Handler {
	base.InitFtpConfig()
	return &FTPUploadHandler{}
}

type FTPUploadHandler struct{}

func (u *FTPUploadHandler) UploadFile(fileName string, localFilePath string) error {
	return nil
}
