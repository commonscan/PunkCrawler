package main

import (
	"coolCrawler"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"log"
)

func main() {
	var fetcher coolCrawler.Fetcher
	_, err := flags.Parse(&fetcher)
	if err != nil {
		log.Fatal(err)
	}
	if fetcher.NoLog {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
	fetcher.Process()
}
