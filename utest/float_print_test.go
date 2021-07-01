package utest

import (
	"fmt"
	"testing"
)

func Test_float_print(t *testing.T) {
	f := 0.1234
	fmt.Printf("f=%.2f\n", f)
}
