package coolCrawler

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"log"
	"os"
	"strings"
)

func (j *Fetcher) OutPutJson(pipe *os.File, output chan Response) {
	var enc = json.NewEncoder(pipe)
	enc.SetEscapeHTML(false)
	for {
		select {
		case response, ok := <-output:
			if ok {
				if !j.WithCert && len(response.Cert) > 0 {
					response.Cert = []string{}
				}
				if err := enc.Encode(&response); err != nil {
					log.Fatal(err)
				}
			} else {
				return
			}
		}
	}
}
func StringWithMax(str string, maxLen int) string {
	if len(str) < maxLen {
		return str
	} else {
		return fmt.Sprintf("%s ...(%d chars)", str[0:maxLen], len(str))
	}
}
func (j *Fetcher) OutputTable(pipe *os.File, output chan Response) {
	t := table.NewWriter()
	t.SetOutputMirror(pipe)
	t.AppendHeader(table.Row{"#", "URL", "状态码", "标题", "证书", "Time"})
	defer t.Render()
	idx := 0
	for {
		select {
		case response, ok := <-output:
			if ok {
				idx += 1
				if response.Succeed {
					cert_byte, _ := json.Marshal(response.Cert)
					cert_str := StringWithMax(string(cert_byte), 32)
					echoTitle := StringWithMax(strings.TrimSpace(response.Title), 32)
					t.AppendRow(table.Row{idx, response.SourceURL, response.StatusCode, echoTitle, cert_str, response.Time.String()})
				}
			} else {
				return
			}
		}
	}

}
