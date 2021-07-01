package goweb

import (
	"log"
	"testing"
)

func TestParsePath(t *testing.T) {
	path := parsePath("/v1/ping")
	log.Printf("path=%#v\n", path)
}
