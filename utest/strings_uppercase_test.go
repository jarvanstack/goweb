package utest

import (
	"fmt"
	"strings"
	"testing"
)

func Test_strings_uppercase_test(t *testing.T) {
	upper := strings.ToUpper("content-type")
	fmt.Printf("%s\n", upper)
}
