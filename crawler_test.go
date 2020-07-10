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
func BenchmarkFetcher_DoRequest(b *testing.B) {
	b.ResetTimer()
	fetcher := Fetcher{}
	fetcher.WithTitle = true
	b.N = 10000
	count := 0
	for i := 0; i < b.N; i++ {
		resp := fetcher.DoRequest("http://127.0.0.1")
		if resp.Succeed {
			count += 1
		}
	}
	fmt.Println(count, "/", b.N)
}
