package upload

import "time"

type Handler interface {
	// UploadFile 上传文件至其他服务器
	UploadFile(fileName string, filePath string, modTime time.Time) error
}
