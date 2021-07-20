package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/dengjiawen8955/go_utils/restfulu"
)

//user http office
func Test_sf_test(t *testing.T) {
	http.HandleFunc("/bmft/v2/sf", func(rw http.ResponseWriter, r *http.Request) {
		f, fh, err := r.FormFile("f1")
		fmt.Printf("f: %v\n", f)
		fmt.Printf("fh.Filename: %v\n", fh.Filename)
		fmt.Printf("err: %v\n", err)
		rw.WriteHeader(200)
		rw.Write(restfulu.Ok("PONG"))
	})
	http.ListenAndServe(":8001", nil)
}
