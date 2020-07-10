package coolCrawler

import (
	"fmt"
	"time"
)

type JSONTime time.Time

type Response struct {
	IP          []string               `json:"ip,omitempty"`
	URL         string                 `json:"url,omitempty"`
	Html        string                 `json:"html,omitempty"`
	Title       string                 `json:"title,omitempty"`
	StatusCode  int                    `json:"status_code,omitempty"`
	Links       []string               `json:"links,omitempty"`
	Time        JSONTime               `json:"time,omitempty"`
	Succeed     bool                   `json:"succeed,omitempty"`
	ErrorReason string                 `json:"error_reason,omitempty"`
	SourceURL   string                 `json:"source_url,omitempty"`
	Tld         string                 `json:"tld,omitempty"`
	Headers     string                 `json:"headers,omitempty"`
	B64Content  string                 `json:"b64,omitempty"`
	Hash        string                 `json:"hash,omitempty"`
	Text        string                 `json:"text,omitempty"`
	Cert        map[string]interface{} `json:"cert,omitempty"`
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
