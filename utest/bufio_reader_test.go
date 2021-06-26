package utest

import (
	"bytes"
	"fmt"
	"testing"
)

//bufio 读取 header 接着读取数据 body

func Test_bufio_reader_test(t *testing.T) {
	ds :=`GET /ping HTTP/1.1\r\nHost img.mukewang.com\r\n\r\nbody data`
	reader := bytes.NewReader([]byte(ds))
	bs := make([]byte,2)
	reader.Read(bs)
	//字符长度不是 byte 长度.
	fmt.Printf("reader.Size()=%#v\n", reader.Size())
	fmt.Printf("reader.Len()=%#v\n", reader.Len())

}