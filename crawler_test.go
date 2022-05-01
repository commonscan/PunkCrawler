package coolCrawler

import (
	"coolCrawler/common"
	"fmt"
	"testing"
)

func TestFetcher_DoRequest(t *testing.T) {
	fetcher := Fetcher{}
	resp := fetcher.EnrichTarget("https://qq.com")
	fmt.Println(resp.Title)
}

func TestHasDisableExtension(t *testing.T) {
	fmt.Println(common.UrlHasDisableExtension("http://qq.com/1.mp4"))
}
