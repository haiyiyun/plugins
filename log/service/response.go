package service

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	id           string
	data         []byte
	writedHeader bool
}

func (lrw *ResponseWriter) Header() http.Header {
	return lrw.ResponseWriter.Header()
}

func (lrw *ResponseWriter) WriteHeader(statusCode int) {
	if !lrw.writedHeader {
		lrw.writedHeader = true
		lrw.ResponseWriter.WriteHeader(statusCode)
	}
}

func (lrw *ResponseWriter) Write(b []byte) (int, error) {
	buf := bytes.NewBuffer(lrw.data)
	buf.Write(b)
	lrw.data = buf.Bytes()

	return lrw.ResponseWriter.Write(b)
}

func (lrw *ResponseWriter) SetID(id string) {
	lrw.id = id
}

func (lrw *ResponseWriter) GetID() string {
	return lrw.id
}

func (lrw *ResponseWriter) GetData() []byte {
	return lrw.data
}
