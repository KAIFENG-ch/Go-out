package Go_out

import (
	"net/http"
	"reflect"
	"testing"
)

func TestContext_Abort(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Context{
				writermem:    tt.fields.writermem,
				request:      tt.fields.request,
				Writer:       tt.fields.Writer,
				Handlers:     tt.fields.Handlers,
				engine:       tt.fields.engine,
				params:       tt.fields.params,
				Params:       tt.fields.Params,
				Keys:         tt.fields.Keys,
				skippedNodes: tt.fields.skippedNodes,
				FullPath:     tt.fields.FullPath,
				Index:        tt.fields.Index,
				Errors:       tt.fields.Errors,
			}
		})
	}
}

func TestContext_Handler(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Context{
				writermem:    tt.fields.writermem,
				request:      tt.fields.request,
				Writer:       tt.fields.Writer,
				Handlers:     tt.fields.Handlers,
				engine:       tt.fields.engine,
				params:       tt.fields.params,
				Params:       tt.fields.Params,
				Keys:         tt.fields.Keys,
				skippedNodes: tt.fields.skippedNodes,
				FullPath:     tt.fields.FullPath,
				Index:        tt.fields.Index,
				Errors:       tt.fields.Errors,
			}
			if got := c.Handler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_IsAborted(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Context{
				writermem:    tt.fields.writermem,
				request:      tt.fields.request,
				Writer:       tt.fields.Writer,
				Handlers:     tt.fields.Handlers,
				engine:       tt.fields.engine,
				params:       tt.fields.params,
				Params:       tt.fields.Params,
				Keys:         tt.fields.Keys,
				skippedNodes: tt.fields.skippedNodes,
				FullPath:     tt.fields.FullPath,
				Index:        tt.fields.Index,
				Errors:       tt.fields.Errors,
			}
			if got := c.IsAborted(); got != tt.want {
				t.Errorf("IsAborted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_Next(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Context{
				writermem:    tt.fields.writermem,
				request:      tt.fields.request,
				Writer:       tt.fields.Writer,
				Handlers:     tt.fields.Handlers,
				engine:       tt.fields.engine,
				params:       tt.fields.params,
				Params:       tt.fields.Params,
				Keys:         tt.fields.Keys,
				skippedNodes: tt.fields.skippedNodes,
				FullPath:     tt.fields.FullPath,
				Index:        tt.fields.Index,
				Errors:       tt.fields.Errors,
			}
		})
	}
}

func TestContext_Reset(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Context{
				writermem:    tt.fields.writermem,
				request:      tt.fields.request,
				Writer:       tt.fields.Writer,
				Handlers:     tt.fields.Handlers,
				engine:       tt.fields.engine,
				params:       tt.fields.params,
				Params:       tt.fields.Params,
				Keys:         tt.fields.Keys,
				skippedNodes: tt.fields.skippedNodes,
				FullPath:     tt.fields.FullPath,
				Index:        tt.fields.Index,
				Errors:       tt.fields.Errors,
			}
		})
	}
}

func Test_responseWriter_Written(t *testing.T) {
	type fields struct {
		ResponseWriter http.ResponseWriter
		size           int
		status         int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &responseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				size:           tt.fields.size,
				status:         tt.fields.status,
			}
			if got := r.Written(); got != tt.want {
				t.Errorf("Written() = %v, want %v", got, tt.want)
			}
		})
	}
}
