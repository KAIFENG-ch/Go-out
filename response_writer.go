package Go_out

import (
	"io"
	"net/http"
	"strconv"
)

type ResponseWriter interface {
	Header() http.Header

	Write([]byte) (int, error)

	WriteHeader(statusCode int)

	WriteHeaderNow()

	Status() int

	WriteString(int) (int, error)

	Written() bool

	Pusher() http.Pusher
}

type responseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func (r *responseWriter) Status() int {
	return r.status
}

func (r *responseWriter) WriteString(s int) (n int, err error) {
	r.WriteHeaderNow()
	n, err = io.WriteString(r.ResponseWriter, strconv.Itoa(s))
	r.size += n
	return
}

func (r *responseWriter) Pusher() http.Pusher {
	if pusher, ok := r.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}

func (r *responseWriter) reset(writer http.ResponseWriter) {
	r.ResponseWriter = writer
	r.size = noWritten
	r.status = defaultStatus
}

func (r *responseWriter) WriteHeader(code int) {
	if code > 0 && r.status != code {
		if r.Written() {
		}
		r.status = code
	}
}

func (r *responseWriter) WriteHeaderNow() {
	if !r.Written() {
		r.size = 0
		r.ResponseWriter.WriteHeader(r.status)
	}
}

func (r *responseWriter) Write(data []byte) (n int, err error) {
	r.WriteHeaderNow()
	n, err = r.ResponseWriter.Write(data)
	r.size += n
	return
}
