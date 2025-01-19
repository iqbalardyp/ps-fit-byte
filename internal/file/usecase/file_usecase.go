package file_usecase

import (
	"context"
	file_dto "fit-byte/internal/file/dto"
	"fit-byte/pkg/dotenv"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type FileUsecase struct {
	S3Uploader *manager.Uploader
	Env        *dotenv.Env
}

const (
	JPEG = "image/jpeg"
	JPG  = "image/jpg"
	PNG  = "image/png"
)

var (
	nameType = map[string]string{
		JPEG: ".jpeg",
		JPG:  ".jpg",
		PNG:  ".png",
	}
)

func NewFileUseCase(uploader *manager.Uploader, env *dotenv.Env) *FileUsecase {
	return &FileUsecase{
		S3Uploader: uploader,
		Env:        env,
	}
}

func (u *FileUsecase) UploadFile(file multipart.File, fileType string) (*file_dto.FileUploadResponse, error) {
	var response file_dto.FileUploadResponse
	defer file.Close()

	filename := u.generateFilename(fileType)
	response.FileUrl = u.generateFileUrl(filename)

	go func(uploader *manager.Uploader, file multipart.File, bucket, name string) {
		params := &s3.PutObjectInput{
			Bucket: aws.String(u.Env.AWS_S3_BUCKET_NAME),
			Key:    aws.String(filename),
			ACL:    types.ObjectCannedACLPublicRead,
			Body:   file,
		}
		_, err := uploader.Upload(context.Background(), params)
		if err != nil {
			fmt.Printf("failed to upload file: %v\n", err)
		}
	}(u.S3Uploader, file, u.Env.AWS_S3_BUCKET_NAME, filename)

	return &response, nil
}

func (c *FileUsecase) generateFilename(fileType string) string {
	postfix := nameType[fileType]
	return uuid.New().String() + postfix
}

func (c *FileUsecase) generateFileUrl(filename string) string {
	return fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		c.Env.AWS_S3_BUCKET_NAME,
		c.Env.AWS_S3_REGION,
		filename,
	)
}
