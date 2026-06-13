package request

import (
	"mime/multipart"
	"net/http"
)

func WithMultipartFile(
	r *http.Request,
	key string,
	maxMemory int64,
	fn func(file multipart.File, header *multipart.FileHeader) error,
) error {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return err
	}

	file, header, err := r.FormFile(key)
	if err != nil {
		return err
	}
	defer file.Close()

	return fn(file, header)
}
