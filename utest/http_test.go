package utest

import (
	"fmt"
	"net/http"
	"testing"
)

//TestHttp is a http demo.
func TestHttp(t *testing.T) {
	addr := ":8888"
	fmt.Printf("addr=%#v\n", addr)
	http.HandleFunc("/ping", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("PONG"))
	})
	http.ListenAndServe(addr, nil)
}
