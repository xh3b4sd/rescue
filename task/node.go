package task

import "github.com/xh3b4sd/rescue/matcher"

type Node map[string]string

func (n *Node) All(key ...string) *Node {
	nod := Node(matcher.All(*n, key...))
	return &nod
}

func (n *Node) Any(key ...string) *Node {
	nod := Node(matcher.Any(*n, key...))
	return &nod
}

func (n *Node) Emp() bool {
	return n.Len() == 0
}

func (n *Node) Eql(x *Node) bool {
	return n != nil && x != nil && n.Len() == x.Len() && n.Has(*x)
}

func (n *Node) Exi(key string) bool {
	if n == nil {
		return false
	}

	nod := *n
	return key != "" && nod[key] != ""
}

func (n *Node) Get(key string) string {
	if n == nil {
		return ""
	}

	nod := *n
	return nod[key]
}

func (n *Node) Has(lab map[string]string) bool {
	return matcher.Has(*n, lab)
}

func (n *Node) Key() []string {
	if n == nil {
		return nil
	}

	var key []string

	for k := range *n {
		key = append(key, k)
	}

	return key
}

func (n *Node) Len() int {
	if n == nil {
		return 0
	}

	nod := *n
	return len(nod)
}

func (n *Node) Set(key string, val string) {
	nod := *n
	nod[key] = val
}
