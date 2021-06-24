package goweb1

import (
	"testing"

	"github.com/dengjiawen8955/go_utils/restful_util"
	"github.com/dengjiawen8955/goweb/goweb1/goweb"
)

func TestGoweb1Server(t *testing.T) {
	w := goweb.NewWeb("/v1")
	w.Get("/ping", func(c *goweb.Context) {
		c.JSON(restful_util.Ok("PONG"))
	})
	w.Run(8000)
	//test
	//curl http://localhost:8000/v1/ping
}
