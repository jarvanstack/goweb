package utest

import (
	"fmt"
	"github.com/dengjiawen8955/go_utils/restful_util"
	"github.com/dengjiawen8955/go_utils/test_util"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func Test_office_server(t *testing.T) {
	http.HandleFunc("/v2/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type","application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(restful_util.Ok("PONG"))
	})
	http.ListenAndServe("localhost:8889",nil)
}
// 类型                    RPS
//tcp to http纯:           1917/2088/1836
//go 官方 http             1939/2097/1946
//goweb 框架               1985/1981/2080
//nice
func Test_office_vs_tcp_to_http_vs_goweb(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	url1 := "http://127.0.0.1:9000/v1/ping"
	url2 := "http://127.0.0.1:8889/v2/ping"
	url3 := "http://127.0.0.1:8888/v1/ping"
	method := "GET"
	loopTimes := 1000
	failTimes := 0
	tu := test_util.NewTestUtil(uint32(loopTimes))
	tu.StartWithComment("预热数据")
	for i := 0; i < loopTimes; i++ {
		reqSuccess := httpReq(url1, method)
		if !reqSuccess {
			failTimes++
		}
	}
	tu.End()
	fmt.Printf("failTimes=%#v\n", failTimes)
	//-----------
	failTimes = 0
	tu = test_util.NewTestUtil(uint32(loopTimes))
	tu.StartWithComment(url1)
	for i := 0; i < loopTimes; i++ {
		reqSuccess := httpReq(url1, method)
		if !reqSuccess {
			failTimes++
		}
	}
	tu.End()
	fmt.Printf("failTimes=%#v\n", failTimes)
	//-----------

	failTimes = 0
	tu = test_util.NewTestUtil(uint32(loopTimes))
	tu.StartWithComment(url2)
	for i := 0; i < loopTimes; i++ {
		reqSuccess := httpReq(url2, method)
		if !reqSuccess {
			failTimes++
		}
	}
	tu.End()
	fmt.Printf("failTimes=%#v\n", failTimes)
	//-----------
	failTimes = 0
	tu = test_util.NewTestUtil(uint32(loopTimes))
	tu.StartWithComment(url3)
	for i := 0; i < loopTimes; i++ {
		reqSuccess := httpReq(url3, method)
		if !reqSuccess {
			failTimes++
		}
	}
	tu.End()
	fmt.Printf("failTimes=%#v\n", failTimes)
}
// 类型                    RPS
//tcp to http纯:           1846
//go 官方 http
func Test_office_vs_goweb(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	url1 := "http://127.0.0.1:8889/v2/ping"
	url2 := "http://127.0.0.1:8888/v1/ping"
	method := "GET"
	loopTimes := 1000
	failTimes := 0
	tu := test_util.NewTestUtil(uint32(loopTimes))
	tu.StartWithComment(url2)
	for i := 0; i < loopTimes; i++ {
		reqSuccess := httpReq(url2, method)
		if !reqSuccess {
			failTimes++
		}
	}
	tu.End()
	fmt.Printf("failTimes=%#v\n", failTimes)
	//-----------
	failTimes = 0
	tu = test_util.NewTestUtil(uint32(loopTimes))
	tu.StartWithComment(url1)
	for i := 0; i < loopTimes; i++ {
		reqSuccess := httpReq(url1, method)
		if !reqSuccess {
			failTimes++
		}
	}
	tu.End()
	fmt.Printf("failTimes=%#v\n", failTimes)
}
func httpReq(url,method string) bool {
	resp, _ := http.Get("http://127.0.0.1:8000/go")
	ioutil.ReadAll(resp.Body)
	return true

}
func Test_function_return_print(t *testing.T) {
	url := "http://localhost:8889/v2/ping"
	method := "GET"
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}