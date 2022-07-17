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
	Request   *http.Request	
	Writer    ResponseWriter 
	Params   Params 		
	handlers HandlersChain	 
	index    int8 		
	fullPath string  	

	engine       *Engine
	params       *Params
	skippedNodes *[]skippedNode 

	Keys map[string]any 

	Errors errorMsgs 	
	Accepted []string
	queryCache url.Values 
	formCache url.Values  
	sameSite http.SameSite
}
```
这里可以看到，context中对http中的request和writer进行了封装，其中还包含着中间件数组和请求的完整路径，
Context相对于http的一大优势在于上下文之间携程的拷贝传递，这里的实现用到了copy函数
```
func (c *Context) Copy() *Context {
	cp := Context{
		writermem: c.writermem,
		Request:   c.Request,
		Params:    c.Params,
		engine:    c.engine,
	}
	cp.writermem.ResponseWriter = nil
	cp.Writer = &cp.writermem
	cp.index = abortIndex
	cp.handlers = nil
	cp.Keys = map[string]any{}
	for k, v := range c.Keys {
		cp.Keys[k] = v
	}
	paramCopy := make([]Param, len(cp.Params))
	copy(paramCopy, cp.Params)
	cp.Params = paramCopy
	return &cp
}
```
在gout中，我们封装了五种请求方式，每一种请求方式通过routerGroup中的handle函数进行处理
```
func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(httpMethod, absolutePath, handlers)
	return group.returnObj()
}
```
这里涉及到gout框架的一个核心：前缀路由树。前缀树是将路由的公共前缀保存在节点中，从根节点开始通过公共节点
进行路由查找可以大大降低查找的时间复杂度
GIN 中常用的声明路由的方式如下
```
// trees 路由树这一部分由一个带有method 和root字段的node列表维护
// 每个node代表了路由树中的每一个节点
// node所具有的字段内容如下

type node struct {
    path      string // 当前节点的绝对路径
    indices   string // 缓存下一节点的第一个字符 在遇到子节点为通配符类型的情况下,indices=''
        // 默认是 false，当 children 是 通配符类型时，wildChild 为 true 即 indices=''
    wildChild bool // 默认是 false，当 children 是 通配符类型时，wildChild 为 true

        // 节点的类型，因为在通配符的场景下在查询的时候需要特殊处理， 
        // 默认是static类型
        // 根节点为 root类型
        // 对于 path 包含冒号通配符的情况，nType 是 param 类型
        // 对于包含 * 通配符的情况，nType 类型是 catchAll 类型
    nType     nodeType
        // 代表了有几条路由会经过此节点，用于在节点
    priority  uint32
        // 子节点列表
    children  []*node // child nodes, at most 1 :param style node at the end of the array
    handlers  HandlersChain
        // 是从 root 节点到当前节点的全部 path 部分；如果此节点为终结节点 handlers 为对应的处理链，否则为 nil；
        // maxParams 是当前节点到各个叶子节点的包含的通配符的最大数量
    fullPath  string
}

// 具体节点类型如下
const (
    static nodeType = iota // default， 静态节点，普通匹配(/user)
    root                   // 根节点 (/)
    param                 // 参数节点(/user/:id)
    catchAll              // 通用匹配，匹配任意参数(*user)
)
```
路由的添加主要通过addRoute函数完成
```
func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
   // 校验
   // 路径必须以 / 开头
   // 请求方法不允许为空
   // 处理方法不允许为空
    if path[0] != '/' {
		panic("path must begin with '/'")
	}
	if method == "" {
		panic("http method can not be empty")
	}
	if len(handler) <= 0 {
		panic("must be at least one handlers")
	}
   // 如果开启了gin的debug模式，则对应处理
   debugPrintRoute(method, path, handlers)
   // 根据请求方式获取对应的树的根
   // 每一个请求方法都有自己对应的一颗紧凑前缀树，这里通过请求方法拿到最顶部的根
   root := engine.trees.get(method)
   // 如果根为空，则表示这是第一个路由，则自己创建一个以 / 为path的根节点
   if root == nil {
      // 如果没有就创建
      root = new(node)
      root.fullPath = "/"
      engine.trees = append(engine.trees, methodTree{method: method, root: root})
   }
   // 此处的path是子路由
   // 以上内容是做了一层预校验，避免书写不规范导致的请求查询不到
   // 接下来是添加路由的正文
   root.addRoute(path, handlers)
}
```

```
// addRoute adds a node with the given handle to the path.
// Not concurrency-safe! 并发不安全
func (n *node) addRoute(path string, handlers HandlersChain) {
    fullPath := path
        // 添加完成后，经过此节点的路由条数将会+1
    n.priority++

    // Empty tree
        // 如果为空树， 即只有一个根节点"/" 则插入一个子节点， 并将当前节点设置为root类型的节点
    if len(n.path) == 0 && len(n.children) == 0 {
        n.insertChild(path, fullPath, handlers)
        n.nType = root
        return
    }

    parentFullPathIndex := 0

walk:
    for {
        // Find the longest common prefix.
        // This also implies that the common prefix contains no ':' or '*'
        // since the existing key can't contain those chars.
                // 找到最长的共有前缀的长度 即到i位置 path[i] == n.path[i]
        i := longestCommonPrefix(path, n.path)

        // Split edge
                // 假设当前节点存在的前缀信息为 hello
                // 现有前缀信息为heo的结点进入， 则当前节点需要被拆分
                // 拆分成为 he节点 以及 (llo 和 o 两个子节点)
        if i < len(n.path) {
            child := node{
                                // 除去公共前缀部分，剩余的内容作为子节点
                path:      n.path[i:],
                wildChild: n.wildChild,
                indices:   n.indices,
                children:  n.children,
                handlers:  n.handlers,
                priority:  n.priority - 1,
                fullPath:  n.fullPath,
            }

            n.children = []*node{&child}
            // []byte for proper unicode char conversion, see #65
            n.indices = bytesconv.BytesToString([]byte{n.path[i]})
            n.path = path[:i]
            n.handlers = nil
            n.wildChild = false
            n.fullPath = fullPath[:parentFullPathIndex+i]
        }

        // Make new node a child of this node
                // 将新来的节点插入新的parent节点作为子节点
        if i < len(path) {
            path = path[i:]
            c := path[0]

            // '/' after param
                        // 如果是参数节点 形如/:i
            if n.nType == param && c == '/' && len(n.children) == 1 {
                parentFullPathIndex += len(n.path)
                n = n.children[0]
                n.priority++
                continue walk
            }

            // Check if a child with the next path byte exists
            for i, max := 0, len(n.indices); i < max; i++ {
                if c == n.indices[i] {
                    parentFullPathIndex += len(n.path)
                    i = n.incrementChildPrio(i)
                    n = n.children[i]
                    continue walk
                }
            }

            // Otherwise insert it
            if c != ':' && c != '*' && n.nType != catchAll {
                // []byte for proper unicode char conversion, see #65
                n.indices += bytesconv.BytesToString([]byte{c})
                child := &node{
                    fullPath: fullPath,
                }
                n.addChild(child)
                n.incrementChildPrio(len(n.indices) - 1)
                n = child
            } else if n.wildChild {
                // inserting a wildcard node, need to check if it conflicts with the existing wildcard
                n = n.children[len(n.children)-1]
                n.priority++

                // Check if the wildcard matches
                if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
                    // Adding a child to a catchAll is not possible
                    n.nType != catchAll &&
                    // Check for longer wildcard, e.g. :name and :names
                    (len(n.path) >= len(path) || path[len(n.path)] == '/') {
                    continue walk
                }

                // Wildcard conflict
                pathSeg := path
                if n.nType != catchAll {
                    pathSeg = strings.SplitN(pathSeg, "/", 2)[0]
                }
                prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
                panic("'" + pathSeg +
                    "' in new path '" + fullPath +
                    "' conflicts with existing wildcard '" + n.path +
                    "' in existing prefix '" + prefix +
                    "'")
            }

            n.insertChild(path, fullPath, handlers)
            return
        }

        // Otherwise add handle to current node
                // 设置处理函数，如果已经存在，则报错
        if n.handlers != nil {
            panic("handlers are already registered for path '" + fullPath + "'")
        }
        n.handlers = handlers
        n.fullPath = fullPath
        return
    }
}
```
node中的priority是优先级的意思，在查找时会根据priority对节点进行排序，常用节点在最前，并且节点中的priority
越大，越优先进行分配

路由树构建完毕后便开始查找路由，第一步是使用HTTPServe解析路由地址，步骤是：
    1.申请一块内存来填充响应体
    2.处理请求信息
    3.拿到请求对应的路由树
    4.获取根节点
```
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    c := engine.pool.Get().(*Context)
    c.writermem.reset(w)
    c.Request = req
    c.reset()

    // 真正开始处理请求
    engine.handleHTTPRequest(c)

    engine.pool.Put(c)
}
```
```
func (engine *Engine) handleHTTPRequest(c *Context) {
    // ...
    t := engine.trees
    for i, tl := 0, len(t); i < tl; i++ {
        // 根据请求方法进行判断
        if t[i].method != httpMethod {
            continue
        }
        root := t[i].root
        // 在该方法树上查找路由
        value := root.getValue(rPath, c.params, unescape)
        if value.params != nil {
            c.Params = *value.params
        }
        // 执行处理函数
        if value.handlers != nil {
            c.handlers = value.handlers
            c.fullPath = value.fullPath
            c.Next() // 涉及到gin的中间件机制
            // 到这里时，请求已经处理完毕，返回的结果也存储在对应的结构体中了
            c.writermem.WriteHeaderNow()
            return
        }
        // ...
      break
   }
   if engine.HandleMethodNotAllowed {
        for _, tree := range engine.trees {
            if tree.method == httpMethod {
                continue
            }
            if value := tree.root.getValue(rPath, nil, c.skippedNodes, unescape); value.handlers != nil {
                c.handlers = engine.allNoMethod
                serveError(c, http.StatusMethodNotAllowed, default405Body)
                return
            }
        }
   }
}
```
最后我们看Run方法，Run实际上是封装了http的监听函数，来达到对目标端口的监听作用
```
func (e *Engine) Run(addr ...string) (err error) {
	address := resolveAddr(addr)
	err = http.ListenAndServe(address, e.Handler())
	return
}
```
其中的ListenAndServe就是对端口的监听
```
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```
```
func (srv *Server) ListenAndServe() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}
```
```
func (srv *Server) Serve(l net.Listener) error {
	if fn := testHookServerServe; fn != nil {
		fn(srv, l) // call hook with unwrapped listener
	}

	origListener := l
	l = &onceCloseListener{Listener: l}
	defer l.Close()

	if err := srv.setupHTTP2_Serve(); err != nil {
		return err
	}

	if !srv.trackListener(&l, true) {
		return ErrServerClosed
	}
	defer srv.trackListener(&l, false)

	baseCtx := context.Background()
	if srv.BaseContext != nil {
		baseCtx = srv.BaseContext(origListener)
		if baseCtx == nil {
			panic("BaseContext returned a nil context")
		}
	}

	var tempDelay time.Duration // how long to sleep on accept failure

	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.logf("http: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		connCtx := ctx
		if cc := srv.ConnContext; cc != nil {
			connCtx = cc(connCtx, rw)
			if connCtx == nil {
				panic("ConnContext returned nil")
			}
		}
		tempDelay = 0
		c := srv.newConn(rw)
		c.setState(c.rwc, StateNew, runHooks) // before Serve can return
		go c.serve(connCtx)
	}
}
```
```
func (c *conn) serve(ctx context.Context) {
    ... ...
    if err := tlsConn.HandShake(); err != nil {
        return
    }	
    ... ...
    for {
        w, err := c.ReadRequest(ctx)
        ... ... 
        serveHandler{c.server}.serveHTTP(w, w.req)
        ... ...
        if !w.conn.serve.keepAlives() {
            ...
            return  
        }
        ...
    }
```