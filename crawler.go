package coolCrawler

import (
	"bufio"
	"bytes"
	"coolCrawler/common"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/imfht/req"
	"github.com/joeguo/tldextract"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Bool bool

func (b *Bool) UnmarshalFlag(value string) error {
	if value == "true" {
		*b = true
	} else if value == "false" {
		*b = false
	} else {
		return fmt.Errorf("only `true' and `false' are valid values, not `%s'", value)
	}

	return nil
}

func (b Bool) MarshalFlag() string {
	if b {
		return "true"
	}

	return "false"
}

type Fetcher struct {
	InputFileName  string `short:"i" long:"input-file" description:"输入文件名" default:"-"`
	OutputFileName string `short:"o" long:"output-file" description:"输出文件名" default:"-"`
	ProcessNum     int    `short:"p" long:"process-num" description:"并发数" default:"100"`
	Timeout        int32  `short:"t" long:"timeout" description:"最大超时数(s)" default:"30"`
	Retries        int    `short:"r" long:"retries" description:"最大重试次数" default:"2"`

	DefaultHTTPS Bool `long:"default-https" description:"没有协议号的域名默认使用https"`

	PreScan bool   `long:"pre-scan" description:"探测前先端口扫描"`
	Ports   string `long:"ports" description:"扫描的端口，用 ,分割" default:"80,8080,443"`

	FilterBinaryExtensions bool `long:"filter-binary" description:"是否过滤已知的二进制后缀URL"`

	UserAgent   string            `long:"user-agent" description:"User-Agent" default:"Mozilla/5.0 (compatible;Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)"`
	HTTPHeaders map[string]string `long:"http_headers" description:"默认的HTTP Header"`
	HTTPMethod  string            `long:"http_method" description:"http请求方法. eg: GET/POST/PATCH/DELETE/OPTIONS...." default:"GET"`
	HTTPBody    string            `long:"http_body" description:"http body. 当 body为合法json的时候自动使用json提交"`
	HTTPUri     string            `long:"endp" description:"endpoint" default:"/"`

	AbortBinaryHeaders bool `long:"disable-binary-data" description:"启用此选项后，如果response header是已经的二进制，则放弃读取数据."`

	Debug        bool   `long:"debug" description:"向httpbin.org 发送请求以调试发包程序"`
	NoLog        bool   `long:"no-log" description:"不输出log信息"`
	OutputMode   string `long:"output-mode" description:"输出 table而不是json" default:"json"`
	WithTitle    bool   `long:"with-title" description:"是否输出Title"`
	WithHTML     bool   `long:"with-html" description:"是否输出HTML"`
	WithTld      bool   `long:"with-tld" description:"是否输出TLD"`
	WithIPv4     bool   `long:"with-ipv4" description:"是否输出IPv4地址"`
	WithIPv6     bool   `long:"with-ip6" description:"是否输出IPv6地址"`
	WithCert     bool   `long:"with-cert" description:"是否输出HTTPS证书"`
	WithLinks    bool   `long:"with-links" description:"是否输出链接信息"`
	WithHeaders  bool   `long:"with-headers" description:"是否输出Headers"`
	WithGeoInfo  bool   `long:"with-geoinfo" description:"是否输出GEO信息"`
	BodyAsBinary bool   `long:"binary-body" description:"输出Body的base64和hash"`
}

var (
	cache      = "/tmp/tld.cache"
	extract, _ = tldextract.New(cache, false)
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	common.SetUlimitMax()
}

// enrich HTTP的response: ip\Cert\tld
func (fetcher *Fetcher) EnrichResponse(response Response) Response {
	if !(fetcher.WithTld || fetcher.WithIPv4 || fetcher.WithIPv6) { // do nothing
		return response
	}
	if !fetcher.WithHTML {
		response.Html = ""
	}
	if fetcher.WithTld {
		response.Tld = getTld(response.URL)
	}
	return response
}

func (fetcher *Fetcher) EnrichTarget(targetUrl string) Response {
	defer func() {
		if err := recover(); err != nil {
			log.Warn().Msgf("recover. %s . reason %s", targetUrl, err)
		}
	}()
	r := req.New()
	r.Client().Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	r.MaxReadSize = 1 * 1024 * 1024 // 1mb
	r.SetTimeout(time.Duration(fetcher.Timeout) * time.Second)
	var (
		rawResp *req.Resp
		err     error
	)
	if fetcher.Retries < 0 {
		fetcher.Retries = 0
	}
	if fetcher.FilterBinaryExtensions && common.UrlHasDisableExtension(targetUrl) {
		return Response{Succeed: false, URL: targetUrl, SourceURL: targetUrl, Time: JSONTime(time.Now()), ErrorReason: "disabled extensions"}
	}
	httpHeaders := req.Header{}
	for k, v := range fetcher.HTTPHeaders {
		httpHeaders[k] = v
	}
	response := Response{SourceURL: targetUrl}
	parsedUrl, _ := url.Parse(response.SourceURL)
	if strings.Contains(parsedUrl.Host, ":") {
		response.Domain = parsedUrl.Host[0:strings.Index(parsedUrl.Host, ":")]
	} else {
		response.Domain = parsedUrl.Host
	}
	if fetcher.WithIPv4 || fetcher.WithIPv6 { // 需要解析IP数据
		response.IPv4Addr, err = getRemoteIPv4Addr(response.Domain)
		response.IPv6Addr, err = getRemoteIPv6Addr(response.Domain)
	}
	if len(response.IPv4Addr) == 0 && len(response.IPv6Addr) > 0 {
		response.Succeed = false
		response.ErrorReason = "DNSFailed/domain_no_ip"
		return response
	}
	for i := 0; i <= fetcher.Retries; i++ {
		rawResp, err = r.Do(strings.ToUpper(fetcher.HTTPMethod), targetUrl, httpHeaders, fetcher.HTTPBody)
		if err == nil {
			break
		}
	}
	if err != nil {
		errorReason := err.Error()
		if strings.Contains(strings.ToLower(err.Error()), "timeout") {
			errorReason = "Timeout"
		}
		if strings.Contains(strings.ToLower(err.Error()), "connection refused") {
			errorReason = "PortClosed."
		}
		log.Warn().Msgf("failed get %s. error reason: %s", targetUrl, errorReason)
		return Response{Succeed: false, ErrorReason: err.Error(), URL: targetUrl, SourceURL: targetUrl, Time: JSONTime(time.Now())}
	}

	response.StatusCode = rawResp.Response().StatusCode
	response.Succeed = true

	if strings.HasPrefix(rawResp.Request().URL.String(), "https://") {
		if len(rawResp.Response().TLS.PeerCertificates) > 0 {
			response.Cert = rawResp.Response().TLS.PeerCertificates[0].DNSNames // tls info
		}
	}
	response.Hash = fmt.Sprintf("%x", sha1.Sum(rawResp.Bytes()))
	if fetcher.BodyAsBinary {
		encoded := base64.StdEncoding.EncodeToString(rawResp.Bytes())
		response.B64Content = encoded
		return response
	} else {
		if fixedHtml, err := FixEncoding(rawResp.Bytes(), rawResp.Response().Header.Get("Content-Type")); err == nil {
			response.Html = fixedHtml
		}
	}

	if fetcher.WithLinks {
		rawLinks := getLinks(response.Html, targetUrl)
		response.Links = rawLinks
	}
	if fetcher.WithTitle {
		title := strings.TrimSpace(getTitle(response.Html))
		response.Title = title
	}
	if fetcher.WithHeaders {
		buf := new(bytes.Buffer)
		rawResp.Response().Header.Write(buf)
		response.Headers = buf.String()
	}
	log.Info().Msgf("HTTP Request Succeed %s. [title: %s]", targetUrl, response.Title)
	return fetcher.EnrichResponse(response)
}

func (fetcher *Fetcher) DialPortService(hostPort string) (isOpen bool, Service string) {
	d := net.Dialer{Timeout: time.Duration(fetcher.Timeout) * time.Second}
	conn, err := d.Dial("tcp", hostPort)
	if err != nil {
		return false, ""
	} else {
		defer conn.Close()
		_, port, _ := net.SplitHostPort(hostPort)
		if strings.Contains(port, "443") {
			if err := tls.Client(conn, &tls.Config{InsecureSkipVerify: true}).Handshake(); err == nil {
				return true, "https"
			} else {
				log.Warn().Msg(err.Error())
			}
		}
		return true, "http"
	}
}

func (fetcher *Fetcher) Crawl(input chan string, output chan Response, group *sync.WaitGroup) {
	defer group.Done()
	// input chan.
	for {
		select {
		case inputUrl, ok := <-input:
			if !ok {
				return
			} else {
				if fetcher.PreScan { // pre scan mode. 1.  scan ip port before send request.
					if strings.Contains(inputUrl, ":") {
						isOpen, Service := fetcher.DialPortService(inputUrl)
						if isOpen {
							log.Debug().Msgf("found [%s] %s port open ", Service, inputUrl)
						}
						if isOpen && strings.Contains(Service, "http") {
							response := fetcher.EnrichTarget(fmt.Sprintf("%s://%s", Service, inputUrl))
							output <- response
						}
						continue
					}
				}
				if !strings.HasPrefix(strings.ToLower(inputUrl), "http:") && !strings.HasPrefix(strings.ToLower(inputUrl), "https:") {
					service := "http"
					if fetcher.DefaultHTTPS {
						service = "https"
					}
					inputUrl = fmt.Sprintf("%s://%s", service, inputUrl)
				}
				response := fetcher.EnrichTarget(fmt.Sprintf("%s%s", inputUrl, fetcher.HTTPUri))
				output <- response
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
			log.Fatal()
		}
	}
	if fetcher.OutputMode == "table" {
		fetcher.OutputTable(pipe, output)
	} else {
		fetcher.OutPutJson(pipe, output)
	}
	defer pipe.Close()
	defer pipe.Sync()
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
			log.Fatal()
		}
	}
	f := bufio.NewScanner(scanner)
	for f.Scan() {
		inputUrl := f.Text()
		if fetcher.PreScan { // 端口扫描模式
			if common.IsCIDR(inputUrl) {
				ips, err := common.GenerateIP(inputUrl)
				if err == nil {
					log.Info().Msgf("convert cidr [%s] -> [%d] ip", inputUrl, len(ips))
					for _, ip := range ips {
						for _, port := range strings.Split(fetcher.Ports, ",") {
							inputChan <- fmt.Sprintf("%s:%s", ip, port)
						}
					}
				} else {
					log.Err(err)
				}
			} else if common.IsIp(inputUrl) {
				for _, port := range strings.Split(fetcher.Ports, ",") {
					inputChan <- fmt.Sprintf("%s:%s", inputUrl, port)
				}
			} else {
				inputChan <- inputUrl
			}
		} else {
			inputChan <- inputUrl
		}

	}
	close(inputChan)
	fetchWg.Wait()
	close(outputChan)
	outputWg.Wait()
	log.Debug().Msg("exit.")
}

func (fetcher *Fetcher) NeedFetch() bool { // 是否需要发送HTTP请求
	if fetcher.WithHTML || fetcher.WithTitle || fetcher.WithLinks {
		return true
	} else {
		return false
	}
}
