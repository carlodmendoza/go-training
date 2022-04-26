package http

import (
	"encoding/json"
	"net/http"
)

type ResponseWriter struct {
	w http.ResponseWriter
}

func (rw *ResponseWriter) WriteMessage(data string) (int, error) {
	resp, _ := json.Marshal(map[string]string{"message": data})
	rw.w.Header().Set("Content-Type", "application/json")

	return rw.w.Write(resp)
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.w.WriteHeader(statusCode)
}

func (rw *ResponseWriter) Writer() http.ResponseWriter {
	return rw.w
}
