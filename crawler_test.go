package coolCrawler

import (
	"fmt"
	"testing"
)

func TestFetcher_DoRequest(t *testing.T) {
	fetcher := Fetcher{}
	resp := fetcher.DoHTTPRequest("https://qq.com")
	fmt.Println(resp.Title)
}

func TestHasDisableExtension(t *testing.T) {
	fmt.Println(HasDisableExtension("http://qq.com/1.mp4"))
}
