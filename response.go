package coolCrawler

import (
	"fmt"
	"time"
)

type JSONTime time.Time

type Response struct {
	IPv4Addr      string   `json:"ipv4_addr,omitempty"`
	IPv6Addr      string   `json:"ipv6_addr,omitempty"`
	IPv4Available bool     `json:"ipv4_ok,omitempty"`
	SslOK         bool     `json:"ssl_ok,omitempty"`
	IPv6Available bool     `json:"ipv6_ok,omitempty"`
	IPv4GeoInfo   string   `json:"ipv4_info,omitempty"`
	IPv6GeoInfo   string   `json:"ipv6_info,omitempty"`
	URL           string   `json:"url,omitempty"`
	Html          string   `json:"html,omitempty"`
	CleanedHtml   string   `json:"cleaned_html,omitempty"`
	Title         string   `json:"title,omitempty"`
	StatusCode    int      `json:"status_code,omitempty"`
	Links         []string `json:"links.omitempty"`
	Time          string   `json:"-"`
	TimeStamp     int      `json:"time"`
	Succeed       bool     `json:"succeed,omitempty"`
	ErrorReason   string   `json:"error_reason,omitempty"`
	SourceURL     string   `json:"source_url,omitempty"`
	Tld           string   `json:"tld,omitempty"`
	Domain        string   `json:"domain"`
	Headers       string   `json:"headers,omitempty"`
	Server        string   `json:"server,omitempty"`
	B64Content    string   `json:"b64,omitempty"`
	WebHash       string   `json:"web_hash,omitempty"`
	DataUUID      string   `json:"data_uuid"`
	Cert          []string `json:"cert,omitempty"`
}
type DNSInfo struct {
	Domain string   `json:"domain"`
	IP     []string `json:"ip"`
}

type LinkURL struct {
	URL        string `json:"url"`
	Host       string `json:"host"`
	RootDomain string `json:"root_domain"`
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (t JSONTime) String() string {
	//do your serializing here
	stamp := fmt.Sprintf("%s", time.Time(t).Format("2006-01-02 15:04:05"))
	return stamp
}
