package utest

import (
	"fmt"
	"testing"

	"github.com/dengjiawen8955/go_utils/stringu"
)

func Test_regx(t *testing.T) {
	str := `Content-Disposition: form-data; name="file"; filename="Snipaste.png"`
	regx := `Content-Disposition: (\S*?); name="(\S*?)"; filename="(\S*?)"`
	subs, _ := stringu.GetSubStringByRegex(str, regx)
	for _, s := range subs {
		fmt.Printf("%s\n", s)
	}
	fmt.Printf("%s\n", "-------------")
	str2 := `Content-Type: image/png`
	regx2 := `Content-Type: (\S*)`
	subs2, _ := stringu.GetSubStringByRegex(str2, regx2)
	fmt.Printf("=%#v\n", len(subs2))
	for _, element := range subs2 {
		fmt.Printf("%s\n", element)
	}
}

func Test_line2_regx_test(t *testing.T) {
	reg1 := `name="(\S*?)"`
	subs, _ := stringu.GetSubStringByRegex(`name="uploadFile"`, reg1)
	fmt.Printf("subs=%#v\n", subs)
	reg1 = `filename="(\S*?)"`
	subs, _ = stringu.GetSubStringByRegex(`filename="f1.txt"`, reg1)
	fmt.Printf("subs=%#v\n", subs)

}
