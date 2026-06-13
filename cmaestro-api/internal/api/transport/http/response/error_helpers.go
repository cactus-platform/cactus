package response

import "net/http"

func Fail(
	w http.ResponseWriter,
	err APIError,
) {
	Write(w, err.Status, Envelope{
		Error: &err,
	})
}
