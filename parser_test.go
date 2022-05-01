package coolCrawler

import (
	"coolCrawler/common"
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
func TestGetIPv6(t *testing.T) {
	fmt.Println(getRemoteIPv6Addr("www.sjtu.edu.cn"))
}
func TestGetTLD(t *testing.T) {
	fmt.Println(getTld("http://www.sjtu.edu.cn"))
}

func TestGeoInfo(t *testing.T) {
	common.GetIPGeoInfo("202.194.14.1")
}
