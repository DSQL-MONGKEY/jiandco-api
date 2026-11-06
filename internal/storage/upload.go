package storage

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	MaxAvatarSize       = 10 * 1024 * 1024
	MaxProductImageSize = 10 * 1024 * 1024
	MaxReceiptSize      = 10 * 1024 * 1024
)

var (
	ErrFileTooLarge    = errors.New("file is too large")
	ErrUnsupportedType = errors.New("unsupported file type")
	ErrEmptyFile       = errors.New("file is empty")
)

var ImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

var ReceiptTypes = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"application/pdf": ".pdf",
}

type ValidatedFile struct {
	File        multipart.File
	ContentType string
	Extension   string
	Size        int64
}

func OpenValidated(fileHeader *multipart.FileHeader, allowed map[string]string, maxSize int64) (*ValidatedFile, error) {
	if fileHeader == nil {
		return nil, ErrEmptyFile
	}
	if fileHeader.Size <= 0 {
		return nil, ErrEmptyFile
	}
	if fileHeader.Size > maxSize {
		return nil, ErrFileTooLarge
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 512)
	n, readErr := file.Read(buffer)
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		_ = file.Close()
		return nil, readErr
	}

	contentType := http.DetectContentType(buffer[:n])
	extension, ok := allowed[contentType]
	if !ok {
		_ = file.Close()
		return nil, ErrUnsupportedType
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		_ = file.Close()
		return nil, err
	}

	return &ValidatedFile{
		File:        file,
		ContentType: contentType,
		Extension:   extension,
		Size:        fileHeader.Size,
	}, nil
}
