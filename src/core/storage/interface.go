package storage

import "time"

type Handler interface {
	// MarkFile 标记已经上传过的文件
	MarkFile(filePath string, modTime time.Time) bool
	// CheckFile 检查文件是否已被标记
	CheckFile(filePath string, modTime time.Time) bool
}
