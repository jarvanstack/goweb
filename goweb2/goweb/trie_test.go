package goweb

import (
	"fmt"
	"testing"
)
//
func Test_trie(t *testing.T) {
	n := &node{}
	path := "/p/:lang/doc"
	strs := parsePattern(path)
	n.insert("/p/:lang",strs,len(strs))
	fmt.Printf("n=%#v\n", n)
	//n.insert()
}
