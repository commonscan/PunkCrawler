package common

import (
	"context"
	"fmt"
	"github.com/imfht/req"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"net"

	"net/http"
	"time"
)

var (
	zeroDialer     net.Dialer
	Ipv4httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	Ipv6httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func init() {
	IPv4 := http.DefaultTransport.(*http.Transport).Clone()
	IPv4.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return zeroDialer.DialContext(ctx, "tcp4", addr)
	}
	Ipv4httpClient.Transport = IPv4
	IPv6 := http.DefaultTransport.(*http.Transport).Clone()
	IPv6.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return zeroDialer.DialContext(ctx, "tcp6", addr)
	}
	Ipv6httpClient.Transport = IPv6
}
func IPv4Available(url string) bool {
	_, err := Ipv4httpClient.Head(url)
	if err == nil {
		return true
	}
	return false
}
func IPv6Available(url string) bool {
	_, err := Ipv6httpClient.Head(url)
	if err == nil {
		return true
	}
	return false
}
func SSLAvailable(domain string) bool {
	uri := fmt.Sprintf("https://%s/", domain)
	client := http.Client{Timeout: 10 * time.Second}
	_, err := client.Head(uri)
	if err == nil {
		return true
	} else {
		return false
	}
}

func GetIPv4Info(ipv4addr string) string {
	var req = req.Req{}
	req.SetTimeout(time.Duration(time.Second * 10))
	resp, err := req.Get("http://43.131.50.200:8888/q?ip=" + ipv4addr)
	if err != nil {
		log.Warn().Msgf("获取IP信息失败，url %s", "http://localhost:8888/q?ip="+ipv4addr)
		return ""
	}
	data := resp.String()
	return fmt.Sprintf("%s|%s|%s|%s", gjson.Get(data, "country_name"), gjson.Get(data, "region_name"),
		gjson.Get(data, "city_name"), gjson.Get(data, "isp_domain"))
}
