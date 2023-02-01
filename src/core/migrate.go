package core

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"PieFileMigrate/src/util"
	"fmt"
	"time"
)

// doMigrate 执行迁移操作
func doMigrate() {

	defer func() {
		err := recover()
		if err != nil {
			base.LogHandler.Println(constant.LogErrorTag, err)
		}
	}()

	base.LogHandler.Printf("%s 扫描目录: [ %s ]\n", constant.LogInfoTag, base.ApplicationConfig.Application.SourcePath)

	var children []fileTreeNodeModel
	root := fileTreeNodeModel{
		IsDir:    true,
		Path:     base.ApplicationConfig.Application.SourcePath,
		Children: children,
	}
	dfsFileTree(&root)
}

// 文件树节点
type fileTreeNodeModel struct {
	IsDir    bool                `json:"isDir"`
	Name     string              `json:"name"`
	Path     string              `json:"path"`
	Size     int64               `json:"size"`
	ModTime  time.Time           `json:"modTime"`
	Children []fileTreeNodeModel `json:"children"`
}

// 递归遍历文件树
func dfsFileTree(node *fileTreeNodeModel) {
	// 如果是文件
	if !node.IsDir {
		// 如果该文件没有标记或标记过期 && 更新时间在限制时间内
		if storageHandler.CheckFile(node.Path, node.ModTime) && node.ModTime.Unix() > (time.Now().Unix()-base.ApplicationConfig.Application.MigrateFileTimeLimit) {
			// 异步迁移文件至其他服务器
			asyncMigrateFile(node.Name, node.Path, node.ModTime)
			base.LogHandler.Printf("%s 迁移文件 [ %s ]\n", constant.LogInfoTag, node.Path)
		}
		return
	}
	list, err := util.ReadDir(node.Path)
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, "文件目录读取失败", err)
		return
	}
	if len(list) > 0 {
		var children []fileTreeNodeModel
		for _, v := range list {
			childPath := fmt.Sprintf("%s/%s", node.Path, v.Name)
			child := fileTreeNodeModel{
				IsDir:    v.IsDir,
				Name:     v.Name,
				Size:     v.Size,
				Path:     childPath,
				ModTime:  v.ModTime,
				Children: children,
			}
			dfsFileTree(&child)
			node.Children = append(node.Children, child)
		}
	}
}
