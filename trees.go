package Go_out

type methodTrees struct {
	method string
	root   *node
}

type MethodTree []methodTrees

type skippedNode struct {
	path        string
	node        *node
	paramsCount int16
}

type node struct {
	path      string
	indices   string
	wildChild bool
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
