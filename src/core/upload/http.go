package upload

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func NewHTTPUploadHandler() Handler {
	base.InitHttpConfig()
	return &HTTPUploadHandler{}
}

type HTTPUploadHandler struct{}

func (u *HTTPUploadHandler) UploadFile(fileName string, localFilePath string) error {
	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)

	// 上传文件
	_, err := bodyWriter.CreateFormFile(fileName, localFilePath)
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
	err = bodyWriter.WriteField("filePath", localFilePath)
	if err != nil {
		return err
	}

	// the file data will be the second part of the body
	fh, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	// need to know the boundary to properly close the part myself.
	boundary := bodyWriter.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	requestReader := io.MultiReader(bodyBuf, fh, closeBuf)
	fi, err := fh.Stat()
	if err != nil {
		return err
	}

	// fmt.Printf("迁移本地文件 [ %s ]\n", localFilePath)
	base.LogHandler.Printf("%s 迁移本地文件 [ %s ]\n", constant.LogInfoTag, localFilePath)

	req, err := http.NewRequest("POST", base.HttpConfig.Http.Url, requestReader)
	if err != nil {
		return err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("status code != 200")
	}

	return nil
}
