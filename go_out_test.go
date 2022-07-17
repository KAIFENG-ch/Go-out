package Go_out

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"sync"
	"testing"
)

func TestAddRoute(t *testing.T) {
	router := New()
	router.addRoute("GET", "/", HandlerChain{func(_ *Context) {}})

	assert.Len(t, router.tree, 1)
	assert.Panics(t, func() { router.addRoute("GET", "a", HandlerChain{func(_ *Context) {}}) })
	assert.Panics(t, func() { router.addRoute("GET", "/", HandlerChain{}) })

	router.addRoute("POST", "/post", HandlerChain{func(_ *Context) {}})
	assert.Panics(t, func() {
		router.addRoute("POST", "/post", HandlerChain{func(_ *Context) {}})
	})
}

func TestNew(t *testing.T) {
	router := New()
	assert.Len(t, router.Handlers, 2)
}

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		want *Engine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Default(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Default() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEngine_AllocateContext(t *testing.T) {
	type fields struct {
		RouterGroup            RouterGroup
		pool                   sync.Pool
		maxParams              uint16
		tree                   MethodTrees
		allNoMethod            HandlerChain
		allNoRoute             HandlerChain
		UseRawPath             bool
		UnescapePathValue      bool
		RemoveExtraSlash       bool
		RedirectTrailingSlash  bool
		RedirectFixedPath      bool
		HandleMethodNotAllowed bool
		UseH2C                 bool
	}
	tests := []struct {
		name   string
		fields fields
		want   *Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				RouterGroup:            tt.fields.RouterGroup,
				pool:                   tt.fields.pool,
				maxParams:              tt.fields.maxParams,
				tree:                   tt.fields.tree,
				allNoMethod:            tt.fields.allNoMethod,
				allNoRoute:             tt.fields.allNoRoute,
				UseRawPath:             tt.fields.UseRawPath,
				UnescapePathValue:      tt.fields.UnescapePathValue,
				RemoveExtraSlash:       tt.fields.RemoveExtraSlash,
				RedirectTrailingSlash:  tt.fields.RedirectTrailingSlash,
				RedirectFixedPath:      tt.fields.RedirectFixedPath,
				HandleMethodNotAllowed: tt.fields.HandleMethodNotAllowed,
				UseH2C:                 tt.fields.UseH2C,
			}
			if got := e.AllocateContext(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllocateContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEngine_Handler(t *testing.T) {
	type fields struct {
		RouterGroup            RouterGroup
		pool                   sync.Pool
		maxParams              uint16
		tree                   MethodTrees
		allNoMethod            HandlerChain
		allNoRoute             HandlerChain
		UseRawPath             bool
		UnescapePathValue      bool
		RemoveExtraSlash       bool
		RedirectTrailingSlash  bool
		RedirectFixedPath      bool
		HandleMethodNotAllowed bool
		UseH2C                 bool
	}
	tests := []struct {
		name   string
		fields fields
		want   http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				RouterGroup:            tt.fields.RouterGroup,
				pool:                   tt.fields.pool,
				maxParams:              tt.fields.maxParams,
				tree:                   tt.fields.tree,
				allNoMethod:            tt.fields.allNoMethod,
				allNoRoute:             tt.fields.allNoRoute,
				UseRawPath:             tt.fields.UseRawPath,
				UnescapePathValue:      tt.fields.UnescapePathValue,
				RemoveExtraSlash:       tt.fields.RemoveExtraSlash,
				RedirectTrailingSlash:  tt.fields.RedirectTrailingSlash,
				RedirectFixedPath:      tt.fields.RedirectFixedPath,
				HandleMethodNotAllowed: tt.fields.HandleMethodNotAllowed,
				UseH2C:                 tt.fields.UseH2C,
			}
			if got := e.Handler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEngine_Run(t *testing.T) {
	type fields struct {
		RouterGroup            RouterGroup
		pool                   sync.Pool
		maxParams              uint16
		tree                   MethodTrees
		allNoMethod            HandlerChain
		allNoRoute             HandlerChain
		UseRawPath             bool
		UnescapePathValue      bool
		RemoveExtraSlash       bool
		RedirectTrailingSlash  bool
		RedirectFixedPath      bool
		HandleMethodNotAllowed bool
		UseH2C                 bool
	}
	type args struct {
		addr []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				RouterGroup:            tt.fields.RouterGroup,
				pool:                   tt.fields.pool,
				maxParams:              tt.fields.maxParams,
				tree:                   tt.fields.tree,
				allNoMethod:            tt.fields.allNoMethod,
				allNoRoute:             tt.fields.allNoRoute,
				UseRawPath:             tt.fields.UseRawPath,
				UnescapePathValue:      tt.fields.UnescapePathValue,
				RemoveExtraSlash:       tt.fields.RemoveExtraSlash,
				RedirectTrailingSlash:  tt.fields.RedirectTrailingSlash,
				RedirectFixedPath:      tt.fields.RedirectFixedPath,
				HandleMethodNotAllowed: tt.fields.HandleMethodNotAllowed,
				UseH2C:                 tt.fields.UseH2C,
			}
			if err := e.Run(tt.args.addr...); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEngine_ServeHTTP(t *testing.T) {
	type fields struct {
		RouterGroup            RouterGroup
		pool                   sync.Pool
		maxParams              uint16
		tree                   MethodTrees
		allNoMethod            HandlerChain
		allNoRoute             HandlerChain
		UseRawPath             bool
		UnescapePathValue      bool
		RemoveExtraSlash       bool
		RedirectTrailingSlash  bool
		RedirectFixedPath      bool
		HandleMethodNotAllowed bool
		UseH2C                 bool
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Engine{
				RouterGroup:            tt.fields.RouterGroup,
				pool:                   tt.fields.pool,
				maxParams:              tt.fields.maxParams,
				tree:                   tt.fields.tree,
				allNoMethod:            tt.fields.allNoMethod,
				allNoRoute:             tt.fields.allNoRoute,
				UseRawPath:             tt.fields.UseRawPath,
				UnescapePathValue:      tt.fields.UnescapePathValue,
				RemoveExtraSlash:       tt.fields.RemoveExtraSlash,
				RedirectTrailingSlash:  tt.fields.RedirectTrailingSlash,
				RedirectFixedPath:      tt.fields.RedirectFixedPath,
				HandleMethodNotAllowed: tt.fields.HandleMethodNotAllowed,
				UseH2C:                 tt.fields.UseH2C,
			}
		})
	}
}

func TestEngine_Use(t *testing.T) {
	type fields struct {
		RouterGroup            RouterGroup
		pool                   sync.Pool
		maxParams              uint16
		tree                   MethodTrees
		allNoMethod            HandlerChain
		allNoRoute             HandlerChain
		UseRawPath             bool
		UnescapePathValue      bool
		RemoveExtraSlash       bool
		RedirectTrailingSlash  bool
		RedirectFixedPath      bool
		HandleMethodNotAllowed bool
		UseH2C                 bool
	}
	type args struct {
		middleware []HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   IRoutes
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				RouterGroup:            tt.fields.RouterGroup,
				pool:                   tt.fields.pool,
				maxParams:              tt.fields.maxParams,
				tree:                   tt.fields.tree,
				allNoMethod:            tt.fields.allNoMethod,
				allNoRoute:             tt.fields.allNoRoute,
				UseRawPath:             tt.fields.UseRawPath,
				UnescapePathValue:      tt.fields.UnescapePathValue,
				RemoveExtraSlash:       tt.fields.RemoveExtraSlash,
				RedirectTrailingSlash:  tt.fields.RedirectTrailingSlash,
				RedirectFixedPath:      tt.fields.RedirectFixedPath,
				HandleMethodNotAllowed: tt.fields.HandleMethodNotAllowed,
				UseH2C:                 tt.fields.UseH2C,
			}
			if got := e.Use(tt.args.middleware...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Use() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEngine_addRoute(t *testing.T) {
	type fields struct {
		RouterGroup            RouterGroup
		pool                   sync.Pool
		maxParams              uint16
		tree                   MethodTrees
		allNoMethod            HandlerChain
		allNoRoute             HandlerChain
		UseRawPath             bool
		UnescapePathValue      bool
		RemoveExtraSlash       bool
		RedirectTrailingSlash  bool
		RedirectFixedPath      bool
		HandleMethodNotAllowed bool
		UseH2C                 bool
	}
	type args struct {
		method  string
		path    string
		handler HandlerChain
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Engine{
				RouterGroup:            tt.fields.RouterGroup,
				pool:                   tt.fields.pool,
				maxParams:              tt.fields.maxParams,
				tree:                   tt.fields.tree,
				allNoMethod:            tt.fields.allNoMethod,
				allNoRoute:             tt.fields.allNoRoute,
				UseRawPath:             tt.fields.UseRawPath,
				UnescapePathValue:      tt.fields.UnescapePathValue,
				RemoveExtraSlash:       tt.fields.RemoveExtraSlash,
				RedirectTrailingSlash:  tt.fields.RedirectTrailingSlash,
				RedirectFixedPath:      tt.fields.RedirectFixedPath,
				HandleMethodNotAllowed: tt.fields.HandleMethodNotAllowed,
				UseH2C:                 tt.fields.UseH2C,
			}
		})
	}
}

func TestEngine_handleHttpRequest(t *testing.T) {
	type fields struct {
		RouterGroup            RouterGroup
		pool                   sync.Pool
		maxParams              uint16
		tree                   MethodTrees
		allNoMethod            HandlerChain
		allNoRoute             HandlerChain
		UseRawPath             bool
		UnescapePathValue      bool
		RemoveExtraSlash       bool
		RedirectTrailingSlash  bool
		RedirectFixedPath      bool
		HandleMethodNotAllowed bool
		UseH2C                 bool
	}
	type args struct {
		c *Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Engine{
				RouterGroup:            tt.fields.RouterGroup,
				pool:                   tt.fields.pool,
				maxParams:              tt.fields.maxParams,
				tree:                   tt.fields.tree,
				allNoMethod:            tt.fields.allNoMethod,
				allNoRoute:             tt.fields.allNoRoute,
				UseRawPath:             tt.fields.UseRawPath,
				UnescapePathValue:      tt.fields.UnescapePathValue,
				RemoveExtraSlash:       tt.fields.RemoveExtraSlash,
				RedirectTrailingSlash:  tt.fields.RedirectTrailingSlash,
				RedirectFixedPath:      tt.fields.RedirectFixedPath,
				HandleMethodNotAllowed: tt.fields.HandleMethodNotAllowed,
				UseH2C:                 tt.fields.UseH2C,
			}
		})
	}
}

func TestHandlerChain_Last(t *testing.T) {
	tests := []struct {
		name string
		h    HandlerChain
		want HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Last(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew1(t *testing.T) {
	tests := []struct {
		name string
		want *Engine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redirectFixedPath(t *testing.T) {
	type args struct {
		c             *Context
		root          *node
		trailingSlash bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redirectFixedPath(tt.args.c, tt.args.root, tt.args.trailingSlash); got != tt.want {
				t.Errorf("redirectFixedPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redirectRequest(t *testing.T) {
	type args struct {
		c *Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_redirectTrailingSlash(t *testing.T) {
	type args struct {
		c *Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_serveError(t *testing.T) {
	type args struct {
		c          *Context
		code       int
		defaultMsg []byte
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}