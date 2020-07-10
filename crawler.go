package coolCrawler

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gojektech/heimdall/v6"
	"github.com/gojektech/heimdall/v6/httpclient"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/joeguo/tldextract"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Fetcher struct {
	InputFileName          string `short:"i" long:"input-file" description:"输入文件名" default:"-"`
	OutputFileName         string `short:"o" long:"output-file" description:"输出文件名" default:"-"`
	ProcessNum             int    `short:"p" long:"process-num" description:"并发数" default:"100"`
	Timeout                int32  `short:"t" long:"timeout" description:"最大超时数(s)" default:"30"`
	Retries                int    `short:"r" long:"retries" description:"最大重试次数" default:"2"`
	WithTitle              bool   `long:"with-title" description:"是否输出Title"`
	WithHTML               bool   `long:"with-html" description:"是否输出HTML"`
	WithTld                bool   `long:"with-tld" description:"是否输出TLD"`
	WithIP                 bool   `long:"with-ip" description:"是否输出IP"`
	IconMode               bool   `long:"icon-mode" description:"输出Body的base64和hash"`
	WithHeaders            bool   `long:"with-headers" description:"是否输出Headers"`
	UserAgent              string `long:"user-agent" description:"User-Agent" default:"Mozilla/5.0 (compatible;Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)"`
	WithCert               bool   `long:"with-cert" description:"是否输出HTTPS证书"`
	WithLinks              bool   `long:"with-links" description:"是否输出链接信息"`
	WithText               bool   `long:"with-text" description:"是否解析纯文本"`
	FilterBinaryExtensions bool   `long:"filter-binary" description:"是否输出链接信息"`
}

var (
	cache           = "/tmp/tld.cache"
	extract, _      = tldextract.New(cache, false)
	backoffInterval = 2 * time.Millisecond
	// Define a maximum jitter interval. It must be more than 1*time.Millisecond
	maximumJitterInterval = 5 * time.Millisecond
	backoff               = heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)

	retrier = heimdall.NewRetrier(backoff)
	//binaryExtensions = []string{""}
)

func init() {

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
	timeout := 1000 * time.Millisecond

	// Use the clients GET method to create and execute the request
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(4),
	)
	m := make(map[string]string)
	m["User-Agent"] = "Baiduspider+(+http://www.baidu.com/search/spider.htm)"
	header := http.Header{}
	header.Add("User-Agent", "qq.com")
	res, err := client.Get(targetUrl, header)
	if err != nil {
		return Response{}
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Response{}
	}
	if err != nil {
		panic(err)
	}
	if fetcher.Retries < 0 {
		fetcher.Retries = 0
	}
	if err != nil {
		return Response{Succeed: false, ErrorReason: err.Error(), URL: targetUrl, SourceURL: targetUrl, Time: JSONTime(time.Now())}
	}
	html := string(body)
	statusCode := res.StatusCode
	response := Response{
		URL:        res.Request.URL.String(),
		StatusCode: statusCode,
		Succeed:    true,
		Time:       JSONTime(time.Now()),
		SourceURL:  targetUrl,
	}
	if fetcher.WithCert && strings.HasPrefix(res.Request.URL.String(), "https://") {
		var cert_interface map[string]interface{}
		inrec, _ := json.Marshal(res.TLS.PeerCertificates[0])
		err := json.Unmarshal(inrec, &cert_interface)
		if err != nil {
			response.Cert = cert_interface
		}

	}
	if fetcher.IconMode {
		encoded := base64.StdEncoding.EncodeToString(body)
		response.B64Content = encoded
		response.Hash = fmt.Sprintf("%x", sha1.Sum(body))
		return response
	}
	if fixedHtml, err := FixEncoding(html, res.Header.Get("Content-Type")); err == nil {
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
		item, _ := httputil.DumpResponse(res, false)
		response.Headers = string(item)
	}
	if fetcher.WithText {
		//p := bluemonday.StripTagsPolicy()
		//html := p.Sanitize(string(body))
		html := strip.StripTags(html)
		if strings.HasPrefix(html, "<!DOCTYPE html>") {
			html = strings.TrimLeft(html, "<!DOCTYPE html>")
		}
		// <!doctype html>
		//html = strings.ReplaceAll(html, "\n ", "\n")
		//html = strings.ReplaceAll(html, "\r\n ", "\n")
		response.Text = removeLBR(removeDoctype(html))
	}
	return fetcher.EnrichResponse(response)
}
func removeDoctype(text string) string {
	re := regexp.MustCompile(`(?i)<!doctype html>[ ]*`)
	return re.ReplaceAllString(text, ``)

}
func removeLBR(text string) string {
	re := regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`)
	return re.ReplaceAllString(text, ``)
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
