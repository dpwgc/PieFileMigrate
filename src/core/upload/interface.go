package upload

import "time"

type Handler interface {
	// UploadFiles 批量上传文件至其他服务器
	UploadFiles(fileNames []string, filePaths []string, modTimes []time.Time) error
}
