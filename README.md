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

