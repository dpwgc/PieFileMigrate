package upload

type Handler interface {
	// UploadFile 上传文件至其他服务器
	UploadFile(fileName string, localFilePath string) error
}
