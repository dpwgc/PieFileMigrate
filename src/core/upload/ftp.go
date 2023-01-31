package upload

func NewFTPUploadHandler() Handler {
	return &FTPUploadHandler{}
}

type FTPUploadHandler struct{}

func (u *FTPUploadHandler) UploadFile(fileName string, localFilePath string) error {
	return nil
}
