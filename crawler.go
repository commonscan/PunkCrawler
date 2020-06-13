package coolCrawler

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/imfht/req"
	"github.com/joeguo/tldextract"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Fetcher struct {
	InputFileName  string `short:"i" long:"input-file" description:"输入文件名" default:"-"`
	OutputFileName string `short:"o" long:"output-file" description:"输出文件名" default:"-"`
	ProcessNum     int    `short:"p" long:"process_num" description:"并发数" default:"100"`
	Timeout        int32  `short:"t" long:"timeout" description:"最大超时数(s)" default:"30"`
	Retries        int    `short:"r" long:"retries" description:"最大重试次数" default:"2"`
	WithTitle      bool   `long:"with-title" description:"是否输出Title"`
	WithHTML       bool   `long:"with-html" description:"是否输出HTML"`
	WithTld        bool   `long:"with-tld" description:"是否输出TLD"`
	WithIP         bool   `long:"with-ip" description:"是否输出IP"`
	IconMode       bool   `long:"icon-mode" description:"输出Body的base64和hash"`
	WithHeaders    bool   `long:"with-headers" description:"是否输出Headers"`
	//WithCert       bool   `long:"with-cert" description:"是否输出HTTPS证书"` todo: add  cert support.
	WithLinks bool `long:"with-links" description:"是否输出链接信息"`
}

var (
	cache      = "/tmp/tld.cache"
	extract, _ = tldextract.New(cache, false)
)

func init() {
	req.Client().Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

// enrich HTTP的response: ip\Cert\tld
func (fetcher *Fetcher) EnrichResponse(response Response) Response {
	if !(fetcher.WithTld || fetcher.WithIP) { // do nothing
		return response
	}
	parsedUrl, _ := url.Parse(response.SourceURL)
	var host string
	if strings.Contains(host, ":") {
		host = host[0:strings.Index(host, ":")]
	} else {
		host = parsedUrl.Host
	}
	if fetcher.WithIP {
		if ip, err := getRemoteIP(host); err == nil {
			response.IP = ip
		}
	}
	if fetcher.WithTld {
		response.Tld = getTld(response.URL)
	}
	return response
}
func (fetcher *Fetcher) DoRequest(targetUrl string) Response {
	r := req.New()
	r.MaxReadSize = 1 * 1024 * 1024 * 10 // 10Mb
	req.SetTimeout(time.Duration(fetcher.Timeout) * time.Second)
	var (
		rawResp *req.Resp
		err     error
	)
	if fetcher.Retries < 0 {
		fetcher.Retries = 0
	}

	for i := 0; i <= fetcher.Retries; i++ {
		rawResp, err = req.Get(targetUrl, req.Header{"User-Agent": "Mozilla/5.0 (compatible;Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)"})
		if err == nil {
			break
		}
	}
	if err != nil {
		return Response{Succeed: false, ErrorReason: err.Error(), URL: targetUrl, SourceURL: targetUrl, Time: JSONTime(time.Now())}
	}
	html, _ := rawResp.ToString()
	statusCode := rawResp.Response().StatusCode
	response := Response{
		URL:        rawResp.Request().URL.String(),
		StatusCode: statusCode,
		Succeed:    true,
		Time:       JSONTime(time.Now()),
		SourceURL:  targetUrl,
	}
	if fetcher.IconMode {
		encoded := base64.StdEncoding.EncodeToString(rawResp.Bytes())
		response.B64Content = encoded
		response.Hash = fmt.Sprintf("%x", sha1.Sum(rawResp.Bytes()))
		return response
	}
	if fixedHtml, err := FixEncoding(html, rawResp.Response().Header.Get("Content-Type")); err == nil {
		html = fixedHtml
	}
	if fetcher.WithLinks {
		rawLinks := getLinks(html, targetUrl)
		response.Links = rawLinks
	}
	if fetcher.WithHTML {
		response.Html = html
	}
	if fetcher.WithTitle {
		title := strings.TrimSpace(getTitle(html))
		response.Title = title
	}
	if fetcher.WithHeaders {
		buf := new(bytes.Buffer)
		rawResp.Response().Header.Write(buf)
		response.Headers = buf.String()
	}
	return fetcher.EnrichResponse(response)
}
func (fetcher *Fetcher) Crawl(input chan string, output chan Response, group *sync.WaitGroup) {
	defer group.Done()
	// input chan.
	for {
		select {
		case inputUrl, ok := <-input:
			if ok {
				response := fetcher.DoRequest(inputUrl)
				output <- response
			} else {
				return
			}
		}
	}
}
func (fetcher *Fetcher) DNSInfo(input chan string, output chan DNSInfo, group *sync.WaitGroup) {
	defer group.Done()
	// input chan.
	for {
		select {
		case inputDomain, ok := <-input:
			if ok {
				if response, err := getRemoteIP(inputDomain); err == nil {
					output <- DNSInfo{
						Domain: inputDomain,
						IP:     response,
					}
				}
			} else {
				return
			}
		}
	}
}

func (fetcher *Fetcher) OutputWorker(output chan Response, group *sync.WaitGroup) {
	defer group.Done()
	//var pipe io.Writer
	var pipe *os.File
	var err error
	if fetcher.OutputFileName == "-" {
		pipe = os.Stdout
	} else {
		pipe, err = os.OpenFile(fetcher.OutputFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
	var enc = json.NewEncoder(pipe)
	enc.SetEscapeHTML(false)
	defer pipe.Close()
	defer pipe.Sync()
	for {
		select {
		case response, ok := <-output:
			if ok {
				//buf := new(bytes.Buffer)
				if err = enc.Encode(&response); err != nil {
					log.Fatal(err)
				}
			} else {
				return
			}
		}
	}
	// generate a big file.
}

func (fetcher *Fetcher) Process() {
	inputChan := make(chan string, fetcher.ProcessNum*4)
	outputChan := make(chan Response, fetcher.ProcessNum*4)
	fetchWg := sync.WaitGroup{}
	outputWg := sync.WaitGroup{}
	if fetcher.ProcessNum <= 0 {
		fetcher.ProcessNum = 10
	}
	for i := 0; i < fetcher.ProcessNum; i++ {
		go fetcher.Crawl(inputChan, outputChan, &fetchWg) // start all workers
		fetchWg.Add(1)
	}
	outputWg.Add(1)
	go fetcher.OutputWorker(outputChan, &outputWg)
	var scanner *os.File
	if fetcher.InputFileName == "-" {
		scanner = os.Stdin
	} else {
		var err error
		scanner, err = os.Open(fetcher.InputFileName)
		if err != nil {
			log.Fatal(err)
		}
	}
	f := bufio.NewScanner(scanner)
	for f.Scan() {
		inputTxt := f.Text()
		if !strings.HasPrefix(inputTxt, "http") {
			inputTxt = "http://" + inputTxt
		}
		inputChan <- inputTxt
	}
	close(inputChan)
	fetchWg.Wait()
	close(outputChan)
	outputWg.Wait()
}

func (fetcher *Fetcher) NeedFetch() bool { // 是否需要发送HTTP请求
	if fetcher.WithHTML || fetcher.WithTitle || fetcher.WithLinks {
		return true
	} else {
		return false
	}
}
