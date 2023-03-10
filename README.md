# PieFileMigrate 
## 一个基于Go开发的小型文件批量迁移脚本
### `待完善`

***

### 简介

PieFileMigrate是一个基于Go开发的小型文件批量迁移脚本，可将服务器内的一个或多个目录内的全部文件通过HTTP、S3（已实现，待优化）、FTP（未实现）的方式定时增量/全量同步上传至其他服务器。

***

### 说明

* `文件迁移流程` ：定时扫描本地目录下的所有文件，将未被标记的文件上传至指定服务器中，如果上传成功，则标记该文件。
* `聚合发送文件` : 为减少网络IO压力，使用HTTP方式推送文件时，可将多个文件写入一个HTTP请求里推送。
* `如何标记文件` ：如果一个文件上传成功，则将它的完整路径及它的最近修改时间写入KV数据库中（可以选择用本地BoltDB存储或者用线上Redis存储），每次上传文件前都会检查该文件路径是否已经在KV数据库中，如果已存在就说明该文件已经被上传过了，无需重复上传。
* `增量迁移方式` ：自行修改 config/application.yaml 配置文件中的 migrate-file-age-limit 配置项，迁移目录下最近一段时间内更新的文件。
* `全量迁移方式` ：将 config/application.yaml 配置文件中的 migrate-file-age-limit 设成0，迁移目录下的所有文件。
* `重新全量迁移` ：如果想重新全量迁移整个目录，可以先停用脚本，将标记用的KV数据库表删除/清空，重新运行脚本即可。

***

### 【HTTP迁移示例】以HTTP方式迁移文件至指定服务接口

* 将本地 `./data` 目录下的文件增量迁移上传至 `http://127.0.0.1:8011/demo/uploadTargetBatchApi` 接口
* 迁移周期为每60秒一次，每次仅迁移在当前时间前`172800`秒内新增/更新的文件，分配五个消费者并发推送文件，单次HTTP请求推送100个文件。

#### 填写脚本应用配置 config/application.yaml

```yaml
application:

  # http服务端口（用于监控脚本运行情况）
  server-port: 8012
  # 文件迁移方式（http、ftp、s3）
  migrate-mode: http
  # 文件上传记录存储介质（boltdb、redis）
  storage-media: boltdb

  # 迁移工作者数组（可同时迁移多个目录）
  workers:
    # 需要进行迁移的源文件目录（两个工作者不能同时迁移一个目录）
    - source-path: ./data
      # 增量迁移文件时间限制（单位：秒，只迁移{migrate-file-time-limit}秒内更新的文件，小于等于0则全部迁移）
      migrate-file-age-limit: 172800
      # 定时迁移任务的cron表达式
      job-cron: 0/60 * * * * ?

  # 内置mq
  mq:
    # 内置mq容量大小
    max-size: 10000000
    # 内置mq单批次消费消息数量（必须大于0）
    consume-batch: 100
    # 内置mq消费者数量（文件迁移并发数）
    consumer-num: 5

  # 日志配置
  log:
    # 日志文件最长存活时间（单位：天）
    maxAge: 7
```

#### 填写HTTP上传配置 config/http.yaml

```yaml
http:
  # 文件上传的目标地址
  target-url: http://127.0.0.1:8011/demo/uploadTargetBatchApi
```

#### 自定义文件接收接口示例

* 该接口的 `Content-Type` 必须为 `multipart/form-data` ，需要接收5个参数

| 参数名       | 参数类型     | 参数说明       |
|-----------|----------|------------|
| token     | string   | 权限校验令牌     |
| fileNames | string[] | 文件名称列表     |
| filePaths | string[] | 文件原路径列表    |
| modTimes  | string[] | 文件最近修改时间列表 |
| files     | file[]   | 文件列表       |

* 需要自行编写接收逻辑，以下是示例接口代码

```
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	// 启动HTTP服务
	r := gin.Default()
	r.POST("/demo/uploadTargetBatchApi", uploadTargetBatchApi)
	//加载端口号
	_ = r.Run(":8011")
}

func uploadTargetBatchApi(c *gin.Context) {

	// 身份校验token（自行处理校验逻辑）
	token := c.PostForm("token")
	fmt.Println("token: ", token)

	// 文件名称列表
	fileNames := c.PostFormArray("fileNames")
	fmt.Println("fileNames: ", fileNames)

	// 文件原路径列表
	filePaths := c.PostFormArray("filePaths")
	fmt.Println("filePaths: ", filePaths)

	// 文件最后修改时间列表
	modTimes := c.PostFormArray("modTimes")
	fmt.Println("modTimes: ", modTimes)

	// 文件列表
	form, _ := c.MultipartForm()
	files := form.File["files"]
	
	// 遍历文件列表
	for i := 0; i < len(files); i++ {
		//保存文件到本地
		if err := c.SaveUploadedFile(files[i], "./data/"+fileNames[i]); err != nil {
			c.String(http.StatusBadRequest, "fail")
			fmt.Println("SaveUploadedFile: ", err)
			return
		}
	}

	// 注意，必须要确保文件保存成功后才返回 code 200
	c.String(http.StatusOK, "ok")
}
```

***

### 【S3示例】以S3方式迁移文件至指定对象存储桶

#### 填写脚本应用配置 config/application.yaml

```yaml
# 参考上文HTTP示例中的应用配置
migrate-mode: s3
```

#### 填写S3上传配置 config/s3.yaml

```yaml
s3:
  endpoint: http://xxx.xxx.cn
  region: default
  bucket: test
  access-key: xxxxxx
  secret-key: xxxxxxxxx
  path-style: true
```

***

### 迁移工作者运行状态监控

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