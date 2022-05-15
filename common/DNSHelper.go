package common

import (
	"context"
	"fmt"
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
	resp, err := client.Head(uri)
	if err != nil {
		return false
	}
	if resp.StatusCode == 200 || (resp.StatusCode == 301 || resp.StatusCode == 302) { // SSLOK: 仅仅200或者30x才认为OK
		return true
	}
	return false
}
