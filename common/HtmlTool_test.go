package common

import (
	"fmt"
	"github.com/imfht/req"
	"testing"
)

func TestHTMLEncoding(t *testing.T) {
	resp, _ := req.Get("http://www.qq.com")
	fmt.Println(GetResponseEncodingRemoved(resp))
}
