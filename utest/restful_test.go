package goweb1

import (
	"fmt"
	"github.com/dengjiawen8955/go_utils/restful_util"
	"github.com/dengjiawen8955/go_utils/string_util"
	"testing"
)

//TestHttp is a http demo.
func TestRestFul(t *testing.T) {
	fmt.Printf("hello=%#v\n", "hello")
	string_util.GetCurrentDirectory()
	b := restful_util.Ok("ok")
	fmt.Printf("data=%#v\n", string(b))
}
