package upload

import (
	"PieFileMigrate/src/base"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func NewHTTPUploadHandler() Handler {
	base.InitHttpConfig()
	return &HTTPUploadHandler{}
}

type HTTPUploadHandler struct{}

func (u *HTTPUploadHandler) UploadFile(fileName string, filePath string, modTime time.Time) error {
	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)

	// 上传文件
	fw, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}
	fh, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = io.Copy(fw, fh)
	if err != nil {
		return err
	}

	// 校验令牌
	err = bodyWriter.WriteField("token", base.HttpConfig.Http.Token)
	if err != nil {
		return err
	}

	// 文件命名
	err = bodyWriter.WriteField("fileName", fileName)
	if err != nil {
		return err
	}

	// 上传路径
	err = bodyWriter.WriteField("filePath", filePath)
	if err != nil {
		return err
	}

	// 文件最后修改时间
	err = bodyWriter.WriteField("modTime", fmt.Sprintf("%v", modTime.UnixMilli()))
	if err != nil {
		return err
	}

	boundary := bodyWriter.Boundary()
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	requestReader := io.MultiReader(bodyBuf, fh, closeBuf)

	req, err := http.NewRequest("POST", base.HttpConfig.Http.TargetUrl, requestReader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("status code != 200")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		return err
	}

	if string(body) != "ok" {
		return errors.New("response body != ok")
	}

	// 只有当http状态码为200且返回数据为小写字符串'ok'时才会判定为上传成功
	return nil
}

func (u *HTTPUploadHandler) UploadFiles(fileNames []string, filePaths []string, modTimes []time.Time) error {

	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)

	boundary := bodyWriter.Boundary()

	// 校验令牌
	err := bodyWriter.WriteField("token", base.HttpConfig.Http.Token)
	if err != nil {
		return err
	}

	for i := 0; i < len(fileNames); i++ {
		// 上传文件
		fw, err := bodyWriter.CreateFormFile("files", fileNames[i])
		if err != nil {
			return err
		}

		fh, err := os.Open(filePaths[i])
		if err != nil {
			fh.Close()
			return err
		}
		_, err = io.Copy(fw, fh)
		if err != nil {
			fh.Close()
			return err
		}

		// 文件命名
		err = bodyWriter.WriteField("fileNames", fileNames[i])
		if err != nil {
			fh.Close()
			return err
		}

		// 上传路径
		err = bodyWriter.WriteField("filePaths", filePaths[i])
		if err != nil {
			fh.Close()
			return err
		}

		// 文件最后修改时间
		err = bodyWriter.WriteField("modTimes", fmt.Sprintf("%v", modTimes[i].UnixMilli()))
		if err != nil {
			fh.Close()
			return err
		}
		fh.Close()
	}

	requestReader := io.MultiReader(bodyBuf)

	req, err := http.NewRequest("POST", base.HttpConfig.Http.TargetUrl, requestReader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("status code != 200")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		return err
	}

	if string(body) != "ok" {
		return errors.New("response body != ok")
	}

	// 只有当http状态码为200且返回数据为小写字符串'ok'时才会判定为上传成功
	return nil
}
