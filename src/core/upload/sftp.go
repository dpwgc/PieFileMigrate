package upload

func NewSFTPUploadHandler() Handler {
	return &SFTPUploadHandler{}
}

type SFTPUploadHandler struct{}

func (u *SFTPUploadHandler) UploadFile(fileName string, localFilePath string) error {
	return nil
}
