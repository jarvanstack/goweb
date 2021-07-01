package utest

import (
	"fmt"
	"testing"
)

type Parent struct {
}
func (p *Parent) p() {
	fmt.Printf("%s\n", "p()")
}
type Son struct {
	*Parent
}
func Test_extend(t *testing.T) {
	son := &Son{}
	son.p()
}

type In interface {
	Hello()
}
type Im struct {
}

func (i *Im) Hello()  {
	fmt.Printf("%s\n", "hello")
}

func Test_implements(t *testing.T) {
	i := &Im{}
	i.Hello()
}
