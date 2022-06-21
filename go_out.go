package Go_out

import (
	"net/http"
	"path"
	"sync"
)

type HandlerFunc func(*Context)

type HandlerChain []HandlerFunc

var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
)

//Param param是url参数，由键值对构成
type Param struct {
	Key   string
	Value string
}

type params []Param

//Engine engine是gin的引擎
type Engine struct {
	RouterGroup
	pool                   sync.Pool
	maxParams              uint16
	tree                   MethodTree
	allNoMethod            HandlerChain
	allNoRoute             HandlerChain
	UseRawPath             bool
	UnescapePathValue      bool
	RemoveExtraSlash       bool
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
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

func (e *Engine) AllocateContext() *Context {
	v := make(params, 0, e.maxParams)
	skippedNode := make([]skippedNode, 0)
	return &Context{
		engine:       e,
		params:       &v,
		skippedNodes: &skippedNode,
	}
}

//func (e *Engine) Run(addr ...string) (err error) {
//	defer func() {}()
//address :=
//err = http.ListenAndServe(address, e.Handlers)
//return
//}

//ServeHttp 处理程序接口
func (e *Engine) ServeHttp(w http.ResponseWriter, req *http.Request) {
	c := e.pool.Get().(*Context)
	c.request = req
	c.Reset()

	e.handleHttpRequest(c)
	e.pool.Put(c)

}

// handleHttpRequest 接收http請求
func (e *Engine) handleHttpRequest(c *Context) {
	httpMethod := c.request.Method
	rPath := c.request.URL.Path
	unescape := false
	if e.UseRawPath && len(c.request.URL.RawPath) == 0 {
		rPath = c.request.URL.RawPath
		unescape = e.UnescapePathValue
	}

	if e.RemoveExtraSlash {
		rPath = cleanPath(rPath)
	}
	t := e.tree
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// 在路由樹中尋找路由
		value := root.getValues(rPath, c.params, c.skippedNodes, unescape)
		if value.params != nil {
			c.Params = *value.params
		}
		if value.handlers != nil {
			c.Handlers = value.handlers
			c.FullPath = value.fullPath
			c.Next()
			c.writermem.WriteHeaderNow()
			return
		}
		if httpMethod != http.MethodConnect && rPath != "/" {
			if value.tsr && e.RedirectTrailingSlash {
				redirectTrailingSlash(c)
				return
			}
			if e.RedirectFixedPath && redirectFixedPath(c, root, e.RedirectFixedPath) {
				return
			}
		}
		break
	}

	if e.HandleMethodNotAllowed {
		for _, tree := range e.tree {
			if tree.method == httpMethod {
				continue
			}
			if value := tree.root.getValues(rPath, nil, c.skippedNodes, unescape); value.handlers != nil {
				c.Handlers = e.allNoMethod
				serveError(c, http.StatusMethodNotAllowed, default405Body)
				return
			}
		}
	}
	c.Handlers = e.allNoRoute
	serveError(c, http.StatusNotFound, default404Body)
}

func redirectTrailingSlash(c *Context) {
	req := c.request
	p := req.URL.Path
	if prefix := path.Clean(c.request.Header.Get("X-Forwarded-Prefix")); prefix != "." {
		p = prefix + "/" + req.URL.Path
	}
	req.URL.Path = p + "/"
	if length := len(p); length > 1 && p[length-1] == '/' {
		req.URL.Path = p[:length-1]
	}
	redirectRequest(c)
}

func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
	req := c.request
	rPath := req.URL.Path
	if fixedPath, ok := root.findCaseInsensitivePath(cleanPath(rPath), trailingSlash); ok {
		req.URL.Path = bytesconv.BytesToString(fixedPath)
		redirectRequest(c)
		return true
	}
	return false
}

func serveError(c *Context, code int, defaultMsg []byte) {
	// TODO
}

func redirectRequest(c *Context) {
	//req := c.request
	//rPath := req.URL.Path
	//rURL := req.URL.String()
	//
	//code := http.StatusMovedPermanently
	//if req.Method != http.MethodGet {
	//	code = http.StatusTemporaryRedirect
	//}
	////debugPrint("redirecting request %d: %s --> %s", code, rPath, rURL)
	//http.Redirect(c.Writer, req, rURL, code)
	//c.writermem.WriteHeaderNow()
}
