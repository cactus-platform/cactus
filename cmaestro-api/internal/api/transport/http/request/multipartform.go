package request

import (
	"mime/multipart"
	"net/http"
)

func WithMultipartFile(
	r *http.Request,
	key string,
	maxMemory int64,
	fn func(file multipart.File, header *multipart.FileHeader) (any, error),
) (any, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, err
	}

	file, header, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return fn(file, header)
}
