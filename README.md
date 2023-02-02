# PieFileMigrate 
## 文件迁移脚本

***

### 简介

PieFileMigrate是一个简易的小型文件批量迁移脚本，可将服务器内的一个或多个目录内的全部文件通过HTTP或FTP（未实现）的方式定时增量/全量同步上传至其他服务器。

***

### 原理

`TODO`

***

### 【示例1】以HTTP方式迁移文件

* 将本地`./data`目录下的文件增量同步上传至`http://127.0.0.1:8011/demo/uploadTargetApi`接口
* 同步周期为每30秒一次，每次仅同步在当前时间前`172800`秒内新增/更新的文件

#### 填写脚本应用配置 config/application.yaml

```yaml
application:

  # http服务端口（用于监控脚本运行情况）
  server-port: 8012
  # 文件迁移方式（http、ftp）
  migrate-mode: http
  # 文件上传记录存储介质（boltdb、redis）
  storage-media: boltdb

  # 迁移工作者数组（可同时迁移多个目录）
  workers:
      # 需要进行迁移的源文件目录（两个工作者不能同时迁移一个目录）
    - source-path: ./data
      # 增量迁移文件时间限制（单位：秒，只迁移{migrate-file-time-limit}秒内更新的文件）
      migrate-file-age-limit: 172800
      # 定时迁移任务的cron表达式
      job-cron: 0/30 * * * * ?

  # 内置mq
  mq:
    # 内置mq容量大小
    max-size: 1000
    # 内置mq消费者数量（文件迁移并发数）
    consumer-num: 5
```

#### 填写HTTP上传配置 config/http.yaml

```yaml
http:
  # 文件上传的目标地址
  target-url: http://127.0.0.1:8011/demo/uploadTargetApi
```

#### 自定义文件接收接口示例

* 需要接收5个参数，分别是：`fileName` 文件名 `filePath` 文件原路径 `modTime` 文件最近修改时间 `token` 权限校验令牌（自定义校验逻辑） `file` 文件
* 需要自行编写接收逻辑（例如把文件传入对象存储或者云盘里）

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

#### 迁移工作者运行状态监控

* 监控接口地址 `/worker/monitor`

> `GET` http://127.0.0.1:8012/worker/monitor

* 接口返回数据

```json
{
  "./data": {
    "lastMigrateStartTime": "2023-02-02T14:21:30.00088+08:00",
    "lastMigrateEndTime": "2023-02-02T14:21:30.001168+08:00"
  }
}
```

* 返回字段说明：`./data` 该工作者负责的目录 `lastMigrateStartTime` 最近一次迁移开始时间 `lastMigrateEndTime` 最近一次迁移结束时间
