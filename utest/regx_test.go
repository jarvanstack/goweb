package utest

import (
	"fmt"
	"testing"

	"github.com/dengjiawen8955/go_utils/stringu"
)

func Test_regx(t *testing.T) {
	str := `Content-Disposition: form-data; name="file"; filename="Snipaste_2021-07-07_15-49-56.png"`
	regx := `Content-Disposition: (.*?); name="(.*?)"; filename="(.*?)"`
	subs := stringu.GetSubStringByRegex(str, regx)
	for _, s := range subs {
		fmt.Printf("%s\n", s)
	}
}
