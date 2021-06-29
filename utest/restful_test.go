package utest

import (
	"fmt"
	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/go_utils/stringu"
	"testing"
)

//TestHttp is a http demo.
func TestRestFul(t *testing.T) {
	fmt.Printf("hello=%#v\n", "hello")
	stringu.GetCurrentDirectory()
	b := restfulu.Ok("ok")
	fmt.Printf("data=%#v\n", string(b))
}
