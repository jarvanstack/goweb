package utest

import (
	"fmt"
	"testing"
)

//进制转化
func Test_base_transform(t *testing.T) {
	fmt.Printf("'0x81'=%d\n", 0x81)
	fmt.Printf("2**3=%#v\n", 1<<3)
}
