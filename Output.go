package coolCrawler

import (
	"encoding/json"
	"github.com/jedib0t/go-pretty/table"
	"log"
	"os"
	"strings"
)

func OutPutJson(pipe *os.File, output chan Response) {
	var enc = json.NewEncoder(pipe)
	enc.SetEscapeHTML(false)
	for {
		select {
		case response, ok := <-output:
			if ok {
				if err := enc.Encode(&response); err != nil {
					log.Fatal(err)
				}
			} else {
				return
			}
		}
	}
}

func OutputTable(pipe *os.File, output chan Response) {
	t := table.NewWriter()
	t.SetOutputMirror(pipe)
	t.AppendHeader(table.Row{"#", "URL", "状态码", "标题", "Time"})
	defer t.Render()
	idx := 0
	for {
		select {
		case response, ok := <-output:
			if ok {
				idx += 1
				if response.Succeed {
					t.AppendRow(table.Row{idx, response.SourceURL, response.StatusCode, strings.TrimSpace(response.Title), response.Time.String()})
				}
			} else {
				return
			}
		}
	}

}
