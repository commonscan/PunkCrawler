package coolCrawler

import (
	"fmt"
	"github.com/imfht/req"
	"testing"
)

func TestFixEncoding(t *testing.T) {
	resp, _ := req.Get("http://www.qq.com")
	fmt.Println(FixEncoding(resp.Bytes(), resp.Response().Header.Get("Content-Type")))
}

func TestGetIP(t *testing.T) {
	fmt.Println(getRemoteIPv4Addr("QQ.COM"))
}
