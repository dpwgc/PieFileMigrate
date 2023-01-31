package util

import (
	"io/ioutil"
	"os"
	"time"
)

// WriteFile 写入文件
func WriteFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0766)
	return err
}

// ReadFile 读取文件
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	return data, err
}

// FileListNodeModel 文件列表节点
type FileListNodeModel struct {
	IsDir   bool      `json:"isDir"`
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

// ReadDir 读取目录
func ReadDir(path string) ([]FileListNodeModel, error) {
	fileInfoList, err := ioutil.ReadDir(path)
	var nodeList []FileListNodeModel
	for _, v := range fileInfoList {
		node := FileListNodeModel{
			IsDir:   v.IsDir(),
			Name:    v.Name(),
			Size:    v.Size(),
			ModTime: v.ModTime(),
		}
		nodeList = append(nodeList, node)
	}
	return nodeList, err
}
