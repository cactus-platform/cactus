package response

import "net/http"

var (
	ErrBadRequest = APIError{
		Status:  http.StatusBadRequest,
		Code:    "BAD_REQUEST",
		Message: "Bad request",
	}

	ErrInternalServer = APIError{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
	}
)
