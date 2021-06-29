package utest

import (
	"fmt"
	"github.com/dengjiawen8955/go_utils/stringu"
	"testing"
)

func Test_base64(t *testing.T) {
	m := stringu.GetMd5ByStr("x3JJHMbDL1EzLkh9GBhXDw==")
	fmt.Printf("m=%#v\n", m)

}
