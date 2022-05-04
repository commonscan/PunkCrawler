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

func TestStack(t *testing.T) {
	fmt.Println("qq.com IPv4 Available", common.IPv4Available("http://qq.com"))
	fmt.Println("qq.com IPv6 Available", common.IPv6Available("http://qq.com"))

	fmt.Println("www.sjtu.edu.cn IPv4 Available", common.IPv4Available("http://www.sjtu.edu.cn"))
	fmt.Println("www.sjtu.edu.cn IPv6 Available", common.IPv6Available("http://www.sjtu.edu.cn"))
}

func TestIPInfo(t *testing.T) {
	fmt.Println("202.194.14.1 info", common.GetIPv4Info("202.194.7.118"))
}

func TestDescription(t *testing.T) {
	getKeyWordDescription("<html>  <title>PuerkitoBio/goquery: A little like that j-thing, only in Go.</title>\n    <meta name=\"description\" content=\"A little like that j-thing, only in Go. Contribute to PuerkitoBio/goquery development by creating an account on GitHub.\">\n    <link rel=\"search\" type=\"application/opensearchdescription+xml\" href=\"/opensearch.xml\" title=\"GitHub\"> </html>")
}
