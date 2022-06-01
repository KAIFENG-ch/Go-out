package Go_out

import (
	"sync"
)

type HandlerFunc func(*Context)

type HandlerChain []HandlerFunc

type Param struct {
	Key   string
	Value string
}

type params []Param

type Engine struct {
	RouterGroup
	pool      sync.Pool
	maxParams uint16
}

func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			bashPath: "/",
			root:     true,
		},
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() any {
		return engine.AllocateContext()
	}
	return engine
}

func Default() *Engine {
	engine := new(Engine)
	return engine
}

func (e *Engine) AllocateContext() func(context *Context) {
	//v := make(params, 0, e.maxParams)
	return func(context *Context) {}
}
