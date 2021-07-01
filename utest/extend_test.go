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
