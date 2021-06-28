package utest

import (
	"encoding/base64"
	"fmt"
	"github.com/dengjiawen8955/go_utils/string_util"
	"github.com/dengjiawen8955/goweb/goweb1/goweb"
	"testing"
)
//浏览器自动生成 key
//服务器返回 accept
//TdO+OLnE6RKxxgBLbWznsg==
//JB6zD9rWx2LVdjvViIWPY9B8BR4=
func Test_key_to_accept(t *testing.T) {
	key := "dGhlIHNhbXBsZSBub25jZQ=="
	//理论值s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
	key = key + goweb.WsMagicKeyPost
	s1 := string_util.GetSha1ByStr(key)
	sEnc := base64.StdEncoding.EncodeToString([]byte(s1))
	fmt.Println(sEnc)
}
