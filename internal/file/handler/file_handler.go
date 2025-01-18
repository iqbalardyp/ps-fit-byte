package handler

import (
	"fit-byte/internal/file/dto"
	"fit-byte/internal/file/usecase"
	customErrors "fit-byte/pkg/custom-errors"
	"fit-byte/pkg/response"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FileHandler struct {
	Log     *logrus.Logger
	Usecase *usecase.FileUsecase
}

func NewFileHandler(Usecase *usecase.FileUsecase, logger *logrus.Logger) *FileHandler {
	return &FileHandler{Log: logger,
		Usecase: Usecase}
}

func (c *FileHandler) UploadFile(ctx echo.Context) error {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	file, err := fileHeader.Open()
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	fileType, isValid := c.isValidFile(fileHeader, file)
	if !isValid {
		err = errors.Wrap(customErrors.ErrBadRequest, "file is invalid")
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	defer file.Close()

	fileUri := c.Usecase.UploadFile(file, *fileType)

	return ctx.JSON(http.StatusOK, &dto.FileUploadResponse{
		FileUri: fileUri,
	})
}

func (c *FileHandler) isValidFile(fileHeader *multipart.FileHeader, file multipart.File) (*string, bool) {

	if fileHeader.Size > 100*1024 {
		return nil, false
	}

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return nil, false
	}
	// Reset the read pointer of the file
	if _, err := file.Seek(0, 0); err != nil {
		return nil, false
	}
	fileType := http.DetectContentType(buffer)

	switch fileType {
	case usecase.JPEG, usecase.JPG, usecase.PNG:
		return &fileType, true
	default:
		return nil, false
	}
}
