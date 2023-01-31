package storage

type Handler interface {
	// MarkFile 标记已经上传过的文件
	MarkFile(filePath string) bool
	// CheckFile 检查文件是否已被标记
	CheckFile(filePath string) bool
}
