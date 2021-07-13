package goweb

import (
	"log"
	"testing"
)

func TestParsePath(t *testing.T) {
	path := parsePath("/goweb/ping")
	log.Printf("path=%#v\n", path)
}
