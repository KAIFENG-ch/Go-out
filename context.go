package Go_out

import "net/http"

type Context struct {
	request *http.Request
	Writer http.ResponseWriter
}

func (c *Context) Handler() HandlerFunc {
	panic("")
}
