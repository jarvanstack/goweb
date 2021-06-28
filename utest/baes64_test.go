package utest

import (
	"fmt"
	"github.com/dengjiawen8955/go_utils/string_util"
	"testing"
)

func Test_base64(t *testing.T) {
	m := string_util.GetMd5ByStr("x3JJHMbDL1EzLkh9GBhXDw==")
	fmt.Printf("m=%#v\n", m)

}
