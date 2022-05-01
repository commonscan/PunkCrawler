package coolCrawler

import (
	"coolCrawler/common"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/table"
	"log"
	"os"
	"strings"
)

func getuuid4String() string {
	var geneator, _ = uuid.NewUUID()
	return geneator.String()
}
func (j *Fetcher) OutPutJson(pipe *os.File, output chan Response) {
	var enc = json.NewEncoder(pipe)
	enc.SetEscapeHTML(false)
	for {
		select {
		case response, ok := <-output:
			if ok {
				response.DataUUID = getuuid4String()
				if j.WithIPInfo && len(response.IPv4Addr) > 0 {
					response.IPv4GeoInfo = common.GetIPv4Info(response.IPv4Addr)
				}
				if !j.NoCerts && len(response.Cert) > 0 {
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
		return fmt.Sprintf("%s ... (%d chars more)", str[0:maxLen], len(str))
	}
}

func (j *Fetcher) OutputTable(pipe *os.File, output chan Response) {
	t := table.NewWriter()
	t.SetOutputMirror(pipe)
	t.AppendHeader(table.Row{"#", "URL", "IP", "状态码", "标题", "证书", "WEB服务器"})
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
					t.AppendRow(table.Row{idx, response.SourceURL, response.IPv4Addr, response.StatusCode, echoTitle, cert_str, response.Server})
				}
			} else {
				return
			}
		}
	}
}
