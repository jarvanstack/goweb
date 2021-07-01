package utest

import (
	"fmt"
	"testing"
)

func Test_go_coroutine_debug_test(t *testing.T) {
	x := 0
	go func() {
		x  = 1
		fmt.Printf("x=%#v\n", x)
	}()
	fmt.Printf("x=%#v\n", x)
}