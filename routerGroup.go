package Go_out

import (
	"net/http"
	"regexp"
)

var (
	regEnLetter = regexp.MustCompile("^[A-Z]+$")
	//methods     = []string{
	//	http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut,
	//}
)

type (
	RouterGroup struct {
		Handlers HandlerChain
		bashPath string
		engine   *Engine
		root     bool
	}
	IRouter interface {
		Group(string ...HandlerFunc) *RouterGroup
	}
	IRoutes interface {
		Use(...HandlerFunc) IRoutes
		Handle(string, string, ...HandlerFunc) IRoutes
		GET(string, ...HandlerFunc) IRoutes
		POST(string, ...HandlerFunc) IRoutes
		DELETE(string, ...HandlerFunc) IRoutes
		PUT(string, ...HandlerFunc) IRoutes
	}
)

func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}

func (group *RouterGroup) Group(relativePath string, handler ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandler(handler),
		bashPath: joinPaths(group.bashPath, relativePath),
		engine:   group.engine,
	}
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.engine
	}
	return group
}

func (group *RouterGroup) Handle(httpMethod string, relativePath string, handlers ...HandlerFunc) IRoutes {
	if matched := regEnLetter.MatchString(httpMethod); !matched {
		panic("http method" + httpMethod + "is not valid")
	}
	return group.Handle(httpMethod, relativePath, handlers...)
}

func (group *RouterGroup) GET(relativePath string, handlerFunc ...HandlerFunc) IRoutes {
	return group.Handle(http.MethodGet, relativePath, handlerFunc...)
}

func (group *RouterGroup) POST(relativePath string, handlerFunc ...HandlerFunc) IRoutes {
	return group.Handle(http.MethodPost, relativePath, handlerFunc...)
}

func (group *RouterGroup) DELETE(relativePath string, handlerFunc ...HandlerFunc) IRoutes {
	return group.Handle(http.MethodDelete, relativePath, handlerFunc...)
}

func (group *RouterGroup) PUT(relativePath string, handlerFunc ...HandlerFunc) IRoutes {
	return group.Handle(http.MethodPut, relativePath, handlerFunc...)
}

func (group *RouterGroup) combineHandler(handlers HandlerChain) HandlerChain {
	finalSize := len(group.Handlers) + len(handlers)
	//if finalSize < int(abortIndex) {
	//	panic("too many arguments")
	//}
	mergedHandlers := make(HandlerChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.bashPath, relativePath)
}
