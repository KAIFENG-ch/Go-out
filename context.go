package Go_out

import (
	"math"
	"net/http"
)

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

type ErrorType int64

type Error struct {
	Err  error
	Type ErrorType
	Meta interface{}
}

type errorMsg []*Error

const abortIndex int8 = math.MaxInt8 >> 1

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

func (c *Context) Handler() HandlerFunc {
	return c.Handlers.Last()
}

// Next next用于调用请求内部handler并传递给下一个handler
func (c *Context) Next() {
	c.Index++
	for c.Index < int8(len(c.Handlers)) {
		c.Handlers[c.Index](c)
		c.Index++
	}
}

func (c *Context) IsAborted() bool {
	return c.Index >= abortIndex
}

func (c *Context) Abort() {
	c.Index = abortIndex
}

func (r *responseWriter) Written() bool {
	return r.size != noWritten
}
