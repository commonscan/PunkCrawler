package main

import (
	"coolCrawler"
	"github.com/jessevdk/go-flags"
	"log"
)

func main() {
	var fetcher coolCrawler.Fetcher
	_, err := flags.Parse(&fetcher)
	if err != nil {
		log.Fatal(err)
	}
	fetcher.Process()
}
