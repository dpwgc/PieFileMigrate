package upload

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/util"
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
)

func NewS3UploadHandler() Handler {
	base.InitS3Config()
	return &S3UploadHandler{
		S3: newS3Client(),
	}
}

type S3UploadHandler struct {
	S3 *s3.S3
}

func (u *S3UploadHandler) UploadFiles(fileNames []string, filePaths []string, modTimes []time.Time) error {
	for i := 0; i < len(fileNames); i++ {

		file, err := util.ReadFile(filePaths[i])
		if err != nil {
			return err
		}

		_, err = u.S3.PutObjectWithContext(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(base.S3Config.S3.Bucket),
			Key:    aws.String(filePaths[i]),
			Body:   bytes.NewReader(file),
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
				// If the SDK can determine the request or retry delay was canceled
				// by a context the CanceledErrorCode error code will be returned.
				return fmt.Errorf("upload canceled due to timeout, %v", err)
			}
			return fmt.Errorf("failed to upload object, %v", err)
		}
	}
	return nil
}

func newS3Client() *s3.S3 {
	creds := credentials.NewStaticCredentials(base.S3Config.S3.AccessKey, base.S3Config.S3.SecretKey, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String("default"),
		Endpoint:         &base.S3Config.S3.Endpoint,
		S3ForcePathStyle: &base.S3Config.S3.PathStyle, // 因为使用的IP:Port/bucket的形式，使用path风格
		Credentials:      creds,
	}))
	svc := s3.New(sess)
	return svc
}
