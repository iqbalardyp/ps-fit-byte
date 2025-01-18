package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type FileUsecase struct {
	S3Uploader *manager.Uploader
}

const (
	JPEG = "image/jpeg"
	JPG  = "image/jpg"
	PNG  = "image/png"
)

var (
	AWS_S3_REGION      = os.Getenv("S3_REGION")
	AWS_S3_ID          = os.Getenv("S3_ID")
	AWS_S3_SECRET_KEY  = os.Getenv("S3_SECRET_KEY")
	AWS_S3_BUCKET_NAME = os.Getenv("S3_BUCKET_NAME")
	nameType           = map[string]string{
		JPEG: ".jpeg",
		JPG:  ".jpg",
		PNG:  ".png",
	}
)

func NewFileUseCase(uploader *manager.Uploader) *FileUsecase {
	return &FileUsecase{
		S3Uploader: uploader,
	}
}

func (c *FileUsecase) UploadFile(file multipart.File, fileType string) string {
	filename := c.generateFilename(fileType)
	fileUri := c.generateFileUri(filename)

	go func(uploader *manager.Uploader, file multipart.File, bucket, name string) {
		params := &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(name),
			Body:   file,
			ACL:    types.ObjectCannedACLPublicRead, // Allowed public read
		}

		_, err := uploader.Upload(context.Background(), params)
		if err != nil {
			// Log the error for diagnostics (optional)
			fmt.Printf("failed to upload file: %v\n", err)
		}
	}(c.S3Uploader, file, AWS_S3_BUCKET_NAME, filename)

	return fileUri
}

func (c *FileUsecase) generateFilename(fileType string) string {
	postfix := nameType[fileType]
	return uuid.New().String() + postfix
}

func (c *FileUsecase) generateFileUri(filename string) string {
	return fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		AWS_S3_BUCKET_NAME,
		AWS_S3_REGION,
		filename,
	)
}
