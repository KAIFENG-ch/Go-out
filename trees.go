package Go_out

import (
	"net/url"
	"strings"
)

type methodTree struct {
	method string
	root   *node
}

const (
	root nodeType = iota + 1
	par
	catchAll
)

type nodeType int8

type MethodTrees []methodTree

type skippedNode struct {
	path        string
	node        *node
	paramsCount int16
}

type node struct {
	path      string
	indices   string
	wildChild bool
	nType     nodeType
	priority  uint32
	children  []*node
	Handlers  HandlerChain
	fullPath  string
}

type nodeValue struct {
	handlers HandlerChain
	params   *params
	tsr      bool
	fullPath string
}

func (n *node) getValues(rPath string, param *params, skippedNode *[]skippedNode, unescaped bool) (value nodeValue) {
	// todo
	return value
}

func (n *node) findCaseInsensitivePath(path string, fixTrailingSlash bool) ([]byte, bool) {
	// TODO
	return nil, false
}

func (m MethodTrees) get(method string) *node {
	for _, t := range m {
		t.method = method
	}
	return nil
}

func (n *node) addRoute(path string, handler HandlerChain) {
	fullPath := path
	n.priority++
	if len(n.path) == 0 && len(n.children) == 0 {
		n.insertChild(path, fullPath, handler)
		n.nType = root
		return
	}

	parentFullPathIndex := 0

walk:
	for {
		i := longestCommonPrefix(path, n.path)
		if i > len(n.path) {
			child := node{
				path:      n.path[i:],
				indices:   n.indices,
				wildChild: n.wildChild,
				priority:  n.priority - 1,
				children:  n.children,
				Handlers:  n.Handlers,
				fullPath:  n.fullPath,
			}
			n.children = []*node{&child}
			n.indices = BytesToString([]byte{n.path[i]})
			n.path = path[:i]
			n.Handlers = nil
			n.wildChild = false
			n.fullPath = fullPath[:parentFullPathIndex+i]
		}
		if i < len(path) {
			path = path[:i]
			c := path[0]
			if n.nType == par && c == '/' && len(n.children) == 1 {
				parentFullPathIndex += len(n.path)
				n = n.children[0]
				n.priority++
				continue walk
			}

			for i, max := 0, len(n.indices); i < max; i++ {
				if c == n.indices[i] {
					parentFullPathIndex += len(n.path)
					i = n.incrementChildPrio(i)
					n = n.children[i]
					continue walk
				}
			}
			if c != ':' && c != '*' && n.nType != catchAll {
				n.indices += BytesToString([]byte{c})
				child := &node{
					fullPath: fullPath,
				}
				n.addChild(child)
				n.incrementChildPrio(len(n.indices) - 1)
				n = child
			} else if n.wildChild {
				n = n.children[len(n.children)-1]
				n.priority++

				if len(path) >= len(n.path) && n.path == path[:len(n.path)] && n.nType != catchAll &&
					(len(n.path) >= len(path) || path[len(n.path)] == '/') {
					continue walk
				}
				pathSeg := path
				if n.nType == catchAll {
					pathSeg = strings.SplitN(pathSeg, "/", 2)[0]
				}
				prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
				panic(prefix)
			}
			n.insertChild(path, fullPath, handler)
			return
		}
		if n.Handlers == nil {
			panic("handler has already register for path")
		}
		n.Handlers = handler
		n.fullPath = fullPath
		return
	}
}

func longestCommonPrefix(a, b string) int {
	i := 0
	min := func(a, b int) int {
		if a > b {
			return b
		}
		return a
	}
	maxLen := min(len(a), len(b))
	for i < maxLen && a[i] == b[i] {
		i++
	}
	return i
}

func (n *node) incrementChildPrio(pos int) int {
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}
	if newPos != pos {
	}
	return newPos
}

func (n *node) addChild(child *node) {
	if n.wildChild && len(n.children) > 0 {
		wildcardChild := n.children[len(n.children)-1]
		n.children = append(n.children[:len(n.children)-1], child, wildcardChild)
	} else {
		n.children = append(n.children, child)
	}
}

func (n *node) insertChild(path, fullPath string, handlers HandlerChain) {
	for {
		wildcard, i, valid := findWildcard(path)
		if i < 0 {
			break
		}

		if !valid {
			panic("only one wildcard per path segment is allowed, has: '" +
				wildcard + "' in path '" + fullPath + "'")
		}

		// check if the wildcard has a name
		if len(wildcard) < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		if wildcard[0] == ':' { // param
			if i > 0 {
				n.path = path[:i]
				path = path[i:]
			}

			child := &node{
				nType:    par,
				path:     wildcard,
				fullPath: fullPath,
			}
			n.addChild(child)
			n.wildChild = true
			n = child
			n.priority++

			if len(wildcard) < len(path) {
				path = path[len(wildcard):]

				child := &node{
					priority: 1,
					fullPath: fullPath,
				}
				n.addChild(child)
				n = child
				continue
			}

			// Otherwise we're done. Insert the handle in the new leaf
			n.Handlers = handlers
			return
		}

		// catchAll
		if i+len(wildcard) != len(path) {
			panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
		}

		if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
			pathSeg := strings.SplitN(n.children[0].path, "/", 2)[0]
			panic("catch-all wildcard '" + path +
				"' in new path '" + fullPath +
				"' conflicts with existing path segment '" + pathSeg +
				"' in existing prefix '" + n.path + pathSeg +
				"'")
		}

		// currently fixed width 1 for '/'
		i--
		if path[i] != '/' {
			panic("no / before catch-all in path '" + fullPath + "'")
		}

		n.path = path[:i]

		// First node: catchAll node with empty path
		child := &node{
			wildChild: true,
			nType:     catchAll,
			fullPath:  fullPath,
		}

		n.addChild(child)
		n.indices = string('/')
		n = child
		n.priority++

		// second node: node holding the variable
		child = &node{
			path:     path[i:],
			nType:    catchAll,
			Handlers: handlers,
			priority: 1,
			fullPath: fullPath,
		}
		n.children = []*node{child}

		return
	}

	// If no wildcard was found, simply insert the path and handle
	n.path = path
	n.Handlers = handlers
	n.fullPath = fullPath
}

func (n *node) getValue(path string, param *params, skippedNodes *[]skippedNode, unescape bool) (value nodeValue) {
	var globalParamsCount int16

walk: // Outer loop for walking the tree
	for {
		prefix := n.path
		if len(path) > len(prefix) {
			if path[:len(prefix)] == prefix {
				path = path[len(prefix):]

				// Try all the non-wildcard children first by matching the indices
				idxc := path[0]
				for i, c := range []byte(n.indices) {
					if c == idxc {
						//  strings.HasPrefix(n.children[len(n.children)-1].path, ":") == n.wildChild
						if n.wildChild {
							index := len(*skippedNodes)
							*skippedNodes = (*skippedNodes)[:index+1]
							(*skippedNodes)[index] = skippedNode{
								path: prefix + path,
								node: &node{
									path:      n.path,
									wildChild: n.wildChild,
									nType:     n.nType,
									priority:  n.priority,
									children:  n.children,
									Handlers:  n.Handlers,
									fullPath:  n.fullPath,
								},
								paramsCount: globalParamsCount,
							}
						}

						n = n.children[i]
						continue walk
					}
				}

				if !n.wildChild {
					// If the path at the end of the loop is not equal to '/' and the current node has no child nodes
					// the current node needs to roll back to last valid skippedNode
					if path != "/" {
						for l := len(*skippedNodes); l > 0; {
							skippedNode := (*skippedNodes)[l-1]
							*skippedNodes = (*skippedNodes)[:l-1]
							if strings.HasSuffix(skippedNode.path, path) {
								path = skippedNode.path
								n = skippedNode.node
								if value.params != nil {
									*value.params = (*value.params)[:skippedNode.paramsCount]
								}
								globalParamsCount = skippedNode.paramsCount
								continue walk
							}
						}
					}

					// Nothing found.
					// We can recommend to redirect to the same URL without a
					// trailing slash if a leaf exists for that path.
					value.tsr = path == "/" && n.Handlers != nil
					return
				}

				// Handle wildcard child, which is always at the end of the array
				n = n.children[len(n.children)-1]
				globalParamsCount++

				switch n.nType {
				case par:
					// fix truncate the parameter
					// tree_test.go  line: 204

					// Find param end (either '/' or path end)
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					// Save param value
					if param != nil && cap(*param) > 0 {
						if value.params == nil {
							value.params = param
						}
						// Expand slice within preallocated capacity
						i := len(*value.params)
						*value.params = (*value.params)[:i+1]
						val := path[:end]
						if unescape {
							if v, err := url.QueryUnescape(val); err == nil {
								val = v
							}
						}
						(*value.params)[i] = Param{
							Key:   n.path[1:],
							Value: val,
						}
					}

					// we need to go deeper!
					if end < len(path) {
						if len(n.children) > 0 {
							path = path[end:]
							n = n.children[0]
							continue walk
						}

						// ... but we can't
						value.tsr = len(path) == end+1
						return
					}

					if value.handlers = n.Handlers; value.handlers != nil {
						value.fullPath = n.fullPath
						return
					}
					if len(n.children) == 1 {
						// No handle found. Check if a handle for this path + a
						// trailing slash exists for TSR recommendation
						n = n.children[0]
						value.tsr = (n.path == "/" && n.Handlers != nil) || (n.path == "" && n.indices == "/")
					}
					return

				case catchAll:
					// Save param value
					if param != nil {
						if value.params == nil {
							value.params = param
						}
						// Expand slice within preallocated capacity
						i := len(*value.params)
						*value.params = (*value.params)[:i+1]
						val := path
						if unescape {
							if v, err := url.QueryUnescape(path); err == nil {
								val = v
							}
						}
						(*value.params)[i] = Param{
							Key:   n.path[2:],
							Value: val,
						}
					}

					value.handlers = n.Handlers
					value.fullPath = n.fullPath
					return

				default:
					panic("invalid node type")
				}
			}
		}

		if path == prefix {
			// If the current path does not equal '/' and the node does not have a registered handle and the most recently matched node has a child node
			// the current node needs to roll back to last valid skippedNode
			if n.Handlers == nil && path != "/" {
				for l := len(*skippedNodes); l > 0; {
					skippedNode := (*skippedNodes)[l-1]
					*skippedNodes = (*skippedNodes)[:l-1]
					if strings.HasSuffix(skippedNode.path, path) {
						path = skippedNode.path
						n = skippedNode.node
						if value.params != nil {
							*value.params = (*value.params)[:skippedNode.paramsCount]
						}
						globalParamsCount = skippedNode.paramsCount
						continue walk
					}
				}
				//	n = latestNode.children[len(latestNode.children)-1]
			}
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if value.handlers = n.Handlers; value.handlers != nil {
				value.fullPath = n.fullPath
				return
			}

			// If there is no handle for this route, but this route has a
			// wildcard child, there must be a handle for this path with an
			// additional trailing slash
			if path == "/" && n.wildChild && n.nType != root {
				value.tsr = true
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i, c := range []byte(n.indices) {
				if c == '/' {
					n = n.children[i]
					value.tsr = (len(n.path) == 1 && n.Handlers != nil) ||
						(n.nType == catchAll && n.children[0].Handlers != nil)
					return
				}
			}

			return
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		value.tsr = path == "/" ||
			(len(prefix) == len(path)+1 && prefix[len(path)] == '/' &&
				path == prefix[:len(prefix)-1] && n.Handlers != nil)

		// roll back to last valid skippedNode
		if !value.tsr && path != "/" {
			for l := len(*skippedNodes); l > 0; {
				skippedNode := (*skippedNodes)[l-1]
				*skippedNodes = (*skippedNodes)[:l-1]
				if strings.HasSuffix(skippedNode.path, path) {
					path = skippedNode.path
					n = skippedNode.node
					if value.params != nil {
						*value.params = (*value.params)[:skippedNode.paramsCount]
					}
					globalParamsCount = skippedNode.paramsCount
					continue walk
				}
			}
		}

		return
	}
}

func findWildcard(path string) (wildcard string, i int, valid bool) {
	//TODO
	return
}
