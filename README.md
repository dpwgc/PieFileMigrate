# PieFileMigrate 
## 文件迁移脚本

***

### 以HTTP方式迁移文件

#### PieFileMigrate 配置文件 http.yaml

```yaml
http:
  # 文件上传的目标地址
  target-url: http://127.0.0.1:8011/demo/uploadTargetApi
```

#### 自定义文件接收接口示例

```
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 启动HTTP服务
	r := gin.Default()
	r.POST("/demo/uploadTargetApi", uploadTargetApi)
	//加载端口号
	_ = r.Run(":8011")
}

func uploadTargetApi(c *gin.Context) {

	// 文件原命名
	fileName := c.PostForm("fileName")
	fmt.Println("fileName: ", fileName)

	// 文件原路径
	filePath := c.PostForm("filePath")
	fmt.Println("filePath: ", filePath)

	// 文件最后修改时间
	modTime := c.PostForm("modTime")
	fmt.Println("modTime: ", modTime)

	// 身份校验token（自行处理校验逻辑）
	token := c.PostForm("token")
	fmt.Println("token: ", token)

	// 接收文件数据
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "fail")
		fmt.Println("FormFile: ", err)
		return
	}
	
	//保存文件到本地
	if err := c.SaveUploadedFile(file, "./data/"+fileName); err != nil {
		c.String(http.StatusBadRequest, "fail")
		fmt.Println("SaveUploadedFile: ", err)
		return
	}
	
	// 注意，必须要确保文件保存成功后才返回 code 200
	c.String(http.StatusOK, "ok")
}
```