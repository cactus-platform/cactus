package response

import "net/http"

func OK(w http.ResponseWriter, data any) {
	Write(w, http.StatusOK, Envelope{
		Data: data,
	})
}

func Created(w http.ResponseWriter, data any) {
	Write(w, http.StatusCreated, Envelope{
		Data: data,
	})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
