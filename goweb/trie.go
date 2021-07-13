package goweb

import "strings"

type node struct {
	keyOfHandler string  // 例如 GET-/bmft/goweb/ping
	currentPath  string  // 目前的路由路径，例如 goweb
	children     []*node // 子节点，例如 [ping, hi, hello]
}

func newNode(currentPath string) *node {
	return &node{
		keyOfHandler: "",
		currentPath: currentPath,
		children: make([]*node,0),
	}
}
func (n *node) getChild(path string) *node {
	for _, child := range n.children {
		if child.currentPath == path {
			return child
		}
	}
	return nil
}
func parsePath(path string) []string {
	split := strings.Split(path, "/")
	return split[1:]
}
func (n *node) insert(keyOfHandler string,paths []string,height int)  {
	//如果最后一层都插入结束
	if height == len(paths) {
		//将 key 放到最后一个叶子节点里面.
		n.keyOfHandler = keyOfHandler
		return
	}
	//继续插入
	currentPath := paths[height]
	child := n.getChild(currentPath)
	if child==nil {
		n.children = append(n.children, newNode(currentPath))
		child = n.getChild(currentPath)
	}
	child.insert(keyOfHandler,paths,height+1)
}
//node 自己是通过 map 得到的 map["GET-/bmft']
//所以传进来的 paths 需要将第一个 /bmft 进行裁剪，否者匹配失败
func (n *node) search(paths []string)*node  {
	tmp := n
	for i := 0; i < len(paths); i++ {
		tmp = tmp.getChild(paths[i])
		if tmp == nil {
			return nil
		}
	}
	return tmp
}

