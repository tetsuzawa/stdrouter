package stdrouter

import (
	"fmt"
	"path"
	"strings"
)

// HandlerFunc is wrapper of http handler function for Node.
type HandlerFunc struct {
	Package string
	Func    string
}

// Node is node of tree structure for path.
type Node struct {
	Depth       int
	Endpoint    string
	IsPathParam bool
	Methods     map[string]HandlerFunc
	Parent      *Node
	Children    []*Node
}

// Add creates new node to node tree.
func (n *Node) Add(p string, httpMethod string, handlerFunc HandlerFunc) error {
	// initialize search queue for BFS.
	queue := make([]*Node, 0)
	queue = append(queue, n)

	childHasDesiredEndPoint := false
	desiredDepth := strings.Count(p, "/")
	head, p := SeparatePath(p, 1)

	// start BFS
	for len(queue) != 0 {
		var node *Node
		node, queue = queue[0], queue[1:]

		// if current node endpoint is same as head of path
		if head[1:] == node.Endpoint {
			// if the path is the end, register the function and end to search
			if strings.Count(p, "/") == 0 {
				if node.Methods == nil {
					node.Methods = make(map[string]HandlerFunc)
				}
				node.Methods[httpMethod] = handlerFunc
				return nil
			}
			// shift paths to search child elements
			head, p = SeparatePath(p, 1)
			// reset flag
			childHasDesiredEndPoint = false
		}

		// check whether head of path is a path parameter
		isPathParam := false
		if len(head) > 1 {
			if head[1] == ':' {
				isPathParam = true
				// "/:param" to "/param"
				head = head[:1] + head[2:]
			}
		}

		// add children nodes to the search queue
		for _, cn := range node.Children {
			queue = append(queue, cn)

			// if the child node has the desired endpoint, set the flag to search next node
			if head[1:] == cn.Endpoint {
				childHasDesiredEndPoint = true
				break
			}
		}

		// if the flag is true, search next node
		if childHasDesiredEndPoint {
			continue
		}

		// if depth of current node is same as the node to add, search next node
		if node.Depth >= desiredDepth {
			continue
		}

		newChild := &Node{
			Depth:       node.Depth + 1,
			Endpoint:    head[1:],
			IsPathParam: isPathParam,
			Parent:      node,
		}
		queue = append(queue, newChild)
		node.Children = append(node.Children, newChild)
	}
	return nil
}


// Print prints general information of each nodes.
func (n *Node) Print() {
	Walk(n, func(node *Node) bool {
		fmt.Printf("Depth: %v, Endpoint: %v, IsPathPram: %v, Methods: %v, hasParent: %v, numChildren: %v\n",
			node.Depth, node.Endpoint, node.IsPathParam, node.Methods, node.Parent == nil, len(node.Children))
		return true
	})
}

// Walk performs a BFS (breadth-first search) on the tree structure of nodes
// and executes the argument function on each node.
// If the return value is false, the search ends.
func Walk(node *Node, fn func(*Node) bool) {
	// initialize search queue for BFS.
	queue := make([]*Node, 0)
	queue = append(queue, node)

	// start BFS
	for len(queue) != 0 {
		// get the first node from queue
		var pn *Node
		pn, queue = queue[0], queue[1:]

		// add children nodes to the search queue
		for _, cn := range pn.Children {
			queue = append(queue, cn)
		}

		// execute function
		if !fn(pn) {
			return
		}
	}
}

// BuildBasePath builds base path from root to argument node.
// If path param node is found, stop building and returns the path.
func BuildBasePath(node *Node) (p string) {
	n := node
	for n.Parent != nil {
		if n.Parent.IsPathParam{
			break
		}
		p = path.Clean("/" + n.Parent.Endpoint + p)
		n = n.Parent
	}
	return p
}
