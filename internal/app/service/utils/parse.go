package utils

import (
	"fmt"
	"log/slog"
	"mime/multipart"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	FileType struct {
		Type     string
		MimeType string
		Size     int64
	}
)

func ValidateFileHeader(fileHeader *multipart.FileHeader, fileTypes []FileType, keyName string) (err error) {
	filenameSplit := strings.Split(fileHeader.Filename, ".")
	fileType := strings.ToLower(filenameSplit[len(filenameSplit)-1])
	mimeType := fileHeader.Header.Get("Content-Type")
	for _, allowedType := range fileTypes {
		if fileType == allowedType.Type {
			if allowedType.Size != 0 && fileHeader.Size > allowedType.Size<<20 {
				err = fmt.Errorf("Error:Field validation for %s. File is too large", keyName)
			}
			return
		}
		if mimeType == allowedType.MimeType {
			if allowedType.Size != 0 && fileHeader.Size > allowedType.Size<<20 {
				err = fmt.Errorf("Error:Field validation for %s. File is too large", keyName)
			}
			return
		}
	}
	err = fmt.Errorf("Error:Field validation for %s. Invalid filetype", keyName)
	return
}

func ParseMultipartForm(c echo.Context, fileTypes []FileType) (files map[string][]*multipart.FileHeader, values map[string][]string, err error) {
	ctx := c.Request().Context()
	err = c.Request().ParseMultipartForm(0)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parsemultipartformfile", "error", err)
		err = fmt.Errorf("Error:Field validation for parse multipart form")
		return
	}
	files = c.Request().MultipartForm.File
	values = c.Request().MultipartForm.Value
	if len(files) < 1 && len(values) < 1 && len(fileTypes) > 6 && len(values) > 6 {
		err = fmt.Errorf("Error:Field validation for parse multipart form")
		slog.ErrorContext(ctx, "failed to parsemultipartformfile", "error", err)
		return
	}
	if len(fileTypes) > 0 {
		for keyName, valuess := range files {
			for _, fileHeader := range valuess {
				err = ValidateFileHeader(fileHeader, fileTypes, keyName)
				if err != nil {
					err = fmt.Errorf("%s, %s", keyName, err)
					slog.ErrorContext(ctx, "failed to parsemultipartformfile", "error", err)
					return
				}
			}
		}
	}
	return
}
