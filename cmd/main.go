package main

import (
	"coolCrawler"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
)

func main() {
	var fetcher coolCrawler.Fetcher
	_, err := flags.Parse(&fetcher)
	if err != nil {
		//fmt.Println(err)
		return
	}
	if fetcher.NoLog {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
	fetcher.Process()
}
