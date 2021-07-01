package utest

import (
	"fmt"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	s := "/bmft/v1/ping"
	split := strings.Split(s, "/")
	for i, v := range split {
		fmt.Printf("i=%#v\n", i)
		fmt.Printf("v=%#v\n", v)
	}

}
func Test_split(t *testing.T) {
	s := "/bmft/v1/ping"
	split := strings.Split(s, "/bmft")
	fmt.Printf("split=%#v\n", split)
}