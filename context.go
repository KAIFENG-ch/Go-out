package Go_out

import "net/http"

type ResponseWriter interface {
	Header() http.Header

	WriteHeaderNow()
}

type responseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func (r responseWriter) WriteHeaderNow() {

}

type ErrorType int64

type Error struct {
	Err  error
	Type ErrorType
	Meta interface{}
}

type errorMsg []*Error

type Context struct {
	writermem    responseWriter
	request      *http.Request
	Writer       ResponseWriter
	Handlers     HandlerChain
	engine       *Engine
	params       *params
	Params       params
	Keys         map[string]interface{}
	skippedNodes *[]skippedNode
	FullPath     string
	Index        int8
	Errors       errorMsg
}

func (c *Context) Reset() {
	c.Writer = &c.writermem
	c.Params = c.Params[:0]
	c.Handlers = nil
	c.Index = -1

	c.FullPath = ""
	c.Keys = nil
	c.Errors = c.Errors[:0]
	*c.params = (*c.params)[:0]
	*c.skippedNodes = (*c.skippedNodes)[:0]
}

// Next next用于调用请求内部handler并传递给下一个handler
func (c *Context) Next() {
	c.Index++
	for c.Index < int8(len(c.Handlers)) {
		c.Handlers[c.Index](c)
		c.Index++
	}
}
