# Go-out
我们可以从最简单的实现开始看

```
package main

import (
  gout "Go-out"
)

func hello() {
  r := gout.Default()
  r.GET("/hello", func(c *gin.Context) {
    c.JSON(200, "pong")
  })
  _ = r.Run("8080")
}
```

从这段代码中可以看出我们gout框架的引擎创建入口是Default,然后我们使用run来启动这个引擎
这样我们看到Default这个函数，里面就是调用到了New这个方法并封装了Logger和Recover两个中间件，这样我们在聚焦New函数，可以看到我们New函数里面实例化了一个Engine引擎
```
func New() *Engine {
	debugPrintWARNINGNew()
	engine := &Engine{ 
		RouterGroup: RouterGroup{ 
			Handlers: nil,
			basePath: "/",
			root:     true, 
		},
	}
    engine.RouterGroup.engine = engine
    engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}

```
在这里我们先创建了一个引擎对象，这个对象包括这个路由对应的方法，并将其路由设置为根路由，然后我们在路由组中添加这个引擎并在线程池中为这个路由对应的方法分配了内存空间，这里都很好理解。
这里我们还涉及到一个路由组的实现，这个路由器的实现有一个继承的概念，为什么会用到继承呢，因为这里针对所有路由我们创建了一个HandlerChain
```
type HandlerChain []HandlerFunc

```

但我们的路由是一层一层实现的，所以我们用到了combineHandlers方法来连接路由
```
func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	assert1(finalSize < int(abortIndex), "too many handlers")
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

```
这里的finalSize是将旧的handler和新的handler连接，然后对这个handler数量进行一个判断，如果超过64则会推出，然后用copy方法创建新的切片。
```
func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}
```
之后我们会将绝对路径转换为相对路径，这样就可以在原有handler的基础上调用下一个handler

既然New了一个engine对象，然后我们就要用这个engine来处理请求，这里就涉及到了我们的Context对象，Context在中文中的含义为上下文,在Go的标准库中，Context被用作进程上下文之间的切换,
而在我们的gout中，context可以被理解为一个请求的接受与处理中心，它通过接收用户的请求来完成上下文handler之间的切换
我们可以来看看context的实现
```
type Context struct {
	writermem responseWriter
	Request   *http.Request	 // 请求对象
	Writer    ResponseWriter // 响应对象
	Params   Params 		 // 路由参数 /user/:id 这个id
	handlers HandlersChain	 // 中间件数组 
	index    int8 		// 当前执行中间件的下标
	fullPath string  	// 请求的完整路径

	engine       *Engine
	params       *Params
	skippedNodes *[]skippedNode 

	mu sync.RWMutex 	// 保证Keys map的线程安全
	Keys map[string]any // 对每一个请求进行处理存储

	Errors errorMsgs 	// 存储错误的列表
	Accepted []string
	queryCache url.Values // 存放url请求参数
	formCache url.Values  // 存放form参数
	sameSite http.SameSite
}
```
这里可以看到，context中
