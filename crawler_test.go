package coolCrawler

import (
	"fmt"
	"testing"
)

func TestFetcher_DoRequest(t *testing.T) {
	fetcher := Fetcher{}
	resp := fetcher.DoRequest("https://qq.com")
	fmt.Println(resp.Title)
}
