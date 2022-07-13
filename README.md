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
