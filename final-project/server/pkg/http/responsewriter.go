package http

import (
	"encoding/json"
	"net/http"
)

type ResponseWriter struct {
	w http.ResponseWriter
}

func (rw *ResponseWriter) Write(data []byte) (int, error) {
	rw.w.Header().Set("Content-Type", "application/json")

	return rw.w.Write(data)
}

func (rw *ResponseWriter) WriteMessage(msg string) (int, error) {
	resp, _ := json.Marshal(map[string]string{"message": msg})
	rw.w.Header().Set("Content-Type", "application/json")

	return rw.w.Write(resp)
}

func (rw *ResponseWriter) WriteError(msg string) (int, error) {
	resp, _ := json.Marshal(map[string]string{"error": msg})
	rw.w.Header().Set("Content-Type", "application/json")

	return rw.w.Write(resp)
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.w.WriteHeader(statusCode)
}

func (rw *ResponseWriter) Writer() http.ResponseWriter {
	return rw.w
}
